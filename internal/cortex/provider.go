// Copyright 2021 The Terraform Provider Cortex developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cortex

import (
	"context"
	"net/http"

	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	EnvCortexAddress  = "CORTEX_ADDRESS"
	EnvCortexApiKey   = "CORTEX_API_KEY"
	EnvCortexTenantId = "CORTEX_TENANT_ID"
	EnvCortexUser = "CORTEX_API_USER"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(EnvCortexAddress, ""),
				Description: "Address of the Cortex cluster",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(EnvCortexApiKey, ""),
				Description: "API key, used as basic auth password.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(EnvCortexUser, ""),
				Description: "used as basic auth user.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(EnvCortexTenantId, ""),
				Description: "Tenant ID, passed as X-Org-ScopeID HTTP header.",
			},
			"backend": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cortex",
				Description: "Backend type to interact with: <cortex|loki>",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cortex_rules":        resourceRules(),
			"cortex_alertmanager": resourceAlertmanager(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var (
		address         = d.Get("address").(string)
		apiKey          string
		defaultTenantID string
		backend         string
		user            string
	)
	if data, ok := d.GetOk("api_key"); ok {
		apiKey = data.(string)
	}
	if data, ok := d.GetOk("tenant_id"); ok {
		defaultTenantID = data.(string)
	}	
	if data, ok := d.GetOk("backend"); ok {
		backend = data.(string)
	}
	if data, ok := d.GetOk("user"); ok {
		user = data.(string)
	}
	var diags diag.Diagnostics

	var f cortexClientFunc = func(tenantID string) (*client.CortexClient, error) {
		if tenantID == "" {
			tenantID = defaultTenantID
		}
		clientConfig := client.Config{
			Key:     apiKey,
			Address: address,
			ID:      tenantID,
			User:    user,
		}
		if backend == "loki" {
			clientConfig.UseLegacyRoutes = true
		}

		client, err := client.New(clientConfig)

		if err != nil {
			return nil, err
		}

		// Setup Terraform-SDK transport to enable debugging via TF_LOGS=debug.
		tr := client.Client.Transport
		if client.Client.Transport == nil {
			tr = http.DefaultTransport
		}
		client.Client.Transport = logging.NewTransport("cortex", tr)

		return client, err
	}

	return f, diags
}

type cortexClientFunc func(string) (*client.CortexClient, error)
