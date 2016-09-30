package main

import (
	"fmt"
	"flag"
	"github.com/rajatchopra/f5test/f5"
)

func main() {
	flag.Parse()
	f5cfg := f5.F5PluginConfig{
		Host: "10.3.88.146",
		Username: "admin",
		Password: "openshift",
		HttpVserver: "ose-vserver",
		HttpsVserver: "https-ose-vserver",
		PrivateKey: "/root/.ssh/libra.pem",
		Insecure: true,
		PartitionPath: "",
		VxlanGateway: "10.130.0.5/16",
		SetupOSDNVxLAN: true,
	}
	p, err := f5.NewF5Plugin(f5cfg)
	fmt.Printf("Testing f5 - %v, %v", p, err)
}
