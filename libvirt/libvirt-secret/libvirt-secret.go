package main

import (
	"flag"
	"os"

	libvirtgo "github.com/alexzorin/libvirt-go"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	var conn libvirtgo.VirConnection
	var libvirtdAddress = "qemu:///system"

	conn, err := libvirtgo.NewVirConnection(libvirtdAddress)
	if err != nil {
		glog.Error("fail to connect to libvirtd ", libvirtdAddress, err)
		os.Exit(1)
	}

	sec, err := conn.LookupSecretByUsage(libvirtgo.VIR_SECRET_USAGE_TYPE_CEPH, "client.cinder secret")
	if err != nil {
		glog.Errorf("%v", err)
	}

	uuid, err := sec.GetUUIDString()
	if err != nil {
		glog.Errorf("%v", err)
	}
	glog.Infof("Secret: %v", uuid)

	err = conn.SecretSetValue(uuid, "AA+QiU19VWaju5PSi/Vgz4g==")
	if err != nil {
		glog.Errorf("%v", err)
	}
}
