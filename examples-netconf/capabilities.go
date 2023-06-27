package main

import (
    "fmt"
    "github.com/Juniper/go-ez/eznc"
)

func main() {
    username := "regress"
    password := "MaRtInI"
    hostname := "10.209.16.149"

    s, err := eznc.Device(hostname, username, password) 
    if err != nil {
        panic(err)
    }

    defer s.Close()

    fmt.Printf("Server Capabilities: '%+v'\n", s.ServerCapabilities)
    fmt.Printf("Session Id: %d\n\n", s.SessionID)
}
