package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GetFloatingIPs(container string) (fips []string, err error) {
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", server+"/fips?filters={\"container\":"+container+"}}", nil)
	if err != nil {
		return
	}

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
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
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("Error response from server: %s", string(data))
		return
	}

	err = json.Unmarshal(data, &fips)
	if err != nil {
		return
	}

	return fips, err
}

func AllocateFloatingIPs() (fips []string, err error) {
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", server+"/fips/allocate?count=2", nil)
	if err != nil {
		return
	}

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
	req.Header.Set("Date", time.Unix(time.Now().Unix(), 0).Format("Mon, 2 Jan 2006 15:04:05 -0700"))
	req.Header.Set("Content-Type", "application/json")
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
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("Error response from server: %s", string(data))
		return
	}

	err = json.Unmarshal(data, &fips)
	if err != nil {
		return
	}

	return fips, err
}
