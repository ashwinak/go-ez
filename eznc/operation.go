package eznc 

import (
    "github.com/Juniper/go-netconf/netconf"
    "fmt"
)

// Get's connection handle to device
func Device(host string, username string, password string) (*netconf.Session, error) {
    s, err := netconf.DialSSH(host,
        netconf.SSHConfigPassword(username, password))
    if err != nil {
        panic(err)
    }
    return s, err
}

// Get's RPC execution reply
func Rpc(s *netconf.Session, rpc string) (*netconf.RPCReply, error) {
    concat_string := fmt.Sprintf("<%s/>", rpc)
    reply, err := s.Exec(netconf.RawMethod(concat_string))
    if err != nil {
       panic(err)
    }
    return reply, err
}

// Get's operational command output
func Cli(s *netconf.Session, rpc string, format string) (*netconf.RPCReply, error) {
    concat_string := fmt.Sprintf("<command format=\"%s\">%s</command>", format, rpc)
    reply, err := s.Exec(netconf.RawMethod(concat_string))
    if err != nil {
       //panic(err)
    }
    return reply, err
}
