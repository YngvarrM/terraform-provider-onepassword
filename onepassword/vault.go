package onepassword

import (
	"encoding/json"
	"fmt"
)

const VaultResource = "vault"

type Vault struct {
	UUID string
	Name string
}

func (o *OnePassClient) ReadVault(id string) (*Vault, error) {
	vault := &Vault{}
	res, err := o.runCmd(opPasswordGet, VaultResource, id)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, vault); err != nil {
		return nil, err
	}
	return vault, nil
}

func (o *OnePassClient) CreateVault(v *Vault) (*Vault, error) {
	args := []string{opPasswordCreate, VaultResource, v.Name}
	res, err := o.runCmd(args...)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (o *OnePassClient) DeleteVault(id string) error {
	return o.Delete(VaultResource, id)
}

func (o *OnePassClient) CreateVaultMember(vaultID string, userID string) error {
	args := []string{opPasswordAdd, UserResource, userID, vaultID}
	_, err := o.runCmd(args...)
	return err
}

func (o *OnePassClient) DeleteVaultMember(vaultID string, userID string) error {
	args := []string{opPasswordRemove, UserResource, userID, vaultID}
	_, err := o.runCmd(args...)
	return err
}

func (o *OnePassClient) ListVaultMembers(id string) ([]User, error) {
	users := []User{}
	if id == "" {
		return users, fmt.Errorf("Must provide an identifier to list group members")
	}

	res, err := o.runCmd(opPasswordList, "users", "--"+VaultResource, id)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (o *OnePassClient) UpdateVault(id string, v *Vault) error {
	args := []string{opPasswordEdit, VaultResource, id, "--name=" + v.Name}
	_, err := o.runCmd(args...)
	return err
}
