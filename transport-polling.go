package socketio

import (
	"net/http"

	"github.com/gildas/go-logger"
)

type PollingTransport struct {
	logger *logger.Logger
}

type PollingTransportOptions struct {
	Logger *logger.Logger
}

func NewPollingTransport(options *PollingTransportOptions) *PollingTransport {
	return &PollingTransport{
		logger: logger.CreateIfNil(options.Logger, "socketio").Child("transport", "transport"),
	}
}

func (transport *PollingTransport) Accept(w http.ResponseWriter, r *http.Request) (Socket, error) {
	log := transport.logger.Child(nil, "accept")
	query := r.URL.Query()
	jsonp := query.Get("j")
	supportsBinary := query.Get("b64") == "" && len(jsonp) == 0

	log.Debugf("Accepting request")
	return NewPollingSocket(&PollingSocketOptions{
		UseJSONP:       true,
		Jsonp:          jsonp,
		SupportsBinary: supportsBinary,
		Transport:      transport,
		Logger:         transport.logger,
	}), nil
}

func (Transport *PollingTransport) String() string {
	return "polling"
}