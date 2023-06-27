package main

import (
    "fmt"
    "github.com/Juniper/go-ez/ezrest"
)

func main() {
    
    username := "regress"
    password := "MaRtInI"
    hostname := "10.209.1.207:3000"
    format   := "json"
    rpc      := "get-interface-information"

    r := ezrest.NewRequest()
    r.BasicAuth(username, password)

    uri := fmt.Sprintf("http://%s/rpc/%s@format=%s?interface-name=fxp0", hostname, rpc, format)
    response := r.Send("get", uri, nil, nil, nil)

    if response.Error != nil {
       fmt.Println(response.Error) 
    } else {
       fmt.Println(string(response.Body))
    }
}
