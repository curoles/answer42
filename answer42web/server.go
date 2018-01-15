package main

import (
    "fmt"
    "net/http"
)

func runServer(prgOptions *ProgramOptions) {
    http.HandleFunc("/", handleHttp)

    hostName := "localhost"
    Logger.Println("Start serving on port:", prgOptions.HttpPort)
    http.ListenAndServe(fmt.Sprintf("%s:%d", hostName, prgOptions.HttpPort), nil)
}


func handleHttp(
    resp http.ResponseWriter,
    req *http.Request) {

    fmt.Fprintf(resp, "Answer is 42")

}

