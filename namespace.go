package socketio

// ConnectHandler called when a socket connects
type ConnectHandler func(socket Socket)

// DisconnectHandler called when a socket disconnects
type DisconnectHandler func(socket Socket)

// ErrorHandler called when an error occurs on a socket
type ErrorHandler func(socket Socket, err error)

// EventHandler called when an event occurs on a socket
type EventHandler func(socket Socket, message string)

// namespace is a collection of handlers attached to a socket.io namespace
type Namespace struct {
	Name              string
	connectHandler    ConnectHandler
	disconnectHandler DisconnectHandler
	errorHandler      ErrorHandler
	eventHandlers     map[string]EventHandler
}

// NewNamespace creates a new Namespace
func NewNamespace(name string) *Namespace {
	return &Namespace{
		Name:          name,
		eventHandlers: map[string]EventHandler{},
	}
}

// EmitConnect emits a Connect event
func (namespace *Namespace) EmitConnect(socket Socket) {
	if namespace.connectHandler != nil {
		namespace.connectHandler(socket)
	}
}


// EmitDisconnect calls the DisconnectHandler if any
func (namespace *Namespace) EmitDisconnect(socket Socket) {
	if namespace.disconnectHandler != nil {
		namespace.disconnectHandler(socket)
	}
}


// EmitError calls the ErrorHandler if any
func (namespace *Namespace) EmitError(socket Socket, err error) {
	if namespace.errorHandler != nil {
		namespace.errorHandler(socket, err)
	}
}

// EmitEvent calls the EventHandler if any
func (namespace *Namespace) EmitEvent(eventname string, socket Socket, message string) {
	if handler, found := namespace.eventHandlers[eventname]; found {
		handler(socket, message)
	}
}
