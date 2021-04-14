package onepassword

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVaultMember() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVaultMemberRead,
		CreateContext: resourceVaultMemberCreate,
		DeleteContext: resourceVaultMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"Vault": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceVaultMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	VaultID, userID, err := resourceVaultMemberExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	v, err := m.onePassClient.ListVaultMembers(VaultID)
	if err != nil {
		return diag.FromErr(err)
	}

	var found string
	for _, member := range v {
		if member.UUID == userID {
			found = member.UUID
		}
	}

	if found == "" {
		d.SetId("")
		return nil
	}

	d.SetId(resourceVaultMemberBuildID(VaultID, found))
	d.Set("Vault", VaultID)
	d.Set("user", found)
	return nil
}

func resourceVaultMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	err := m.onePassClient.CreateVaultMember(
		d.Get("Vault").(string),
		d.Get("user").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceVaultMemberBuildID(d.Get("Vault").(string), d.Get("user").(string)))
	return resourceVaultMemberRead(ctx, d, meta)
}

func resourceVaultMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	VaultID, userID, err := resourceVaultMemberExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	err = m.onePassClient.DeleteVaultMember(
		VaultID,
		userID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourceVaultMemberBuildID will conjoin the Vault ID and user ID into a single string
// This is used as the resource ID.
//
// Note that user ID is being lowercased. Some operations require this user ID to be uppercased.
// Use the resourceVaultMemberExtractID function to correctly reverse this encoding.
func resourceVaultMemberBuildID(VaultID, userID string) string {
	return strings.ToLower(VaultID + "-" + strings.ToLower(userID))
}

// resourceVaultMemberExtractID will split the Vault ID and user ID from a given resource ID
//
// Note that user ID is being uppercased. Some operations require this user ID to be uppercased.
func resourceVaultMemberExtractID(id string) (VaultID, userID string, err error) {
	spl := strings.Split(id, "-")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("Improperly formatted Vault member string. The format \"Vaultid-userid\" is expected")
	}
	return spl[0], strings.ToUpper(spl[1]), nil
}
