package vault

import (
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
)

func GetConfig(path string, options ...Option) (*Config, error) {
	// Read options
	vaultConfig := readOptions(path, options...)

	// Init Vault client
	client, err := createVaultClient(vaultConfig.addr, vaultConfig.token)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
		return nil, err
	}

	// Read Vault data
	secret, err := read(client, fmt.Sprintf("%s%s", vaultConfig.mountPath, vaultConfig.secretPath))
	if err != nil {
		log.Fatalf("unable to read Vault secret: %v", err)
		return nil, err
	}

	var vaultData map[string]interface{}
	var ok bool
	if vaultConfig.kvversion == 2 {
		vaultData, ok = secret.Data["data"].(map[string]interface{})
		if !ok {
			vaultData = secret.Data
		}
	} else {
		vaultData = secret.Data
	}

	vaultConfig.data = vaultData

	return vaultConfig, nil
}

func createVaultClient(address, token string) (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = address

	client, err := vault.NewClient(config)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	client.SetToken(token)

	return client, nil
}

func read(client *vault.Client, path string) (*vault.Secret, error) {
	resp, err := client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
