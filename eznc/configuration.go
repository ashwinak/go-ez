package eznc

import (
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

// ExecConfigRPC is a func
func ExecConfigRPC(ns *netconf.Session, rpc string, attributes string, configuration string) (*netconf.RPCReply, error) {
	if rpc == "lock-configuration" || rpc == "unlock-configuration" {
		reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s/>", rpc)))
		return reply, err
	} else {
		if rpc == "commit-configuration" {
			if configuration != "" {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s>%s</%s>", rpc, configuration, rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			}
			reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s/>", rpc)))
			if err != nil {
				panic(err)
			}
			return reply, err

		} else if rpc == "load-configuration" {
			if attributes != "" && configuration != "" {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s %s>%s</%s>", rpc, attributes, configuration, rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			} else if attributes != "" && configuration == "" {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s %s></%s>", rpc, attributes, rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			} else if attributes == "" && configuration != "" {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s>%s</%s>", rpc, configuration, rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			} else {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s/>", rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			}
		} else if rpc == "get-configuration" {
			if attributes != "" {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s %s></%s>", rpc, attributes, rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			} else {
				reply, err := ns.Exec(netconf.RawMethod(fmt.Sprintf("<%s/>", rpc)))
				if err != nil {
					panic(err)
				}
				return reply, err
			}
		} else {
			err := "Invalid RPC request!!"
			panic(err)
		}
	}
}

func Get_configuration(ns *netconf.Session, attributes string) (*netconf.RPCReply, error) {
	if ns == nil {
		err := "Netconf connection handle missing!!"
		panic(err)
	} else {
		reply, err := ExecConfigRPC(ns, "get-configuration", attributes, "")
		return reply, err
	}
}

func Commit_configuration(ns *netconf.Session, configuration string) (*netconf.RPCReply, error) {
	if ns == nil {
		err := "Netconf connection handle missing!!"
		panic(err)
	} else {
		reply, err := ExecConfigRPC(ns, "commit-configuration", "", configuration)
		return reply, err
	}
}

func Load_configuration(ns *netconf.Session, attributes string, configuration string) (*netconf.RPCReply, error) {
	if ns == nil {
		err := "Netconf connection handle missing!!"
		panic(err)
	} else {
		reply, err := ExecConfigRPC(ns, "load-configuration", attributes, configuration)
		return reply, err
	}
}

func Lock_configuration(ns *netconf.Session) (*netconf.RPCReply, error) {
	if ns == nil {
		err := "Netconf connection handle missing!!"
		panic(err)
	} else {
		reply, err := ExecConfigRPC(ns, "lock-configuration", "", "")
		return reply, err
	}
}

func Unlock_configuration(ns *netconf.Session) (*netconf.RPCReply, error) {
	if ns == nil {
		err := "Netconf connection handle missing!!"
		panic(err)
	} else {
		reply, err := ExecConfigRPC(ns, "unlock-configuration", "", "")
		return reply, err
	}
}
