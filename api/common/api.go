// Package common implements all the things that an IPFS Cluster API component
// must do, except the actual routes that it handles.
//
// This is meant for re-use when implementing actual REST APIs by saving most
// of the efforts and automatically getting a lot of the setup and things like
// authentication handled.
//
// The API exposes the routes in two ways: the first is through a regular
// HTTP(s) listener. The second is by tunneling HTTP through a libp2p stream
// (thus getting an encrypted channel without the need to setup TLS). Both
// ways can be used at the same time, or disabled.
//
// This is used by rest and pinsvc packages.
package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	types "github.com/ipfs/ipfs-cluster/api"

	cid "github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	gopath "github.com/ipfs/go-path"
	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	rpc "github.com/libp2p/go-libp2p-gorpc"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	noise "github.com/libp2p/go-libp2p-noise"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	manet "github.com/multiformats/go-multiaddr/net"

	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Common errors
var (
	// ErrNoEndpointEnabled is returned when the API is created but
	// no HTTPListenAddr, nor libp2p configuration fields, nor a libp2p
	// Host are provided.
	ErrNoEndpointsEnabled = errors.New("neither the libp2p nor the HTTP endpoints are enabled")

	// ErrHTTPEndpointNotEnabled is returned when trying to perform
	// operations that rely on the HTTPEndpoint but it is disabled.
	ErrHTTPEndpointNotEnabled = errors.New("the HTTP endpoint is not enabled")
)

// When passed to SendResponse(), it will figure out which http status
// to set by itself.
const SetStatusAutomatically = -1

// API implements an API and aims to provides
// a RESTful HTTP API for Cluster.
type API struct {
	ctx    context.Context
	cancel func()

	config *Config

	rpcClient *rpc.Client
	rpcReady  chan struct{}
	router    *mux.Router
	routes    func(*rpc.Client) []Route

	server *http.Server
	host   host.Host

	httpListeners  []net.Listener
	libp2pListener net.Listener

	shutdownLock sync.Mutex
	shutdown     bool
	wg           sync.WaitGroup
}

// Route defines a REST endpoint supported by this API.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type logWriter struct {
	logger *logging.ZapEventLogger
}

func (lw logWriter) Write(b []byte) (int, error) {
	lw.logger.Info(string(b))
	return len(b), nil
}

// NewAPI creates a new common API component with the given configuration.
func NewAPI(ctx context.Context, cfg *Config, routes func(*rpc.Client) []Route) (*API, error) {
	return NewAPIWithHost(ctx, cfg, nil, routes)
}

// NewAPIWithHost creates a new common API component and enables
// the libp2p-http endpoint using the given Host, if not nil.
func NewAPIWithHost(ctx context.Context, cfg *Config, h host.Host, routes func(*rpc.Client) []Route) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	// Our handler is a gorilla router wrapped with:
	// - a custom strictSlashHandler that uses 307 redirects (#1415)
	// - the cors handler,
	// - the basic auth handler.
	//
	// Thus every request will need to have valid credentials first, then
	// comply with CORS, then it may be redirected if the path ends with a
	// "/" and finally it hits one of our routes and handlers.
	router := mux.NewRouter()
	handler := basicAuthHandler(
		cfg.BasicAuthCredentials,
		cors.New(*cfg.CorsOptions()).
			Handler(
				strictSlashHandler(router),
			),
		cfg.Logger,
	)
	if cfg.Tracing {
		handler = &ochttp.Handler{
			IsPublicEndpoint: true,
			Propagation:      &tracecontext.HTTPFormat{},
			Handler:          handler,
			StartOptions:     trace.StartOptions{SpanKind: trace.SpanKindServer},
			FormatSpanName:   func(req *http.Request) string { return req.Host + ":" + req.URL.Path + ":" + req.Method },
		}
	}

	writer, err := cfg.LogWriter()
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		Handler:           handlers.LoggingHandler(writer, handler),
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	// See: https://github.com/ipfs/go-ipfs/issues/5168
	// See: https://github.com/ipfs/ipfs-cluster/issues/548
	// on why this is re-enabled.
	s.SetKeepAlivesEnabled(true)
	s.MaxHeaderBytes = cfg.MaxHeaderBytes

	ctx, cancel := context.WithCancel(ctx)

	api := &API{
		ctx:      ctx,
		cancel:   cancel,
		config:   cfg,
		server:   s,
		host:     h,
		router:   router,
		routes:   routes,
		rpcReady: make(chan struct{}, 2),
	}

	// Set up api.httpListeners if enabled
	err = api.setupHTTP()
	if err != nil {
		return nil, err
	}

	// Set up api.libp2pListeners if enabled
	err = api.setupLibp2p()
	if err != nil {
		return nil, err
	}

	if len(api.httpListeners) == 0 && api.libp2pListener == nil {
		return nil, ErrNoEndpointsEnabled
	}

	api.run(ctx)
	return api, nil
}

