package main

import (
    "fmt"
    "github.com/Juniper/go-ez/ezrest"
)

func main() {
    
    username := "regress"
    password := "MaRtInI"
    hostname := "10.209.1.207:3000"

    rpc      := `
      <lock-configuration/>
      <load-configuration>
        <configuration>
          <system>
            <syslog>
              <file><name>sample</name><contents><name>any</name><any/></contents></file>
            </syslog>
          </system>
        </configuration>
      </load-configuration>
      <commit-configuration/>
      <unlock-configuration/>
    `

    r := ezrest.NewRequest()
    r.BasicAuth(username, password)

    // Running multiple RPC's with single POST request 
    xmlBody := fmt.Sprintf(rpc, hostname, username, password)
    uri := fmt.Sprintf("http://%s/rpc?stop-on-error=1", hostname) 

    headers := map[string]string{
        "Content-Type": "application/xml",
        "Accept": "application/xml",  
    }
    response := r.Send("post", uri, []byte(xmlBody), headers, nil)

    if response.Error != nil {
       fmt.Println(response.Error) 
    } else {
       fmt.Println(string(response.Body))
    }
}
