package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "path"
    "regexp"
    "github.com/braintree/manners"
)

// Type associated with HTTP request handler function.
// Implements http.Handler interface.
// Also has map{key,handler} to find handler by URL path.
type httpHandler struct {
    handlers map[string]http.HandlerFunc
    rexpHandlers map[string]http.HandlerFunc
    rexpCache map[string]*regexp.Regexp
}

// Creates new instance of httpHandler and returns pointer to it.
func newHttpHandler() *httpHandler {
    return &httpHandler{
        handlers     : make(map[string]http.HandlerFunc),
        rexpHandlers : make(map[string]http.HandlerFunc),
        rexpCache    : make(map[string]*regexp.Regexp),
    }
}

// Add new pair {path, handler} to internal handler lookup map.
func (h *httpHandler) Add(path string, handler http.HandlerFunc, isRegexp bool) {
    if isRegexp == true {
        h.rexpHandlers[path] = handler
        compiledRexp, _ := regexp.Compile(path)
        h.rexpCache[path] = compiledRexp
    } else {
        h.handlers[path] = handler
    }
}

// Handles HTTP request, sends back a response.
func (h *httpHandler) ServeHTTP(
    rsp http.ResponseWriter,
    req *http.Request) {

    // Construct string "GET/POST URL"
    methodAndPath := req.Method + " " + req.URL.Path

    // Regexp match
    for pattern, handlerFunc := range h.rexpHandlers {
        if h.rexpCache[pattern].MatchString(methodAndPath) == true {
            handlerFunc(rsp, req)
            return
        }
    }

    // Simple "path" match
    for pattern, handlerFunc := range h.handlers {
        if mismatch, err := path.Match(pattern, methodAndPath); mismatch && err == nil {
            handlerFunc(rsp, req)
            return
        } else if err != nil {
            fmt.Fprint(rsp, err)
        }
    }

    http.NotFound(rsp, req)
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
    httpHandler.Add("GET /*", mainHttpHandler, /*regexp:*/false)

    // Set up monitoring OS signals Interrupt and Kill.
    sigch := make(chan os.Signal)
    signal.Notify(sigch, os.Interrupt, os.Kill)
    go listenForShutdown(sigch) // separate thread listening to the OS signals

    hostName := "" //"localhost"
    hostAndPort := fmt.Sprintf("%s:%d", hostName, prgOptions.HttpPort)
    Logger.Println("Start serving on:", hostAndPort)
    manners.ListenAndServe(hostAndPort, httpHandler)
}


// ?.
func mainHttpHandler(
    rsp http.ResponseWriter,
    req *http.Request) {

    Logger.Println("HTTP request:", req.Method + " " + req.URL.RequestURI())

    query := req.URL.Query()
    name := query.Get("name")
    if name == "" {
        name = "dear guest"
    }

    fmt.Fprintf(rsp, "Answer is 42, %s", name)
}
