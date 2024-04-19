package provider

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/tanka"
)

// Client -
type Client struct {
	Endpoint             string
	Token                string
	ClusterCaCertificate string
}

func NewClient(endpoint, token, cluster_ca_certificate *string) (client *Client, err error) {
	c := Client{
		Endpoint:             *endpoint,
		Token:                *token,
		ClusterCaCertificate: *cluster_ca_certificate,
	}

	c.setCredentials()

	client = &c
	return
}

func (c *Client) setCredentials() {
	cfgJSON := bytes.Buffer{}

	cluster_config_identifier := strings.Replace(c.Endpoint, ".", "_", -1)

	cmd := exec.Command("kubectl", "config", "set-credentials", cluster_config_identifier, "--token", c.Token)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	cmd = exec.Command("kubectl", "config", "set-context", cluster_config_identifier, "--cluster", cluster_config_identifier, "--user", cluster_config_identifier)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	cmd = exec.Command("kubectl", "config", "set-cluster", cluster_config_identifier, "--server", c.Endpoint)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	cmd = exec.Command("kubectl", "config", "set", "clusters."+cluster_config_identifier+".certificate-authority-data", c.ClusterCaCertificate)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
}

func createBaseOpts(api_server, namespace, config_inline, config_local string) (opts tanka.ApplyBaseOpts, err error) {

	var TLACode jsonnet.InjectedCode
	TLACode.Set("apiServer", "\""+api_server+"\"")
	TLACode.Set("namespace", "\""+namespace+"\"")
	_ = config_inline
	// TLACode.Set("config_inline", config_inline)
	TLACode.Set("config_local", config_local)
	opts.TLACode = TLACode

	opts.AutoApprove = "true"
	opts.DryRun = "none"
	opts.Force = true

	return
}

func (c *Client) Apply(api_server, namespace, config_inline, config_local, baseDir string) (err error) {
	opts, _ := createBaseOpts(api_server, namespace, config_inline, config_local)

	var applyOpts tanka.ApplyOpts
	applyOpts.ApplyBaseOpts = opts

	applyOpts.ApplyStrategy = "server"

	err = tanka.Apply(baseDir, applyOpts)
	if err != nil {
		return err
	}

	return
}

func (c *Client) Delete(api_server, namespace, config_inline, config_local, baseDir string) (err error) {
	opts, _ := createBaseOpts(api_server, namespace, config_inline, config_local)

	var deleteOpts tanka.DeleteOpts
	deleteOpts.ApplyBaseOpts = opts

	err = tanka.Delete(baseDir, deleteOpts)
	if err != nil {
		return err
	}

	return
}

func (c *Client) getLocalConfig(config_input string) (config string, err error) {

	config_type := "json"
	if strings.Contains(config_input, "://") {
		split := strings.Split(config_input, "://")
		config_type = split[0]
	}

	var raw []byte

	switch config_type {
	case "json":
		config = config_input
	case "file":
		split := strings.Split(config_input, "://")
		raw, err = os.ReadFile(split[1])
		if err != nil {
			return
		}
		config = string(raw[:])
	case "http", "https":
		raw, err = getHttpContent(config_input)
		if err != nil {
			return
		}
		config = string(raw[:])
	default:
		err = fmt.Errorf("unknown protocol used in config_local")
		return
	}


	return
}

func getHttpContent(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	return data, nil
}













// // HostURL - Default Hashicups URL
// const HostURL string = "http://localhost:19090"

// // AuthStruct -
// type AuthStruct struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// // AuthResponse -
// type AuthResponse struct {
// 	UserID   int    `json:"user_id`
// 	Username string `json:"username`
// 	Token    string `json:"token"`
// }

// NewClient -
// func NewClient(host, username, password *string) (*Client, error) {
// 	c := Client{
// 		HTTPClient: &http.Client{Timeout: 10 * time.Second},
// 		// Default Hashicups URL
// 		HostURL: HostURL,
// 		Auth: AuthStruct{
// 			Username: *username,
// 			Password: *password,
// 		},
// 	}

// 	if host != nil {
// 		c.HostURL = *host
// 	}

// 	ar, err := c.SignIn()
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.Token = ar.Token

// 	return &c, nil
// }


// func (c *Client) doRequest(req *http.Request) ([]byte, error) {
// 	req.Header.Set("Authorization", c.Token)

// 	res, err := c.HTTPClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if res.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
// 	}

// 	return body, err
// }
