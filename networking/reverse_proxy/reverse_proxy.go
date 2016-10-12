package main

/*
 * 反向代理（Reverse Proxy）方式是指以代理服务器来接受internet上的连接请求，然后
 * 将请求转发给内部网络上的服务器，并将从服务器上得到的结果返回给internet上请求
 * 连接的客户端，此时代理服务器对外就表现为一个反向代理服务器。
 */

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type handle struct {
	host string
	port string
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse("http://" + this.host + ":" + this.port)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func startServer() {
	//被代理的服务器host和port
	h := &handle{host: "127.0.0.1", port: "80"}
	err := http.ListenAndServe(":8888", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	startServer()
}
