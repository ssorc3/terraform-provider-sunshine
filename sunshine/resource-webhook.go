package sunshine

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ssorc3/sunshine-client"
)

func resourceCustomIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomIntegrationCreate,
		DeleteContext: resourceCustomIntegrationDelete,
		UpdateContext: resourceCustomIntegrationUpdate,
		ReadContext:   resourceCustomIntegrationRead,
		Schema: map[string]*schema.Schema{
			"displayName": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"triggers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCustomIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sunshine.Client)

	var diags diag.Diagnostics

	displayName := d.Get("displayName").(string)
	target := d.Get("target").(string)
	triggers := d.Get("triggers").([]string)

	resource, err := c.CreateCustomIntegration(displayName, target, triggers, false, false)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while creating webhook",
			Detail:   err.Error(),
		})
	}

	d.SetId(resource.Id)

	resourceCustomIntegrationRead(ctx, d, m)

	return diags
}

func resourceCustomIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sunshine.Client)

	var diags diag.Diagnostics

	id := d.Id()

	resource, err := c.GetCustomIntegration(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch custom integration resource",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("displayName", resource.DisplayName)
	d.Set("target", resource.Webhooks[0].Target)
	d.Set("triggers", resource.Webhooks[0].Triggers)

	return diags
}

func resourceCustomIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sunshine.Client)

	var diags diag.Diagnostics

	integrationId := d.Id()

	err := c.DeleteCustomIntegration(integrationId)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while deleting webhook",
			Detail:   err.Error(),
		})
	}

	d.SetId("")

	return diags
}

func resourceCustomIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sunshine.Client)

	var diags diag.Diagnostics

	id := d.Id()
	displayName := d.Get("displayName").(string)

	c.UpdateCustomIntegration(id, displayName)

	resourceCustomIntegrationRead(ctx, d, m)

	return diags
}
