package main

import (
	"encoding/json"
	"fmt"
)

type RBDVolume struct {
	Keyring     string   `json:"keyring"`
	Authenabled bool     `json:"auth_enabled"`
	AuthUser    string   `json:"auth_username"`
	Hosts       []string `json:"hosts"`
	Ports       []int    `json:"ports"`
	Name        string   `json:"name"`
	AccessMode  string   `json:"access_mode"`
	VolumeType  string   `json:"volume_type"`

	qosSpecs   string `json:"qos_specs",omitempty`
	secretUUID string `json:"secret_uuid",omitempty`
	secretType string `json:"secret_type",omitempty`
}

func main() {
	data := map[string]interface{}{
		"qos_specs":     "",
		"ports":         []int{6789},
		"volume_type":   "rbd",
		"name":          "cinder/volume-d9dbb2f4-e7e3-450d-8fd4-2ed6d6052eb2",
		"secret_uuid":   "3bb14c9f-a468-450e-8a15-727f9eb90f2f",
		"hosts":         []string{"172.31.6.228"},
		"auth_enabled":  true,
		"access_mode":   "rw",
		"auth_username": "cinder",
		"keyring":       "AQAtqv9V3u4nKRAA9xfic687DqPW1FV/rly3nw==",
		"secret_type":   "ceph",
	}
	fmt.Println(data)

	data2, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data2))

	var volume RBDVolume
	err = json.Unmarshal(data2, &volume)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(volume)
}
