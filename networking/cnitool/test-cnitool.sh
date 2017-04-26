#!/bin/bash
set -x
set -o errexit
set -o nounset
set -o pipefail

CNITOOL_ROOT=$(dirname "${BASH_SOURCE}")
NETNS=test

cleanup() {
    sudo ip netns del $NETNS
}

trap "cleanup" EXIT SIGINT

cd $CNITOOL_ROOT
go build -o cnitool

sudo ip netns add $NETNS
sudo CNI_PATH=/opt/cni/bin NETCONFPATH=/etc/cni/net.d ./cnitool add netlist /var/run/netns/test
sudo CNI_PATH=/opt/cni/bin NETCONFPATH=/etc/cni/net.d ./cnitool del netlist /var/run/netns/test
