package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "github.com/braintree/manners"
)

// Type associated with the handler.
// Implements http.Handler interface.
type httpHandler struct {
}

// Creates new instance of httpHandler and returns pointer to it.
func newHttpHandler() *httpHandler {
    return &httpHandler{}
}

// Handles HTTP request, sends back a response.
func (h *httpHandler) ServeHTTP(
    rsp http.ResponseWriter,
    req *http.Request) {

    query := req.URL.Query()
    name := query.Get("name")
    if name == "" {
        name = "dear guest"
    }

    fmt.Fprintf(rsp, "Answer is 42, %s", name)
}

// Waits for OS shutdown signal and stops HTTP server.
func listenForShutdown(ch <-chan os.Signal) {
    <-ch
    Logger.Println("Received shutdown signal from OS, stopping HTTP server...")
    Logger.Println("Stop accepting new connections, shut down after all the current requests are completed.")
    manners.Close()
}

// Starts HTTP server and listens to OS signals.
//
// Since function ListenAndServe blocks execution, to monitor
// OS signals we use go-routine.
//
func runServer(prgOptions *ProgramOptions) {
    httpHandler := newHttpHandler()

    // Set up monitoring OS signals Interrupt and Kill.
    sigch := make(chan os.Signal)
    signal.Notify(sigch, os.Interrupt, os.Kill)
    go listenForShutdown(sigch) // separate thread listening to the OS signals

    hostName := "" //"localhost"
    hostAndPort := fmt.Sprintf("%s:%d", hostName, prgOptions.HttpPort)
    Logger.Println("Start serving on:", hostAndPort)
    manners.ListenAndServe(hostAndPort, httpHandler)
}
