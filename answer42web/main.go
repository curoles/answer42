package main

import (
    "fmt"
    "github.com/curoles/go-fun/utils/filelog"
    "log"
    "flag"
    "os"
    "encoding/json"
)

var Logger *log.Logger = nil

type ProgramOptions struct {
    LogFile string
    ConfigFile string
    HttpPort int
}

func createProgramOptions(options *ProgramOptions) {
    flag.StringVar(&options.ConfigFile, "config-file", options.ConfigFile, "Configuration file")
    flag.StringVar(&options.LogFile, "log-file", options.LogFile, "Log file, envar $A42W_LOGFILE")
    flag.IntVar(&options.HttpPort, "http-port", options.HttpPort, "Web server HTTP port, envar $A42W_PORT")
}

func readConfigFile(options *ProgramOptions) {
    if file, err := os.Open(options.ConfigFile); err == nil {
        defer file.Close()
        fmt.Println("Reading configuration file ", options.ConfigFile)
        decoder := json.NewDecoder(file)
        if err := decoder.Decode(options); err != nil {
            fmt.Println("Error reading configuration file:", err)
        }
    }
}

func readEnv(options *ProgramOptions) {
    if valStr, present := os.LookupEnv("A42W_LOGFILE"); present == true {
        options.LogFile = valStr
    }
}

func main() {

    prgOptions := ProgramOptions{
        LogFile : "./answer42web.log",
        ConfigFile : "./answer42web-config.json",
        HttpPort : 4000}

    readEnv(&prgOptions)
    readConfigFile(&prgOptions)
    createProgramOptions(&prgOptions)
    flag.Parse()


    fl, err := filelog.New(prgOptions.LogFile, "", log.LstdFlags, true)
    if err != nil {
        panic(fmt.Sprintf("can't create log file \"%s\", error: %v", prgOptions.LogFile, err))
    }
    defer fl.Close()

    Logger = fl.Logger

    if cwdStr, err := os.Getwd(); err == nil {
        Logger.Println("Current working directory:", cwdStr)
    }
    Logger.Println("Log file:", prgOptions.LogFile)

    runServer(&prgOptions)
}



