package provider

import (
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "time"

	"bytes"
	"os"
	"os/exec"
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

func (c *Client) setCredentials() {
	cfgJSON := bytes.Buffer{}

	cmd := exec.Command("kubectl", "config", "set-credentials", c.Endpoint, "--token", c.Token)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	cmd = exec.Command("kubectl", "config", "set-context", c.Endpoint, "--cluster", c.Endpoint)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	cmd = exec.Command("kubectl", "config", "set-cluster", c.Endpoint, "--certificate-authority", c.ClusterCaCertificate, "--server", c.Endpoint)
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
}

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
