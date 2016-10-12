package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	HYPER_PROTO       = "unix"
	HYPER_ADDR        = "/var/run/hyper.sock"
	HYPER_SCHEME      = "http"
	HYPER_MINVERSION  = "0.4.0"
	DEFAULT_IMAGE_TAG = "latest"

	KEY_ID             = "id"
	KEY_IMAGEID        = "imageId"
	KEY_IMAGENAME      = "imageName"
	KEY_ITEM           = "item"
	KEY_DNS            = "dns"
	KEY_MEMORY         = "memory"
	KEY_POD_ID         = "podId"
	KEY_POD_NAME       = "podName"
	KEY_RESOURCE       = "resource"
	KEY_VCPU           = "vcpu"
	KEY_TTY            = "tty"
	KEY_TYPE           = "type"
	KEY_VALUE          = "value"
	KEY_NAME           = "name"
	KEY_IMAGE          = "image"
	KEY_VOLUMES        = "volumes"
	KEY_CONTAINERS     = "containers"
	KEY_VOLUME_SOURCE  = "source"
	KEY_VOLUME_DRIVE   = "driver"
	KEY_ENVS           = "envs"
	KEY_CONTAINER_PORT = "containerPort"
	KEY_HOST_PORT      = "hostPort"
	KEY_PROTOCOL       = "protocol"
	KEY_PORTS          = "ports"
	KEY_MOUNTPATH      = "path"
	KEY_READONLY       = "readOnly"
	KEY_VOLUME         = "volume"
	KEY_COMMAND        = "command"
	KEY_CONTAINER_ARGS = "args"
	KEY_WORKDIR        = "workdir"
	VOLUME_TYPE_VFS    = "vfs"
	TYPE_CONTAINER     = "container"
	TYPE_POD           = "pod"
)

const (
	StatusRunning = "running"
	StatusPending = "pending"
	StatusFailed  = "failed"
	StatusSuccess = "succeeded"
)

type HyperImage struct {
	repository  string
	tag         string
	imageID     string
	createdAt   int64
	virtualSize int64
}

// Container JSON Data Structure
type ContainerPort struct {
	Name          string `json:"name"`
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"hostIP"`
}

type EnvironmentVar struct {
	Env   string `json:"env"`
	Value string `json:"value"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly"`
	MountPath string `json:"mountPath"`
}

type WaitingStatus struct {
	Reason string `json:"reason"`
}

type RunningStatus struct {
	StartedAt string `json:"startedAt"`
}

type TermStatus struct {
	ExitCode   int    `json:"exitCode"`
	Reason     string `json:"reason"`
	Message    string `json:"message"`
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
}

type ContainerStatus struct {
	Name        string        `json:"name"`
	ContainerID string        `json:"containerID"`
	Phase       string        `json:"phase"`
	Waiting     WaitingStatus `json:"waiting"`
	Running     RunningStatus `json:"running"`
	Terminated  TermStatus    `json:"terminated"`
}

// Pod JSON Data Structure
type Container struct {
	Name            string           `json:"name"`
	ContainerID     string           `json:"containerID"`
	Image           string           `json:"image"`
	ImageID         string           `json:"imageID"`
	Commands        []string         `json:"commands"`
	Args            []string         `json:"args"`
	Workdir         string           `json:"workingDir"`
	Ports           []ContainerPort  `json:"ports"`
	Environment     []EnvironmentVar `json:"env"`
	Volume          []VolumeMount    `json:"volumeMounts"`
	ImagePullPolicy string           `json:"imagePullPolicy"`
}

type RBDVolumeSource struct {
	Monitors []string `json:"monitors"`
	Image    string   `json:"image"`
	FsType   string   `json:"fsType"`
	Pool     string   `json:"pool"`
	User     string   `json:"user"`
	Keyring  string   `json:"keyring"`
	ReadOnly bool     `json:"readOnly"`
}

