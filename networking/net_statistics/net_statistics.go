package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

var netStats = []string{"tx_packets", "tx_errors", "tx_dropped", "rx_bytes", "rx_packets", "rx_errors", "rx_dropped"}

func getNetworkInterfaceStats(interfaceName string) (map[string]uint64, error) {
	out := make(map[string]uint64)
	for _, netStat := range netStats {
		data, err := readSysfsNetworkStats(interfaceName, netStat)
		if err != nil {
			return nil, err
		}
		out[netStat] = data
	}

	return out, nil
}

// Reads the specified statistics available under /sys/class/net/<EthInterface>/statistics
func readSysfsNetworkStats(ethInterface, statsFile string) (uint64, error) {
	data, err := ioutil.ReadFile(filepath.Join("/sys/class/net", ethInterface, "statistics", statsFile))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

func main() {
	interfaceName := flag.String("eth", "eth0", "the name of interface")
	flag.Parse()

	stats, err := getNetworkInterfaceStats(*interfaceName)
	if err != nil {
		fmt.Printf("getNetworkInterfaceStats failed: %v\n", interfaceName, err)
	}

	fmt.Printf("Get network stats for %v: %v\n", *interfaceName, stats)
}
