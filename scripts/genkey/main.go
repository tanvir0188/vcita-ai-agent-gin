// scripts/genkey/main.go
// Run once during setup to generate the AES-256 encryption key:
//
//	go run ./scripts/genkey
//
// Copy the output into your ENCRYPTION_KEY_B64 environment variable.
// Store it in a secrets manager (AWS Secrets Manager, Vault, etc.) — never in source control.
package main

import (
	"fmt"
	"log"

	"github.com/tanvir0188/vcita-ai-agent/internal/crypto"
)

func main() {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate key: %v", err)
	}

	fmt.Println("Generated AES-256 encryption key (base64):")
	fmt.Println(key)
	fmt.Println()
	fmt.Println("Add to your environment:")
	fmt.Printf("  ENCRYPTION_KEY_B64=%s\n", key)
	fmt.Println()
	fmt.Println("WARNING: Store this key securely. Loss of this key means loss of all encrypted data.")
	fmt.Println("WARNING: Never commit this key to source control.")
}
