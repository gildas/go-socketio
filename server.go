package socketio

import (
	"net/http"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

// Server defines the Socket.io server
type Server struct {
	http.Handler

	namespaces map[string]*Namespace
	sockets    *socketMap
	transports map[string]Transport
	logger     *logger.Logger
}

// ServerOptions defines the options for creating a Server
type ServerOptions struct {
	Logger *logger.Logger
}

// NewServer creates a new Server
func NewServer(options *ServerOptions) (*Server, error) {
	return &Server{
		namespaces: map[string]*Namespace{},
		sockets:    NewSocketMap(),
		transports: map[string]Transport{
			"polling": NewPollingTransport(&PollingTransportOptions{ Logger: options.Logger }),
		},
		logger: logger.CreateIfNil(options.Logger, "socketio").Child("socketio", "socketio"),
	}, nil
}

// Close stops the server
func (server *Server) Close() {
}

// OnConnect sets the function to be called when a connection event is received for the given namespace
func (server *Server) OnConnect(namespace string, callback ConnectHandler) {
	ns, found := server.namespaces[namespace]
	if !found {
		ns = NewNamespace(namespace)
		server.namespaces[namespace] = ns
	}
	ns.connectHandler = callback
}

// OnDisconnect sets the function to be called when a disconnection event is received for the given namespace
func (server *Server) OnDisconnect(namespace string, callback DisconnectHandler) {
	ns, found := server.namespaces[namespace]
	if !found {
		ns = NewNamespace(namespace)
		server.namespaces[namespace] = ns
	}
	ns.disconnectHandler = callback
}

// OnError sets the function to be called when an error event is received for the given namespace
func (server *Server) OnError(namespace string, callback ErrorHandler) {
	ns, found := server.namespaces[namespace]
	if !found {
		ns = NewNamespace(namespace)
		server.namespaces[namespace] = ns
	}
	ns.errorHandler = callback
}

// On sets the function to be called when an event is received for the given namespace
func (server *Server) On(namespace, eventName string, callback EventHandler) {
	ns, found := server.namespaces[namespace]
	if !found {
		ns = NewNamespace(namespace)
		server.namespaces[namespace] = ns
	}
	ns.eventHandlers[eventName] = callback
}

// Serve starts the Socket.io loop
func (server *Server) Serve() {
}

// ServeHTTP handles http requests
//
// implements http.Handler
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	log := server.logger.Child(nil, "servehttp")

	// 1. analyze the query to get the socket.io config
	query := r.URL.Query()

	// EIO (4), t (some short id), transport (polling)
	log.Infof("query: %+v", query)
	transports := query.Get("transport")
	transports = strings.TrimPrefix(strings.TrimSuffix(transports, "]"), "[")
	// TODO: parse transports for ","
	log.Debugf("Looking for transport: %s", transports)
	transport, found := server.transports[transports]
	if !found {
		log.Errorf("Transport %s not found, aborting", transports)
		core.RespondWithError(w, http.StatusBadRequest, errors.HTTPBadRequest.WithMessagef("Invalid Socket.io transport: %s", query.Get("transport")))
		return
	}
	log.Debugf("Using transport: %s", transport)

	socket := server.sockets.Get(query.Get("sid"))
	if socket == nil {
		log.Debugf("Creating a new socket")
		if socket, err = transport.Accept(w, r); err != nil {
			core.RespondWithError(w, http.StatusBadGateway, err)
			return
		}
		log.Debugf("Transport %s accepted request and created socket %s", transport, socket)
		server.sockets.Add(socket)

		// TODO: Start the socket loop
	} else {
		log.Debugf("using socket %s", socket)
	}

	log.Debugf("Socket %s will serve this request", socket)
	// TODO: if transport change, upgrade socket transport
	socket.ServeHTTP(w, r)

	// 2. upgrade connection to websocket
	// 3. start loop
}
