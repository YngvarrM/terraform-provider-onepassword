package onepassword

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVault() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVaultRead,
		CreateContext: resourceVaultCreate,
		DeleteContext: resourceVaultDelete,
		UpdateContext: resourceVaultUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceVaultRead(ctx, d, meta); err.HasError() {
					return []*schema.ResourceData{d}, errors.New(err[0].Summary)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"safety_lock": {
				Type:     schema.TypeBool,
				Optional:    true,
				Default: false,
				Description: "Prevents provider from Vault removal",
			},
			"incognito": {
				Type:     schema.TypeBool,
				Optional:    true,
				Default: false,
				Description: "Remove resource creator from resource",
			},
		},
	}
}

func resourceVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	v, err := m.onePassClient.ReadVault(getID(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Name); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)

	v, err := m.onePassClient.CreateVault(&Vault{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("incognito") == true {
		err1 := m.onePassClient.DeleteVaultMember(v.UUID, m.onePassClient.Email) // deletes user that created resource
		if err1 != nil {
			return diag.FromErr(err)
		}
	}
	return resourceVaultRead(ctx, d, meta)
}

func resourceVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("safety_lock") == false {
		m := meta.(*Meta)
		err := m.onePassClient.DeleteVault(getID(d))
		if err == nil {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.FromErr(fmt.Errorf("Vault removal disabled. Please delete manually or remove safe lock."))
}

func resourceVaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)

	v := &Vault{
		Name: d.Get("name").(string),
	}

	if err := m.onePassClient.UpdateVault(getID(d), v); err != nil {
		return diag.FromErr(err)
	}
	return resourceVaultRead(ctx, d, meta)
}
