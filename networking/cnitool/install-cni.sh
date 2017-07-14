#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

git clone https://github.com/containernetworking/plugins $GOPATH/src/github.com/containernetworking/plugins
cd $GOPATH/src/github.com/containernetworking/plugins

# build cni plugins and copy to specified folder.
./build.sh
sudo mkdir -p /opt/cni/bin
sudo cp bin/* /opt/cni/bin/

# Setup network config.
sudo mkdir -p /etc/cni/net.d
sudo sh -c 'cat >/etc/cni/net.d/10-mynet.conf <<-EOF
{
    "cniVersion": "0.3.0",
    "name": "mynet",
    "type": "bridge",
    "bridge": "cni0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "subnet": "10.100.0.0/16",
        "routes": [
            { "dst": "0.0.0.0/0"  }
        ]
    }
}
EOF'

sudo sh -c 'cat >/etc/cni/net.d/netlist.conflist <<-EOF
{
    "cniVersion": "0.3.1",
    "name": "netlist",
    "plugins": [
        {
            "type": "bridge",
            "bridge": "cni1",
            "isGateway": true,
            "ipMasq": true,
            "ipam": {
                "type": "host-local",
                "subnet": "10.200.0.0/16",
                "routes": [
                    { "dst": "0.0.0.0/0"  }
                ]
            }
        }
    ]
}
EOF'

sudo sh -c 'cat >/etc/cni/net.d/99-loopback.conf <<-EOF
{
    "cniVersion": "0.3.0",
    "type": "loopback"
}
EOF'

