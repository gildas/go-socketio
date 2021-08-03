package socketio

import (
	"fmt"
	"net/http"
)

// Transport describes the transport mechanism used by a session
type Transport interface {
	Accept(http.ResponseWriter, *http.Request) (Socket, error)
	fmt.Stringer
}