type PodVolume struct {
	Name     string          `json:"name"`
	HostPath string          `json:"source"`
	Driver   string          `json:"driver"`
	Rbd      RBDVolumeSource `json:"rbd"`
}

type PodSpec struct {
	Volumes    []PodVolume `json:"volumes"`
	Containers []Container `json:"containers"`
}

type PodStatus struct {
	Phase     string            `json:"phase"`
	Message   string            `json:"message"`
	Reason    string            `json:"reason"`
	HostIP    string            `json:"hostIP"`
	PodIP     []string          `json:"podIP"`
	StartTime string            `json:"startTime"`
	Status    []ContainerStatus `json:"containerStatus"`
}

type PodInfo struct {
	Kind       string    `json:"kind"`
	ApiVersion string    `json:"apiVersion"`
	Vm         string    `json:"vm"`
	Spec       PodSpec   `json:"spec"`
	Status     PodStatus `json:"status"`
}

type HyperPod struct {
	PodID   string
	PodName string
	VmName  string
	Status  string
	PodInfo PodInfo
}

type HyperContainer struct {
	containerID string
	name        string
	podID       string
	status      string
}

type HyperServiceBackend struct {
	HostIP   string `json:"hostip"`
	HostPort int    `json:"hostport"`
}

type HyperService struct {
	ServiceIP   string                `json:"serviceip"`
	ServicePort int                   `json:"serviceport"`
	Protocol    string                `json:"protocol"`
	Hosts       []HyperServiceBackend `json:"hosts"`
}

type HyperClient struct {
	proto  string
	addr   string
	scheme string
}

type AttachToContainerOptions struct {
	Container    string
	InputStream  io.Reader
	OutputStream io.Writer
	ErrorStream  io.Writer
}

type ExecInContainerOptions struct {
	Container    string
	InputStream  io.Reader
	OutputStream io.Writer
	ErrorStream  io.Writer
	Commands     []string
}

type hijackOptions struct {
	in     io.Reader
	stdout io.Writer
	stderr io.Writer
	data   interface{}
}

func NewHyperClient() *HyperClient {
	var (
		scheme = HYPER_SCHEME
		proto  = HYPER_PROTO
		addr   = HYPER_ADDR
	)

	return &HyperClient{
		proto:  proto,
		addr:   addr,
		scheme: scheme,
	}
}

var (
	ErrConnectionRefused = errors.New("Cannot connect to the Hyper daemon. Is 'hyperd' running on this host?")
)

func (cli *HyperClient) encodeData(data string) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != "" {
		if _, err := params.Write([]byte(data)); err != nil {
			return nil, err
		}
	}
	return params, nil
}

// ParseRepositoryTag gets a repos name and returns the right reposName + tag|digest
// The tag can be confusing because of a port in a repository name.
//     Ex: localhost.localdomain:5000/samalba/hipache:latest
//     Digest ex: localhost:5000/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb
func ParseRepositoryTag(repos string) (string, string) {
	n := strings.Index(repos, "@")
	if n >= 0 {
		parts := strings.Split(repos, "@")
		return parts[0], parts[1]
	}
	n = strings.LastIndex(repos, ":")
	if n < 0 {
		return repos, ""
	}
	if tag := repos[n+1:]; !strings.Contains(tag, "/") {
		return repos[:n], tag
	}
	return repos, ""
}

// parseImageName parses a docker image string into two parts: repo and tag.
// If tag is empty, return the defaultImageTag.
func parseImageName(image string) (string, string) {
	repoToPull, tag := ParseRepositoryTag(image)
	// If no tag was specified, use the default "latest".
	if len(tag) == 0 {
		tag = DEFAULT_IMAGE_TAG
	}
	return repoToPull, tag
}

