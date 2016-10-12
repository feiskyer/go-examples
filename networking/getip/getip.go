package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func GetAddr() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return "Erorr"
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func GetExternalAddr() {
	resp, err := http.Get("http://ifconfig.me/ip")
	if err != nil {
		fmt.Println("Can not connect to network")
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Can not get result")
		return
	}

	fmt.Println("External ip addr: ", string(result))
}

func main() {
	addr := GetAddr()
	fmt.Println("Local ip addr: ", addr)
	GetExternalAddr()
}
