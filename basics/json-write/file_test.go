package main

import "fmt"
import "os"
import "encoding/json"
import "io/ioutil"

const testSpec = `{"containers":[{"image":"nginx","name":"kube_993cae13-552a-11e5-83a4-000c29f6dfc5_nginx_ns1_nginx.de0022e5_26e555cc","ports":[{"containerPort":80,"protocol":"TCP"}],"tty":false,"volumes":[{"path":"/var/run/secrets/kubernetes.io/serviceaccount","readOnly":true,"volume":"default-token-yvcun"}]}],"id":"kube_993cae13-552a-11e5-83a4-000c29f6dfc5_nginx_ns1","resource":{"memory":192,"vcpu":1},"tty":true,"type":"pod","volumes":[{"driver":"vfs","name":"default-token-yvcun","source":"/var/lib/kubelet/pods/993cae13-552a-11e5-83a4-000c29f6dfc5/volumes/kubernetes.io~secret/default-token-yvcun"}]}`

func main() {
	h, e := os.Hostname()
	if e != nil {
		fmt.Println("Error %s", e)
	}

	fmt.Println("Hostname%s", h)

	test := make(map[string]interface{})
	err := json.Unmarshal([]byte(testSpec), &test)
	if err != nil {
		fmt.Println("Error", err)
	}

	err = ioutil.WriteFile("aaa.json", []byte(testSpec), 0644)
	if err != nil {
		fmt.Println("Error", err)
	}
}