func (cli *HyperClient) clientRequest(method, path string, in io.Reader, headers map[string][]string) (io.ReadCloser, string, int, *net.Conn, *httputil.ClientConn, error) {
	expectedPayload := (method == "POST" || method == "PUT")
	if expectedPayload && in == nil {
		in = bytes.NewReader([]byte{})
	}
	req, err := http.NewRequest(method, path, in)
	if err != nil {
		return nil, "", -1, nil, nil, err
	}
	req.Header.Set("User-Agent", "kubelet")
	req.URL.Host = cli.addr
	req.URL.Scheme = cli.scheme

	if headers != nil {
		for k, v := range headers {
			req.Header[k] = v
		}
	}

	if expectedPayload && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "text/plain")
	}

	var dial net.Conn
	dial, err = net.DialTimeout(HYPER_PROTO, HYPER_ADDR, 32*time.Second)
	if err != nil {
		return nil, "", -1, nil, nil, err
	}

	clientconn := httputil.NewClientConn(dial, nil)
	resp, err := clientconn.Do(req)
	statusCode := -1
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, "", statusCode, &dial, clientconn, ErrConnectionRefused
		}

		return nil, "", statusCode, &dial, clientconn, fmt.Errorf("An error occurred trying to connect: %v", err)
	}

	if statusCode < 200 || statusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, "", statusCode, &dial, clientconn, err
		}
		if len(body) == 0 {
			return nil, "", statusCode, nil, nil, fmt.Errorf("Error: request returned %s for API route and version %s, check if the server supports the requested API version", http.StatusText(statusCode), req.URL)
		}

		return nil, "", statusCode, &dial, clientconn, fmt.Errorf("%s", bytes.TrimSpace(body))
	}

	return resp.Body, resp.Header.Get("Content-Type"), statusCode, &dial, clientconn, nil
}

func (cli *HyperClient) call(method, path string, data string, headers map[string][]string) ([]byte, int, error) {
	params, err := cli.encodeData(data)
	if err != nil {
		return nil, -1, err
	}

	if data != "" {
		if headers == nil {
			headers = make(map[string][]string)
		}
		headers["Content-Type"] = []string{"application/json"}
	}

	body, _, statusCode, dial, clientconn, err := cli.clientRequest(method, path, params, headers)
	if dial != nil {
		defer (*dial).Close()
	}
	if clientconn != nil {
		defer clientconn.Close()
	}
	if err != nil {
		return nil, statusCode, err
	}

	if body == nil {
		return nil, statusCode, err
	}

	defer body.Close()

	result, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, -1, err
	}

	return result, statusCode, nil
}

func (cli *HyperClient) stream(method, path string, in io.Reader, out io.Writer, headers map[string][]string) error {
	body, contentType, _, dial, clientconn, err := cli.clientRequest(method, path, in, headers)
	if dial != nil {
		defer (*dial).Close()
	}
	if clientconn != nil {
		defer clientconn.Close()
	}
	if err != nil {
		return err
	}

	defer body.Close()

	if MatchesContentType(contentType, "application/json") {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		if out != nil {
			out.Write(buf.Bytes())
		}
		return nil
	}
	return nil

}

func MatchesContentType(contentType, expectedType string) bool {
	mimetype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		glog.V(4).Infof("Error parsing media type: %s error: %v", contentType, err)
	}
	return err == nil && mimetype == expectedType
}

func (client *HyperClient) Version() (string, error) {
	body, _, err := client.call("GET", "/version", "", nil)
	if err != nil {
		return "", err
	}

	var info map[string]interface{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return "", err
	}

	version, ok := info["Version"]
	if !ok {
		return "", fmt.Errorf("Can not get hyper version")
	}

	return version.(string), nil
}

