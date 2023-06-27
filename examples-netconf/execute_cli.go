package main

import (
    "fmt"
    "github.com/Juniper/go-ez/eznc"
)

func main() {
   s, err := eznc.Device("10.209.16.149", "regress", "MaRtInI") 
   if err != nil {
      panic(err)
   }
  
   // eznc.Cli(<connection_handle>, "operational command", "format") 
   reply, err := eznc.Cli(s, "show version", "json")
   if err != nil {
       panic(err)
   }

   defer s.Close()
   fmt.Printf("Reply: %+v", reply)
}
