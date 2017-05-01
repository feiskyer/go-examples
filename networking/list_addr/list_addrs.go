package main

import (
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink"
)

func main() {
	iface, err := netlink.LinkByName("eth0")
	if err != nil {
		fmt.Println(err)
		return
	}

	addrs, err := netlink.AddrList(iface, syscall.AF_INET)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, addr := range addrs {
		fmt.Printf("%#v\n", addr)
	}
}
