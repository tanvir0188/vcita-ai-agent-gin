// scripts/register-webhooks/main.go
// Run this once after deployment to register all required webhook subscriptions
// with the vcita inTandem platform.
//
// Usage:
//
//	VCITA_DIRECTORY_TOKEN=xxx VCITA_BUSINESS_TOKEN=xxx WEBHOOK_CALLBACK_URL=https://yourdomain.com/webhooks/vcita \
//	  go run ./scripts/register-webhooks
package main

import (
	// "context"
	"fmt"
	"log"
	"os"

	// "time"

	// "github.com/tanvir0188/vcita-ai-agent/internal/vcita"
	"go.uber.org/zap"
)

// Subscriptions we need from vcita based on our feature set
var subscriptions = []string{
	"message.client_sent_message",   // main AI trigger
	"message.business_sent_message", // audit logging
	"client.created",                // initial medication check
	"client.updated",                // medication note changes
	"appointment.requested",         // AI scheduling
	"appointment.scheduled",         // audit
	"appointment.rescheduled",       // audit
	"appointment.cancelled",         // audit
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() //nolint:errcheck

	// directoryToken := requireEnv("VCITA_DIRECTORY_TOKEN")
	// businessToken := requireEnv("VCITA_BUSINESS_TOKEN")
	// callbackURL := requireEnv("WEBHOOK_CALLBACK_URL")

	// client := vcita.NewAPIClient(directoryToken, businessToken)
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	fmt.Printf("Registering %d webhook subscriptions...\n", len(subscriptions))

	// for _, eventType := range subscriptions {
	// 	if err := client.RegisterWebhookSubscription(ctx, eventType, callbackURL); err != nil {
	// 		log.Printf("FAILED to register %s: %v", eventType, err)
	// 	} else {
	// 		fmt.Printf("  ✓ %s\n", eventType)
	// 	}
	// }

	fmt.Println("Webhook registration not implemented yet")

	fmt.Println("Done.")
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}