func (client *HyperClient) ListPods() ([]HyperPod, error) {
	v := url.Values{}
	v.Set(KEY_ITEM, TYPE_POD)
	body, _, err := client.call("GET", "/list?"+v.Encode(), "", nil)
	if err != nil {
		return nil, err
	}

	var podList map[string]interface{}
	err = json.Unmarshal(body, &podList)
	if err != nil {
		return nil, err
	}

	var result []HyperPod
	for _, pod := range podList["podData"].([]interface{}) {
		fields := strings.Split(pod.(string), ":")
		var hyperPod HyperPod
		hyperPod.PodID = fields[0]
		hyperPod.PodName = fields[1]
		hyperPod.VmName = fields[2]
		hyperPod.Status = fields[3]

		values := url.Values{}
		values.Set(KEY_POD_NAME, hyperPod.PodID)
		body, _, err = client.call("GET", "/pod/info?"+values.Encode(), "", nil)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &hyperPod.PodInfo)
		if err != nil {
			return nil, err
		}

		result = append(result, hyperPod)
	}

	return result, nil
}

func (client *HyperClient) ListContainers() ([]HyperContainer, error) {
	v := url.Values{}
	v.Set(KEY_ITEM, TYPE_CONTAINER)
	body, _, err := client.call("GET", "/list?"+v.Encode(), "", nil)
	if err != nil {
		return nil, err
	}

	var containerList map[string]interface{}
	err = json.Unmarshal(body, &containerList)
	if err != nil {
		return nil, err
	}

	var result []HyperContainer
	for _, container := range containerList["cData"].([]interface{}) {
		fields := strings.Split(container.(string), ":")
		var h HyperContainer
		h.containerID = fields[0]
		if len(fields[1]) < 1 {
			return nil, errors.New("Hyper container name not resolved")
		}
		h.name = fields[1][1:]
		h.podID = fields[2]
		h.status = fields[3]

		result = append(result, h)
	}

	return result, nil
}

