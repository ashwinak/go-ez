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

   //reply, err := eznc.Get_configuration(s, "")
   reply, err := eznc.Get_configuration(s, "format=\"json\" database=\"candidate\"")
   //reply, err := eznc.Get_configuration(s, "format='xml' database='candidate'")
   if err != nil {
       panic(err)
   }
   fmt.Printf("Reply: %+v", reply)
}
