package main

import (
    "fmt"
    "github.com/Juniper/go-ez/eznc"
    "github.com/Juniper/go-netconf/netconf"
)

func main() {
    username := "regress"
    password := "MaRtInI"
    host     := "10.209.16.149"

    //eznc.Device(host, username, password)
    s, err := eznc.Device(host, username, password) 
    if err != nil {
        panic(err)
    }

    defer s.Close()

    //reply, err := s.Exec(netconf.MethodGetConfig("running"))
    reply, err := s.Exec(netconf.MethodGetConfig("candidate"))
    if err != nil {
        panic(err)
    }
    fmt.Printf("Reply: %+v", reply)
}
