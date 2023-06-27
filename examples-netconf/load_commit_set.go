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

   // lock configuration
   reply, err := eznc.Lock_configuration(s)
   if err != nil {
       panic(err)
   } else {
       fmt.Printf("Reply: %+v", reply)
       // File containing set format configuration
       url := "/var/tmp/set.config"
       attributes := fmt.Sprintf("url=\"%s\" action =\"%s\" format=\"%s\"", url, "set", "text")
       //eznc.Load_configuration(<connection_handle>, attributes, configuration) 
       reply, err := eznc.Load_configuration(s, attributes, "")
       if err != nil {
          //Unlock configuration
          reply, err := eznc.Unlock_configuration(s)
          if err != nil {
             panic(err)
          }
          fmt.Printf("Reply: %+v", reply)
       } else {
          fmt.Printf("Reply: %+v", reply)
          // Commit configuration
          reply, err := eznc.Commit_configuration(s, "")
          if err != nil {
             panic(err)
          } 
          fmt.Printf("Reply: %+v", reply)
          //Unlock configuration session
          unlock_reply, err := eznc.Unlock_configuration(s)
          if err != nil {
             panic(err)
          }
          fmt.Printf("Reply: %+v", unlock_reply)   
       }
   }
}
