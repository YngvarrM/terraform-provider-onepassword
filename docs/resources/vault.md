# onepassword_vault

This resource can create vaults in your 1password account.

## Example Usage

```hcl
resource "onepassword_vault" "this" {
    name = "new-vault"
}
```

## Argument Reference

* `name` - (Required) vault name.

* `incognito` - (Optional) Terraform will automatically remove your user from this resource after creation.

* `safe_lock` - (Optional) It prevents the vault from being accidentally removed. You will get an error if you try to remove the vault with this parameter equal to true.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - vault id.

1Password Vaults can be imported using the `id`, e.g.

```
terraform import onepassword_vault.vault 7kalogoe3kirwf5aizotkbzrpq
```