func (api *API) setupHTTP() error {
	if len(api.config.HTTPListenAddr) == 0 {
		return nil
	}

	for _, listenMAddr := range api.config.HTTPListenAddr {
		n, addr, err := manet.DialArgs(listenMAddr)
		if err != nil {
			return err
		}

		var l net.Listener
		if api.config.TLS != nil {
			l, err = tls.Listen(n, addr, api.config.TLS)
		} else {
			l, err = net.Listen(n, addr)
		}
		if err != nil {
			return err
		}
		api.httpListeners = append(api.httpListeners, l)
	}
	return nil
}

func (api *API) setupLibp2p() error {
	// Make new host. Override any provided existing one
	// if we have config for a custom one.
	if len(api.config.Libp2pListenAddr) > 0 {
		// We use a new host context. We will call
		// Close() on shutdown(). Avoids things like:
		// https://github.com/ipfs/ipfs-cluster/issues/853
		h, err := libp2p.New(
			libp2p.Identity(api.config.PrivateKey),
			libp2p.ListenAddrs(api.config.Libp2pListenAddr...),
			libp2p.Security(noise.ID, noise.New),
			libp2p.Security(libp2ptls.ID, libp2ptls.New),
		)
		if err != nil {
			return err
		}
		api.host = h
	}

	if api.host == nil {
		return nil
	}

	l, err := gostream.Listen(api.host, p2phttp.DefaultP2PProtocol)
	if err != nil {
		return err
	}
	api.libp2pListener = l
	return nil
}

func (api *API) addRoutes() {
	for _, route := range api.routes(api.rpcClient) {
		api.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(
				ochttp.WithRouteTag(
					http.HandlerFunc(route.HandlerFunc),
					"/"+route.Name,
				),
			)
	}
	api.router.NotFoundHandler = ochttp.WithRouteTag(
		http.HandlerFunc(api.notFoundHandler),
		"/notfound",
	)
}

