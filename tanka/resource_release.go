package tanka

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/tanka"
)

func resourceRelease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReleaseCreate,
		ReadContext:   resourceReleaseRead,
		UpdateContext: resourceReleaseUpdate,
		DeleteContext: resourceReleaseDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"source_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tanka/environments/default",
			},
			"config_inline": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"config_local": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceReleaseRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	_ = d.Get("last_updated")

	return nil
}

func resourceReleaseCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	name := fmt.Sprintf("%v", d.Get("name"))

	api_server := fmt.Sprintf("%v", d.Get("api_server"))
	namespace := fmt.Sprintf("%v", d.Get("namespace"))
	source_path := fmt.Sprintf("%v", d.Get("source_path"))
	config_inline := d.Get("config_inline").(map[string]any)
	config_local := fmt.Sprintf("%v", d.Get("config_local"))

	raw, err := json.Marshal(config_inline)
	if err != nil {
		return diag.FromErr(err)
	}
	ci := string(raw[:])

	cl, err := getLocalConfig(config_local)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tankaApply(api_server, namespace, ci, cl, source_path)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceReleaseRead(ctx, d, m)
}

func resourceReleaseUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	api_server := fmt.Sprintf("%v", d.Get("api_server"))
	namespace := fmt.Sprintf("%v", d.Get("namespace"))
	source_path := fmt.Sprintf("%v", d.Get("source_path"))
	config_inline := d.Get("config_inline").(map[string]any)
	config_local := fmt.Sprintf("%v", d.Get("config_local"))

	raw, err := json.Marshal(config_inline)
	if err != nil {
		return diag.FromErr(err)
	}
	ci := string(raw[:])

	cl, err := getLocalConfig(config_local)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Running Tanka Apply")

	err = tankaApply(api_server, namespace, ci, cl, source_path)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceReleaseRead(ctx, d, m)
}

func resourceReleaseDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	api_server := fmt.Sprintf("%v", d.Get("api_server"))
	namespace := fmt.Sprintf("%v", d.Get("namespace"))
	source_path := fmt.Sprintf("%v", d.Get("source_path"))
	config_inline := d.Get("config_inline").(map[string]any)
	config_local := fmt.Sprintf("%v", d.Get("config_local"))

	raw, err := json.Marshal(config_inline)
	if err != nil {
		return diag.FromErr(err)
	}
	ci := string(raw[:])

	cl, err := getLocalConfig(config_local)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Running Tanka Delete")

	err = tankaDelete(api_server, namespace, ci, cl, source_path)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getLocalConfig(config_input string) (config string, err error) {

	config_type := "json"
	if strings.Contains(config_input, "://") {
		split := strings.Split(config_input, "://")
		config_type = split[0]
	}

	var raw []byte

	switch config_type {
	case "json":
		config = config_input
		// raw, err = json.Marshal(config_input)
		// if err != nil {
		// 	return
		// }
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

func createBaseOpts(api_server, namespace, config_inline, config_local string) (opts tanka.ApplyBaseOpts, err error) {

	var TLACode jsonnet.InjectedCode
	TLACode.Set("apiServer", "\""+api_server+"\"")
	TLACode.Set("namespace", "\""+namespace+"\"")
	TLACode.Set("config_inline", config_inline)
	TLACode.Set("config_local", config_local)
	opts.TLACode = TLACode

	opts.AutoApprove = "true"
	opts.DryRun = "none"
	opts.Force = true

	return
}

func tankaApply(api_server, namespace, config_inline, config_local, baseDir string) (err error) {

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

func tankaDelete(api_server, namespace, config_inline, config_local, baseDir string) (err error) {
	opts, _ := createBaseOpts(api_server, namespace, config_inline, config_local)

	var deleteOpts tanka.DeleteOpts
	deleteOpts.ApplyBaseOpts = opts

	err = tanka.Delete(baseDir, deleteOpts)
	if err != nil {
		return err
	}

	return
}
