# Introduction
go client library for JUNOS Automation

# FEATURES

Junos GoEZ is designed to provide the same capabilities as a user would have on the Junos CLI, but in an environment built for automation tasks. These capabilities include, but are not limited to:

. Remote connectivity and management of Junos devices via NETCONF and REST API's

. Retrieve "operational" or "run-state" information

. Retrieve configuration information

. Make configuration changes in unstructured and structured ways (Supports 'json'/'xml'/'set'/'text' format)

# INSTALLATION

Installing from Git is also supported (OS must have git installed).
. go get github.com/Juniper/go-ez/eznc (Or)
. git clone https://github.com/Juniper/go-ez.git (Use Git or checkout with SVN using the web URL)

# HELLO WORLD

The following is a quick "hello, world" example to ensure that the software was installed correctly. This code will simply connect to a device and display the netconf capabilities of the device.

package main

import (

    "fmt"

    "github.com/Juniper/go-ez/eznc"
)

func main() {

    username := "jnpr"
    password := "jnpr123"
    hostname := "my_host_or_ipaddr"

    s, err := eznc.Device(hostname, username, password) 
    if err != nil {
        panic(err)
    }

    defer s.Close()

    fmt.Printf("Server Capabilities: '%+v'\n", s.ServerCapabilities)
    fmt.Printf("Session Id: %d\n\n", s.SessionID)
}

Example output for an VMX device:

$ go run capabilities.go
Server Capabilities: '[urn:ietf:params:netconf:base:1.0 urn:ietf:params:netconf:capability:candidate:1.0 urn:ietf:params:netconf:capability:confirmed-commit:1.0 urn:ietf:params:netconf:capability:validate:1.0 urn:ietf:params:netconf:capability:url:1.0?scheme=http,ftp,file urn:ietf:params:xml:ns:netconf:base:1.0 urn:ietf:params:xml:ns:netconf:capability:candidate:1.0 urn:ietf:params:xml:ns:netconf:capability:confirmed-commit:1.0 urn:ietf:params:xml:ns:netconf:capability:validate:1.0 urn:ietf:params:xml:ns:netconf:capability:url:1.0?protocol=http,ftp,file http://xml.juniper.net/netconf/junos/1.0 http://xml.juniper.net/dmi/system/1.0]'

Session Id: 11969

# LICENSE

BSD 2-Clause License
