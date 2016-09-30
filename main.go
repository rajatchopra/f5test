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
		VxlanGateway: "10.130.0.5/14",
		InternalAddress: "10.3.89.213",
		SetupOSDNVxLAN: true,
	}
	p, err := f5.NewF5Plugin(f5cfg)
	nodes := [2]string{"10.3.89.172", "10.3.89.173"}
	for _, ipStr := range nodes {
		err := p.F5Client.AddVtep(ipStr)
		if err != nil {
			fmt.Printf("Error adding %s - %v\n", ipStr, err)
		}
	}
	fmt.Printf("Testing f5 - %v, %v\n", p, err)
}
