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
   defer s.Close()

   reply, err := eznc.Commit_configuration(s, "<check/>")
   if err != nil {
       panic(err)
   }
   fmt.Printf("Reply: %+v", reply)
}
