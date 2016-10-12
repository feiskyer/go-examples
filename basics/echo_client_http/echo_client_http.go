package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/tmp/echo.sock")
}

func reader(client *http.Client) {
	for {
		resp, err := client.Get("http://justatest.com")
		if err != nil {
			panic(err.Error())
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Client got:", string(body))
	}
}

func main() {
	tr := &http.Transport{Dial: fakeDial}
	client := &http.Client{Transport: tr}

	go reader(client)

	for {
		resp, err := client.Post("http://justatest.com",
			"application/x-www-form-urlencoded",
			strings.NewReader("name=cjb"))
		if err != nil {
			println(err.Error())
			break
		}
		defer resp.Body.Close()
		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//      println(err.Error())
		//      break
		//}
		//fmt.Println("Client resp: ", string(body))
		time.Sleep(1e9)
	}
}