func (client *HyperClient) Info() (map[string]interface{}, error) {
	body, _, err := client.call("GET", "/info", "", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (client *HyperClient) ListImages() ([]HyperImage, error) {
	v := url.Values{}
	v.Set("all", "no")
	body, _, err := client.call("GET", "/images/get?"+v.Encode(), "", nil)
	if err != nil {
		return nil, err
	}

	var images map[string][]string
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil, err
	}

	var hyperImages []HyperImage
	for _, image := range images["imagesList"] {
		imageDesc := strings.Split(image, ":")
		if len(imageDesc) != 5 {
			glog.Warning("Hyper: can not parse image info")
			return nil, fmt.Errorf("Hyper: can not parse image info")
		}

		var imageHyper HyperImage
		imageHyper.repository = imageDesc[0]
		imageHyper.tag = imageDesc[1]
		imageHyper.imageID = imageDesc[2]

		createdAt, err := strconv.ParseInt(imageDesc[3], 10, 0)
		if err != nil {
			return nil, err
		}
		imageHyper.createdAt = createdAt

		virtualSize, err := strconv.ParseInt(imageDesc[4], 10, 0)
		if err != nil {
			return nil, err
		}
		imageHyper.virtualSize = virtualSize

		hyperImages = append(hyperImages, imageHyper)
	}

	return hyperImages, nil
}

func (client *HyperClient) RemoveImage(imageID string) error {
	v := url.Values{}
	v.Set(KEY_IMAGEID, imageID)
	_, _, err := client.call("DELETE", "/images?"+v.Encode(), "", nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *HyperClient) RemovePod(podID string) error {
	v := url.Values{}
	v.Set(KEY_POD_ID, podID)
	_, _, err := client.call("DELETE", "/pod?"+v.Encode(), "", nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *HyperClient) StartPod(podID string) error {
	v := url.Values{}
	v.Set(KEY_POD_ID, podID)
	_, _, err := client.call("POST", "/pod/start?"+v.Encode(), "", nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *HyperClient) StopPod(podID string) error {
	v := url.Values{}
	v.Set(KEY_POD_ID, podID)
	v.Set("stopVM", "yes")
	_, _, err := client.call("POST", "/pod/stop?"+v.Encode(), "", nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *HyperClient) PullImage(image string, credential string) error {
	v := url.Values{}
	v.Set(KEY_IMAGENAME, image)

	headers := make(map[string][]string)
	if credential != "" {
		headers["X-Registry-Auth"] = []string{credential}
	}

	err := client.stream("POST", "/image/create?"+v.Encode(), nil, nil, headers)
	if err != nil {
		return err
	}

	return nil
}

func (client *HyperClient) CreatePod(podArgs string) (map[string]interface{}, error) {
	glog.V(5).Infof("Hyper: starting to create pod %s", podArgs)
	body, _, err := client.call("POST", "/pod/create", podArgs, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *HyperClient) hijack(method, path string, hijackOptions hijackOptions) error {
	var params io.Reader
	if hijackOptions.data != nil {
		buf, err := json.Marshal(hijackOptions.data)
		if err != nil {
			return err
		}
		params = bytes.NewBuffer(buf)
	}

	if hijackOptions.stdout == nil {
		hijackOptions.stdout = ioutil.Discard
	}
	if hijackOptions.stderr == nil {
		hijackOptions.stderr = ioutil.Discard
	}
	req, err := http.NewRequest(method, fmt.Sprintf("/v%s%s", HYPER_MINVERSION, path), params)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "kubelet")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "tcp")
	req.Host = HYPER_ADDR

	dial, err := net.Dial(HYPER_PROTO, HYPER_ADDR)
	if err != nil {
		return err
	}

	clientconn := httputil.NewClientConn(dial, nil)
	defer clientconn.Close()
	clientconn.Do(req)
	rwc, br := clientconn.Hijack()
	defer rwc.Close()
	errChanOut := make(chan error, 1)
	errChanIn := make(chan error, 1)
	exit := make(chan bool)
	go func() {
		defer close(exit)
		defer close(errChanOut)
		_, err := io.Copy(hijackOptions.stdout, br)
		errChanOut <- err
	}()
	go func() {
		if hijackOptions.in != nil {
			_, err := io.Copy(rwc, hijackOptions.in)

			rwc.(interface {
				CloseWrite() error
			}).CloseWrite()

			errChanIn <- err
		}
	}()
	<-exit
	select {
	case err = <-errChanIn:
		return err
	case err = <-errChanOut:
		return err
	}
}

func (client *HyperClient) Attach(opts AttachToContainerOptions) error {
	if opts.Container == "" {
		return fmt.Errorf("No Such Container %s", opts.Container)
	}

	v := url.Values{}
	v.Set(KEY_TYPE, TYPE_CONTAINER)
	v.Set(KEY_VALUE, opts.Container)
	path := "/attach?" + v.Encode()
	return client.hijack("POST", path, hijackOptions{
		in:     opts.InputStream,
		stdout: opts.OutputStream,
		stderr: opts.ErrorStream,
	})
}

func (client *HyperClient) Exec(opts ExecInContainerOptions) error {
	if opts.Container == "" {
		return fmt.Errorf("No Such Container %s", opts.Container)
	}

	command, err := json.Marshal(opts.Commands)
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set(KEY_TYPE, TYPE_CONTAINER)
	v.Set(KEY_VALUE, opts.Container)
	v.Set("command", string(command))
	path := "/exec?" + v.Encode()
	return client.hijack("POST", path, hijackOptions{
		in:     opts.InputStream,
		stdout: opts.OutputStream,
		stderr: opts.ErrorStream,
	})
}

func (client *HyperClient) IsImagePresent(repo, tag string) (bool, error) {
	if outputs, err := client.ListImages(); err == nil {
		for _, imgInfo := range outputs {
			if imgInfo.repository == repo && imgInfo.tag == tag {
				return true, nil
			}
		}
	}
	return false, nil
}

func (client *HyperClient) ListServices(podId string) ([]HyperService, error) {
	v := url.Values{}
	v.Set("podId", podId)
	body, _, err := client.call("GET", "/service/list?"+v.Encode(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "doesn't have services discovery") {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var svcList []HyperService
	err = json.Unmarshal(body, &svcList)
	if err != nil {
		return nil, err
	}

	return svcList, nil
}

func (client *HyperClient) UpdateServices(podId string, services []HyperService) error {
	v := url.Values{}
	v.Set("podId", podId)

	serviceData, err := json.Marshal(services)
	if err != nil {
		return err
	}
	v.Set("services", string(serviceData))
	_, _, err = client.call("POST", "/service/update?"+v.Encode(), "", nil)
	if err != nil {
		return err
	}

	return nil
}

func list(client *HyperClient) {
	for {
		pods, err := client.ListPods()
		if err != nil {
			fmt.Printf("List pod error %s\n", err)
		}

		for _, pod := range pods {
			fmt.Printf("Get pod %s, status %s\n", pod.PodName, pod.Status)
		}
		fmt.Printf("\n")
		time.Sleep(50 * time.Millisecond)
	}
}

func clearPods(hyperClient *HyperClient) {
	pods, err := hyperClient.ListPods()
	if err != nil {
		fmt.Printf("ListPods error %s\n", err)
	}
	for _, pod := range pods {
		fmt.Printf("Get pod %s, status %s\n", pod.PodName, pod.Status)
	}
	for _, pod := range pods {
		err = hyperClient.RemovePod(pod.PodID)
		if err != nil {
			fmt.Printf("Remove pod error %s\n", err)
		}
		fmt.Printf("Pod %s removed\n", pod.PodID)
	}
}

func attach() {
	hyperClient := NewHyperClient()
	clearPods(hyperClient)

	podSpec := `{
	    "containers": [
	        {
	            "image": "busybox",
	            "name": "busybox"
	        }
	    ],
	    "id": "busybox",
	    "dns": ["8.8.8.8"],
	    "tty": true,
	    "type": "pod"
	}`
	body, err := hyperClient.CreatePod(podSpec)
	if err != nil {
		fmt.Printf("Create pod error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Pod %s created.\n", body["ID"])

	podId := body["ID"].(string)
	err = hyperClient.StartPod(podId)
	if err != nil {
		fmt.Printf("Start pod error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Pod started.\n")
	}

	containers, err := hyperClient.ListContainers()
	if err != nil {
		fmt.Printf("List containers error: %s\n", err)
		os.Exit(1)
	}

	containerID := containers[0].containerID
	opts := AttachToContainerOptions{
		Container:    containerID,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
	}

	err = hyperClient.Attach(opts)
	if err != nil {
		fmt.Printf("Attach container error: %s\n", err)
		os.Exit(1)
	}
}

func exec() {
	hyperClient := NewHyperClient()
	clearPods(hyperClient)

	podSpec := `{
	    "containers": [
	        {
	            "image": "busybox",
	            "name": "busybox"
	        }
	    ],
	    "id": "busybox",
	    "dns": ["8.8.8.8"],
	    "tty": true,
	    "type": "pod"
	}`
	body, err := hyperClient.CreatePod(podSpec)
	if err != nil {
		fmt.Printf("Create pod error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Pod %s created.\n", body["ID"])

	podId := body["ID"].(string)
	err = hyperClient.StartPod(podId)
	if err != nil {
		fmt.Printf("Start pod error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Pod started.\n")
	}

	containers, err := hyperClient.ListContainers()
	if err != nil {
		fmt.Printf("List containers error: %s\n", err)
		os.Exit(1)
	}

	containerID := containers[0].containerID
	opts := ExecInContainerOptions{
		Container:    containerID,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Commands:     []string{"echo", "aaaa"},
	}

	err = hyperClient.Exec(opts)
	if err != nil {
		fmt.Printf("Exec container error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	attach()
}
