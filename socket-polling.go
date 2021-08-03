package socketio

import (
	"fmt"
	"net/http"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

type PollingSocket struct {
	SocketID       string `json:"sid"`
	Namespace      string `json:"namespace"`
	useJSONP       bool
	jsonp          string
	supportsBinary bool
	transport      Transport
	logger         *logger.Logger
}

type PollingSocketOptions struct {
	Namespace      string `json:"namespace"`
	UseJSONP       bool
	Jsonp          string
	SupportsBinary bool
	Transport      Transport
	Logger         *logger.Logger
}

func NewPollingSocket(options *PollingSocketOptions) *PollingSocket {
	if options == nil {
		panic(errors.ArgumentMissing.With("options"))
	}
	socket := &PollingSocket{
		Namespace:      options.Namespace,
		transport:      options.Transport,
		useJSONP:       options.UseJSONP,
		jsonp:          options.Jsonp,
		supportsBinary: options.SupportsBinary,
		logger:         logger.CreateIfNil(options.Logger, "SOCKETIO").Child("socket", "socket", "type", "polling"),
	}
	// TODO: If client.protocol (EIO?) == 3 use client.ID
	if len(options.Namespace) == 0 {
		socket.SocketID = base64ID()
	} else {
		socket.SocketID = fmt.Sprintf("%s/%s", options.Namespace, base64ID())
	}
	return socket
}

func (socket *PollingSocket) ID() string {
	return socket.SocketID
}

func (socket *PollingSocket) String() string {
	return socket.SocketID
}

func (socket *PollingSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := socket.logger.Child(nil, "servehttp")
	query  := r.URL.Query()
	jsonp  := query.Get("j")
	origin := r.Header.Get("Origin")

	log.Debugf("Query: %+#v", query)
	// TODO: provide a way to customize. So the lib user can implement their own origin check
	if len(origin) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	log.Debugf("Processing %s", r.Method)
	switch r.Method {
	case http.MethodOptions:
		if len(jsonp) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
		}
	case http.MethodGet:
		if len(jsonp) > 0 {
			w.Header().Set("Content-Type", "text/javascript; charset=UTF-8")
			value := ""
			_, _ = w.Write([]byte(fmt.Sprintf("___eio[%s](\"%s\");", jsonp, value)))
		} else {
			if socket.supportsBinary {
				w.Header().Set("Content-Type", "application/octet-stream")
			} else {
				w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			}
			value := ""
			_, _ = w.Write([]byte(value))
		}
	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")
		// TODO: read stuff in, if supportBinary, read binary if provided
		log.Debugf("Content Type: %s", contentType)
		_, _ = w.Write([]byte("ok"))
	default:
		core.RespondWithError(w, http.StatusMethodNotAllowed, errors.HTTPMethodNotAllowed.WithStack())
	}
}
