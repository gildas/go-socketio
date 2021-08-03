package socketio

import (
	"fmt"
	"net/http"
)

type Socket interface {
	ID() string
	http.Handler
	fmt.Stringer
}