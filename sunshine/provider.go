package sunshine

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ssorc3/sunshine-client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"appId": {
				Type:     schema.TypeString,
				Optional: false,
			},
			"appSecret": {
				Type:     schema.TypeString,
				Optional: false,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: func() (interface{}, error) { return "EU", nil },
			},
			"appKeyId": {
				Type:     schema.TypeString,
				Optional: false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"webhook": resourceCustomIntegration(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	appId := d.Get("appId").(string)
	appSecret := d.Get("appSecret").(string)
	appKeyId := d.Get("appKeyId").(string)

	var diags diag.Diagnostics

	region := d.Get("region").(string)

	var baseUrl string
	if region == "EU" {
		baseUrl = "https://api.eu-1.smooch.io/"
	} else if region == "US" {
		baseUrl = "https://api.smooch.io/"
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid region",
			Detail:   "region must be one of EU, US",
		})

		return nil, diags
	}

	c := sunshine.NewClient(baseUrl, appId, appKeyId, appSecret)

	return c, diags
}
