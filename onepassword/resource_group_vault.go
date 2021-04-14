package onepassword

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupVault() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceGroupVaultRead,
		CreateContext: resourceGroupVaultCreate,
		DeleteContext: resourceGroupVaultDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"vault": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceGroupVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID, vaultID, err := resourceGroupVaultExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	v, err := m.onePassClient.ListGroupVaults(groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	var found string
	for _, member := range v {
		if member.UUID == vaultID {
			found = member.UUID
		}
	}

	if found == "" {
		d.SetId("")
		return nil
	}

	d.SetId(resourceGroupVaultBuildID(groupID, found))
	d.Set("group", groupID)
	d.Set("vault", found)
	return nil
}

func resourceGroupVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	err := m.onePassClient.CreateGroupVault(
		d.Get("vault").(string),
		d.Get("group").(string),

	)
	if err != nil {
		return diag.FromErr(err)
	}
	


	d.SetId(resourceGroupVaultBuildID(d.Get("group").(string), d.Get("vault").(string)))
	return resourceGroupVaultRead(ctx, d, meta)
}

func resourceGroupVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID, vaultID, err := resourceGroupVaultExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	err = m.onePassClient.DeleteGroupVault(
		groupID,
		vaultID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourceGroupVaultBuildID will conjoin the group ID and vault ID into a single string
// This is used as the resource ID.
//
// Note that vault ID is being lowercased. Some operations require this vault ID to be uppercased.
// Use the resourceGroupVaultExtractID function to correctly reverse this encoding.
func resourceGroupVaultBuildID(groupID, vaultID string) string {
	return strings.ToLower(groupID + "-" + vaultID)
}

// resourceGroupVaultExtractID will split the group ID and vault ID from a given resource ID
//
// Note that vault ID is being uppercased. Some operations require this vault ID to be uppercased.
func resourceGroupVaultExtractID(id string) (groupID, vaultID string, err error) {
	spl := strings.Split(id, "-")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("Improperly formatted group vault string. The format \"groupid-vaultID\" is expected")
	}
	return spl[0], spl[1], nil
}
