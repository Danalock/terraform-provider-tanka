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

	// A kube config context name creates issues if it contains dots
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

func createBaseOpts(api_server, namespace, config, config_override string) (opts tanka.ApplyBaseOpts) {

	var TLACode jsonnet.InjectedCode
	TLACode.Set("apiServer", "\""+api_server+"\"")
	TLACode.Set("namespace", "\""+namespace+"\"")
	TLACode.Set("config", config)
	TLACode.Set("config_override", config_override)
	opts.TLACode = TLACode

	opts.AutoApprove = "true"
	opts.DryRun = "none"
	opts.Force = true

	return
}

func (c *Client) Apply(api_server, namespace, config, config_override, baseDir string) (err error) {
	opts := createBaseOpts(api_server, namespace, config, config_override)

	var applyOpts tanka.ApplyOpts
	applyOpts.ApplyBaseOpts = opts

	applyOpts.ApplyStrategy = "server"

	err = tanka.Apply(baseDir, applyOpts)
	if err != nil {
		return err
	}

	return
}

func (c *Client) Delete(api_server, namespace, config, config_override, baseDir string) (err error) {
	opts := createBaseOpts(api_server, namespace, config, config_override)

	var deleteOpts tanka.DeleteOpts
	deleteOpts.ApplyBaseOpts = opts

	err = tanka.Delete(baseDir, deleteOpts)
	if err != nil {
		return err
	}

	return
}

func (c *Client) parseConfig(config_input string) (config string, err error) {

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
		err = fmt.Errorf("unknown protocol used in config")
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
