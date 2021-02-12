package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-socketio"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	var (
		port = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the TCP port for which the server listens to")
	)
	flag.Parse()

	log := logger.Create("SERVER", &logger.StdoutStream{ Unbuffered: true })
	defer log.Flush()

	ioserver, err := socketio.NewServer()
	if err != nil {
		log.Fatalf("Failed to create the Socket IO server")
		os.Exit(1)
	}

	ioserver.OnConnect("/", func (socket socketio.Socket) {
		log.Infof("Connected to socket: %s", socket)
	})

	ioserver.OnDisconnect("/", func (socket socketio.Socket) {
		log.Infof("Disconnected from socket: %s", socket)
	})

	ioserver.OnError("/", func (socket socketio.Socket, err error) {
		log.Errorf("Received error...", err)
	})

	ioserver.On("/", "message", func (socket socketio.Socket, message string) {
		log.Infof("Received: %s", message)
		socket.emit("response", "Received: %s", message)
	})

	go ioserver.Serve()
	defer ioserver.Close()

	server := &http.Server{
		Addr:     fmt.Sprintf("0.0.0.0:%d", *port),
		ErrorLog: log.AsStandardLog(logger.ERROR),
	}

	http.Handle("/socket.io/", ioserver)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Infof("Starting Web Server on port %d", *port)
	if err := server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			log.Fatalf("Failed to start the WEB server on port: %d", *port, err)
			fmt.Fprintf(os.Stderr, "Failed to start the WEB server on port: %d. Error: %s", *port, err)
			os.Exit(1)
		}
	}
}