// basicAuth wraps a given handler with basic authentication
func basicAuthHandler(credentials map[string]string, h http.Handler, lggr *logging.ZapEventLogger) http.Handler {
	if credentials == nil {
		return h
	}

	wrap := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		username, password, ok := r.BasicAuth()
		if !ok {
			resp, err := unauthorizedResp()
			if err != nil {
				lggr.Error(err)
				return
			}
			http.Error(w, resp, http.StatusUnauthorized)
			return
		}

		authorized := false
		for u, p := range credentials {
			if u == username && p == password {
				authorized = true
			}
		}
		if !authorized {
			resp, err := unauthorizedResp()
			if err != nil {
				lggr.Error(err)
				return
			}
			http.Error(w, resp, http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(wrap)
}

// The Gorilla muxer StrictSlash option uses a 301 permanent redirect, which
// results in POST requests becoming GET requests in most clients.  Thus we
// use our own middleware that performs a 307 redirect.  See issue #1415 for
// more details.
func strictSlashHandler(h http.Handler) http.Handler {
	wrap := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/") {
			u, _ := url.Parse(r.URL.String())
			u.Path = u.Path[:len(u.Path)-1]
			http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
			return
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(wrap)
}

func unauthorizedResp() (string, error) {
	apiError := &types.Error{
		Code:    401,
		Message: "Unauthorized",
	}
	resp, err := json.Marshal(apiError)
	return string(resp), err
}

func (api *API) run(ctx context.Context) {
	api.wg.Add(len(api.httpListeners))
	for _, l := range api.httpListeners {
		go func(l net.Listener) {
			defer api.wg.Done()
			api.runHTTPServer(ctx, l)
		}(l)
	}

	if api.libp2pListener != nil {
		api.wg.Add(1)
		go func() {
			defer api.wg.Done()
			api.runLibp2pServer(ctx)
		}()
	}
}

// runs in goroutine from run()
func (api *API) runHTTPServer(ctx context.Context, l net.Listener) {
	select {
	case <-api.rpcReady:
	case <-api.ctx.Done():
		return
	}

	maddr, err := manet.FromNetAddr(l.Addr())
	if err != nil {
		api.config.Logger.Error(err)
	}

	api.config.Logger.Infof(strings.ToUpper(api.config.ConfigKey)+" (HTTP): %s", maddr)
	err = api.server.Serve(l)
	if err != nil && !strings.Contains(err.Error(), "closed network connection") {
		api.config.Logger.Error(err)
	}
}

// runs in goroutine from run()
func (api *API) runLibp2pServer(ctx context.Context) {
	select {
	case <-api.rpcReady:
	case <-api.ctx.Done():
		return
	}

	listenMsg := ""
	for _, a := range api.host.Addrs() {
		listenMsg += fmt.Sprintf("        %s/p2p/%s\n", a, api.host.ID().Pretty())
	}

	api.config.Logger.Infof(strings.ToUpper(api.config.ConfigKey)+" (libp2p-http): ENABLED. Listening on:\n%s\n", listenMsg)

	err := api.server.Serve(api.libp2pListener)
	if err != nil && !strings.Contains(err.Error(), "context canceled") {
		api.config.Logger.Error(err)
	}
}

// Shutdown stops any API listeners.
func (api *API) Shutdown(ctx context.Context) error {
	_, span := trace.StartSpan(ctx, "api/Shutdown")
	defer span.End()

	api.shutdownLock.Lock()
	defer api.shutdownLock.Unlock()

	if api.shutdown {
		api.config.Logger.Debug("already shutdown")
		return nil
	}

	api.config.Logger.Info("stopping Cluster API")

	api.cancel()
	close(api.rpcReady)

	// Cancel any outstanding ops
	api.server.SetKeepAlivesEnabled(false)

	for _, l := range api.httpListeners {
		l.Close()
	}

	if api.libp2pListener != nil {
		api.libp2pListener.Close()
	}

	api.wg.Wait()

	// This means we created the host
	if api.config.Libp2pListenAddr != nil {
		api.host.Close()
	}
	api.shutdown = true
	return nil
}

// SetClient makes the component ready to perform RPC
// requests.
func (api *API) SetClient(c *rpc.Client) {
	api.rpcClient = c
	api.addRoutes()

	// One notification for http server and one for libp2p server.
	api.rpcReady <- struct{}{}
	api.rpcReady <- struct{}{}
}

func (api *API) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	api.SendResponse(w, http.StatusNotFound, errors.New("not found"), nil)
}

// ParsePinPathOrFail parses a pin path and returns it or makes the request
// fail.
func (api *API) ParsePinPathOrFail(w http.ResponseWriter, r *http.Request) *types.PinPath {
	vars := mux.Vars(r)
	urlpath := "/" + vars["keyType"] + "/" + strings.TrimSuffix(vars["path"], "/")

	path, err := gopath.ParsePath(urlpath)
	if err != nil {
		api.SendResponse(w, http.StatusBadRequest, errors.New("error parsing path: "+err.Error()), nil)
		return nil
	}

	pinPath := &types.PinPath{Path: path.String()}
	err = pinPath.PinOptions.FromQuery(r.URL.Query())
	if err != nil {
		api.SendResponse(w, http.StatusBadRequest, err, nil)
	}
	return pinPath
}

// ParsePidOrFail parses a Cid and returns it or makes the request fail.
func (api *API) ParseCidOrFail(w http.ResponseWriter, r *http.Request) *types.Pin {
	vars := mux.Vars(r)
	hash := vars["hash"]

	c, err := cid.Decode(hash)
	if err != nil {
		api.SendResponse(w, http.StatusBadRequest, errors.New("error decoding Cid: "+err.Error()), nil)
		return nil
	}

	opts := types.PinOptions{}
	err = opts.FromQuery(r.URL.Query())
	if err != nil {
		api.SendResponse(w, http.StatusBadRequest, err, nil)
	}
	pin := types.PinWithOpts(c, opts)
	pin.MaxDepth = -1 // For now, all pins are recursive
	return pin
}

// ParsePidOrFail parses a PID and returns it or makes the request fail.
func (api *API) ParsePidOrFail(w http.ResponseWriter, r *http.Request) peer.ID {
	vars := mux.Vars(r)
	idStr := vars["peer"]
	pid, err := peer.Decode(idStr)
	if err != nil {
		api.SendResponse(w, http.StatusBadRequest, errors.New("error decoding Peer ID: "+err.Error()), nil)
		return ""
	}
	return pid
}

// SendResponse wraps all the logic for writing the response to a request:
// * Write configured headers
// * Write application/json content type
// * Write status: determined automatically if given "SetStatusAutomatically"
// * Write an error if there is or write the response if there is
func (api *API) SendResponse(
	w http.ResponseWriter,
	status int,
	err error,
	resp interface{},
) {

	api.SetHeaders(w)
	enc := json.NewEncoder(w)

	// Send an error
	if err != nil {
		if status == SetStatusAutomatically || status < 400 { // set a default error status
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

		errorResp := types.Error{
			Code:    status,
			Message: err.Error(),
		}
		api.config.Logger.Errorf("sending error response: %d: %s", status, err.Error())

		if err := enc.Encode(errorResp); err != nil {
			api.config.Logger.Error(err)
		}
		return
	}

	// Send a body
	if resp != nil {
		if status == SetStatusAutomatically {
			status = http.StatusOK
		}

		w.WriteHeader(status)

		if err = enc.Encode(resp); err != nil {
			api.config.Logger.Error(err)
		}
		return
	}

	// Empty response
	if status == SetStatusAutomatically {
		status = http.StatusNoContent
	}

	w.WriteHeader(status)
}

// SetsHeaders sets all the headers that are common to all responses
// from this API. Called automatically from SendResponse().
func (api *API) SetHeaders(w http.ResponseWriter) {
	for header, values := range api.config.Headers {
		for _, val := range values {
			w.Header().Add(header, val)
		}
	}

	w.Header().Add("Content-Type", "application/json")
}

// These functions below are mostly used in tests.

// HTTPAddresses returns the HTTP(s) listening address
// in host:port format. Useful when configured to start
// on a random port (0). Returns error when the HTTP endpoint
// is not enabled.
func (api *API) HTTPAddresses() ([]string, error) {
	if len(api.httpListeners) == 0 {
		return nil, ErrHTTPEndpointNotEnabled
	}
	var addrs []string
	for _, l := range api.httpListeners {
		addrs = append(addrs, l.Addr().String())
	}

	return addrs, nil
}

// Host returns the libp2p Host used by the API, if any.
// The result is either the host provided during initialization,
// a default Host created with options from the configuration object,
// or nil.
func (api *API) Host() host.Host {
	return api.host
}

// Headers returns the configured Headers.
// Useful for testing.
func (api *API) Headers() map[string][]string {
	return api.config.Headers
}

// Controls the HTTP server Keep Alive settings.
// Useful for testing.
func (api *API) SetKeepAlivesEnabled(b bool) {
	api.server.SetKeepAlivesEnabled(b)
}
