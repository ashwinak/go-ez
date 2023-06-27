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

   //eznc.Load_configuration(<connection_handle>, attributes, configuration)

   configuration := `
     <configuration>
       <system>
         <syslog>
           <file><name>sample</name><contents><name>any</name><any/></contents></file>
         </syslog>
       </system>
     </configuration>
   `
   reply, err := eznc.Load_configuration(s, "", configuration)
   if err != nil {
       panic(err)
   }
   fmt.Printf("Reply: %+v", reply)
}
