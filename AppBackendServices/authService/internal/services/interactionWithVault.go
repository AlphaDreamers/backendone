package services

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func (i *Impl) InteractionWithVault(userId string, phrase string) error {
	config := api.DefaultConfig()
	config.Address = viper.GetString("vault-address")
	client, err := api.NewClient(config)
	if err != nil {
		i.log.Fatalf("Error creating Vault client: %v", err)
		return err
	}
	health, err := client.Sys().Health()
	if err != nil {
		i.log.Errorf("Vault connection check failed: %v", err)
		return fmt.Errorf("vault connection failed: %w", err)
	}
	i.log.Infof("Connected to Vault (version: %s)", health.Version)

	rootToken := viper.GetString("vault-root-token")
	secretPath := viper.GetString("vault-secret-path")

	client.SetToken(rootToken)
	secrets := map[string]interface{}{}
	secrets[userId] = phrase

	_, err = client.KVv2("secret").Put(context.Background(), secretPath, secrets)
	if err != nil {
		i.log.Warn("Error writing secret: %v", err)
	}
	i.log.Println("Secret written successfully")
	result, err := client.KVv2("secret").Get(context.Background(), secretPath)
	if err != nil {
		i.log.Warn("Error reading secret: %v", err)
	}
	i.log.Println("Secret read successfully", result.Data)
	i.log.Println("Secret written successfully", result.Raw)
	i.log.Println(result.VersionMetadata)
	i.log.Info("Successfully interacting with Vault client")
	return nil
}
