package main

import (
	"encoding/xml"
	"flag"
	"os"

	libvirtgo "github.com/alexzorin/libvirt-go"
	"github.com/golang/glog"
)

/*
#cgo LDFLAGS: -lvirt
#include <libvirt/libvirt.h>
#include <libvirt/virterror.h>
#include <stdlib.h>
*/
import "C"

const (
	VIR_DOMAIN_MEMORY_STAT_SWAP_IN        = C.VIR_DOMAIN_MEMORY_STAT_SWAP_IN
	VIR_DOMAIN_MEMORY_STAT_SWAP_OUT       = C.VIR_DOMAIN_MEMORY_STAT_SWAP_OUT
	VIR_DOMAIN_MEMORY_STAT_MAJOR_FAULT    = C.VIR_DOMAIN_MEMORY_STAT_MAJOR_FAULT
	VIR_DOMAIN_MEMORY_STAT_MINOR_FAULT    = C.VIR_DOMAIN_MEMORY_STAT_MINOR_FAULT
	VIR_DOMAIN_MEMORY_STAT_UNUSED         = C.VIR_DOMAIN_MEMORY_STAT_UNUSED
	VIR_DOMAIN_MEMORY_STAT_AVAILABLE      = C.VIR_DOMAIN_MEMORY_STAT_AVAILABLE
	VIR_DOMAIN_MEMORY_STAT_ACTUAL_BALLOON = C.VIR_DOMAIN_MEMORY_STAT_ACTUAL_BALLOON
	VIR_DOMAIN_MEMORY_STAT_RSS            = C.VIR_DOMAIN_MEMORY_STAT_RSS
	VIR_DOMAIN_MEMORY_STAT_NR             = C.VIR_DOMAIN_MEMORY_STAT_NR

	VIR_DOMAIN_CPU_STATS_CPUTIME    = C.VIR_DOMAIN_CPU_STATS_CPUTIME
	VIR_DOMAIN_CPU_STATS_USERTIME   = C.VIR_DOMAIN_CPU_STATS_USERTIME
	VIR_DOMAIN_CPU_STATS_SYSTEMTIME = C.VIR_DOMAIN_CPU_STATS_SYSTEMTIME

	libvirtdAddress = "qemu:///system"
)

type VirDomain struct {
	Name    string    `xml:"name"`
	UUID    string    `xml:"uuid"`
	Memory  string    `xml:"memory"`
	Devices VirDevice `xml:"devices"`
}

type VirDevice struct {
	Disks      []VirDisk      `xml:"disk"`
	Interfaces []VirInterface `xml:"interface"`
}

// TODO: filesystem, rbd
type VirDisk struct {
	Type   string        `xml:"type,attr"`
	Source VirDiskSource `xml:"source"`
	Target VirDiskTarget `xml:"target"`
}
type VirDiskSource struct {
	File string `xml:"file,attr"`
}
type VirDiskTarget struct {
	Dev string `xml:"dev,attr"`
}

type VirInterface struct {
	Type   string             `xml:"type,attr"`
	Device VirInterfaceTarget `xml:"target"`
	Mac    VirInterfaceMac    `xml:"mac"`
}
type VirInterfaceTarget struct {
	Dev string `xml:"dev,attr"`
}
type VirInterfaceMac struct {
	Address string `xml:"mac,attr"`
}

func GetCPUStats(domain libvirtgo.VirDomain) error {
	// Get the number of cpus available to query from the host perspective,
	ncpus, err := domain.GetCPUStats(nil, 0, 0, 0, 0)
	if err != nil {
		return err
	}

	// Get how many statistics are available for the given @start_cpu.
	nparams, err := domain.GetCPUStats(nil, 0, 0, 1, 0)
	if err != nil {
		return err
	}

	// Query per-cpu stats
	var perCPUStats libvirtgo.VirTypedParameters
	_, err = domain.GetCPUStats(&perCPUStats, nparams, 0, uint32(ncpus), 0)
	if err != nil {
		return err
	}

	glog.Infof("Get per-cpu stats: %v", perCPUStats)

	// Query total stats
	var cpuStats libvirtgo.VirTypedParameters
	nparams, err = domain.GetCPUStats(nil, 0, -1, 1, 0)
	if err != nil {
		return err
	}
	_, err = domain.GetCPUStats(&cpuStats, nparams, -1, 1, 0)
	if err != nil {
		return err
	}
	glog.Infof("Get total cpu stats: %v", cpuStats)

	return nil
}

