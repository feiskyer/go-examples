package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// TODO: setup this
	accessKey = ""
	secretKey = ""
	server    = "https://xxxx"
)

type ServiceSpec struct {
	Name                string
	Algorithm           string
	Image               string
	WorkingDir          string
	ContainerSize       string
	SSLCert             string
	ServicePort         uint16
	ContainerPort       uint16
	Replicas            uint32
	HealthCheckInterval uint32
	HealthCheckFall     uint32
	HealthCheckRise     uint32
	Protocol            string
	TTY                 bool
	SessionAffinity     bool
	Command             []string
	Entrypoint          string
	Env                 []string
	Volumes             []string
	Labels              map[string]string
}

type ServiceUpdate struct {
	Replicas *int
	Image    *string
	FIP      *string
}

func encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if _, err := params.Write(buf); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func Request(method, uri string, postData interface{}) ([]byte, error) {
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	var postBody = &bytes.Buffer{}
	if postData != nil {
		postBody, err = encodeData(postData)
		if err != nil {
			log.Fatal(err)
		}
	}
	req, err := http.NewRequest(method, server+uri, postBody)
	if err != nil {
		return nil, err
	}

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Date", time.Unix(time.Now().Unix(), 0).Format("Mon, 2 Jan 2006 15:04:05 -0700"))
	if signature, err := makeSign(accessKey, secretKey, req); err != nil {
		return nil, err
	} else {
		req.Header.Set("Authorization", fmt.Sprintf(" HSC %s:%s", accessKey, signature))
	}

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("Error response from server: %s", string(data))
		return nil, err
	}

	return data, nil
}

func getContainers() {
	containers, err := Request("GET", "/containers/json?filters=%7B%22label%22:%7B%22app%3Dnginx%22:true%7D%7D&limit=0", nil)
	if err != nil {
		log.Fatalf("GetContainers failed: %v", err)
	}

	fmt.Println(string(containers))
}

func getServices() {
	services, err := Request("GET", "/services", nil)
	if err != nil {
		log.Fatalf("Get service list failed: %v", err)
	}

	fmt.Println(string(services))
}

func createService(name string) {
	serviceSpec := ServiceSpec{
		Name:  name,
		Image: "nginx",
		//ServiceSize:   "m1",
		//ContainerSize: "m1",
		ServicePort:   80,
		ContainerPort: 80,
		Replicas:      3,
		Protocol:      "http",
		Labels:        map[string]string{"app": "nginx"},
	}
	service, err := Request("POST", "/services/create", serviceSpec)
	if err != nil {
		log.Fatalf("Create service failed: %v", err)
	}

	fmt.Println(string(service))

}

func deleteService(name string) {
	service, err := Request("DELETE", "/services/"+name, nil)
	if err != nil {
		log.Fatalf("Delete service failed: %v", err)
	}

	fmt.Printf("Service %q delete success\n", service)
}

func updateService(name string) {
	count := 3
	serviceUpdate := ServiceUpdate{
		Replicas: &count,
	}
	service, err := Request("POST", fmt.Sprintf("/services/%s/update", name), serviceUpdate)
	if err != nil {
		log.Fatalf("Update service failed: %v", err)
	}

	fmt.Println(string(service))
}

func main() {
	getContainers()
	//createService("nginx")
	//updateService("nginx")
	//deleteService("nginx")
	//getServices()
}