func GetMemoryStats(domain libvirtgo.VirDomain) error {
	memStats, err := domain.MemoryStats(VIR_DOMAIN_MEMORY_STAT_NR, 0)
	if err != nil {
		return err
	}

	for _, stat := range memStats {
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_SWAP_IN {
			glog.Infof("swap_in %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_SWAP_OUT {
			glog.Infof("swap_out %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_MAJOR_FAULT {
			glog.Infof("major_fault %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_MINOR_FAULT {
			glog.Infof("minor_fault %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_UNUSED {
			glog.Infof("unused %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_AVAILABLE {
			glog.Infof("available %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_ACTUAL_BALLOON {
			glog.Infof("actual %v", stat.Val)
		}
		if stat.Tag == VIR_DOMAIN_MEMORY_STAT_RSS {
			glog.Infof("rss %v", stat.Val)
		}
	}

	return nil
}

/*
type VirDomainInterfaceStats struct {
	RxBytes   int64
	RxPackets int64
	RxErrs    int64
	RxDrop    int64
	TxBytes   int64
	TxPackets int64
	TxErrs    int64
	TxDrop    int64
}*/
func GetNetworkStats(domain libvirtgo.VirDomain, virDomain *VirDomain) error {
	for _, iface := range virDomain.Devices.Interfaces {
		ifaceStats, err := domain.InterfaceStats(iface.Device.Dev)
		if err != nil {
			return err
		}
		glog.Infof("Get iface %s stat: %v", iface.Device.Dev, ifaceStats)
	}

	return nil
}

func GetBlockStats(domain libvirtgo.VirDomain, virDomain *VirDomain) error {
	for _, blk := range virDomain.Devices.Disks {
		blkStats, err := domain.BlockStats(blk.Target.Dev)
		if err != nil {
			return err
		}
		glog.Infof("Get blk %s stat: %v", blk.Target.Dev, blkStats)
	}
	return nil
}

func main() {
	vmName := flag.String("name", "vm1", "the name of vm")
	flag.Parse()
	flag.Set("logtostderr", "true")

	var conn libvirtgo.VirConnection

	conn, err := libvirtgo.NewVirConnection(libvirtdAddress)
	if err != nil {
		glog.Error("fail to connect to libvirtd ", libvirtdAddress, err)
		os.Exit(1)
	}

	domain, err := conn.LookupDomainByName(*vmName)
	if err != nil {
		glog.Errorf("Can't find domain %s", *vmName)
		os.Exit(1)
	}

	p, err := doamin.IsPersistent()
	if err != nil {
		glog.Errorf("Can not get persistent: %v", err)
		os.Exit(1)
	}

	glog.Infof("Persistenst is %v\n", p)

	state, err := domain.GetState()
	if err != nil {
		glog.Errorf("Can't get state for domain %s", *vmName)
	}
	if state[0] != libvirtgo.VIR_DOMAIN_RUNNING {
		glog.Errorf("Domain %s is not running", *vmName)
		os.Exit(1)
	}

	glog.Infof("Got running domain %s", *vmName)

	xmlDesc, err := domain.GetXMLDesc(0)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	glog.V(4).Infof("XML description for domain is %s", xmlDesc)

	var virDomain VirDomain
	err = xml.Unmarshal([]byte(xmlDesc), &virDomain)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}

	glog.V(4).Infof("Get domain description: %v", virDomain)

	err = GetCPUStats(domain)
	if err != nil {
		glog.Infof("GetCPUStats error: %v", err)
	}
	err = GetMemoryStats(domain)
	if err != nil {
		glog.Infof("GetMemoryStats error: %v", err)
	}
	err = GetNetworkStats(domain, &virDomain)
	if err != nil {
		glog.Infof("GetNetworkStats error: %v", err)
	}
	err = GetBlockStats(domain, &virDomain)
	if err != nil {
		glog.Infof("GetBlockStats error: %v", err)
	}
}
