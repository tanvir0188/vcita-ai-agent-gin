package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/tanvir0188/vcita-ai-agent/cmd/api"
	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/crypto"
	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

func main() {
	debug.SetMemoryLimit(800 * 1024 * 1024)
	fmt.Println("main started")

	// clients, err := medicationreminder.GetClientsWithCurrentMedications()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "failed to get clients: %v\n", err)
	// 	os.Exit(1)
	// }

	// data, err := json.MarshalIndent(clients, "", "  ")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "failed to marshal clients: %v\n", err)
	// 	os.Exit(1)
	// }

	// if err := os.WriteFile("data.json", data, 0644); err != nil {
	// 	fmt.Fprintf(os.Stderr, "failed to write data.json: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("saved %d clients to data.json\n", len(clients))

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("run started")

	cfg := config.Envs

	fmt.Println("config loaded")

	if err := logger.Init(); err != nil {
		return err
	}

	defer logger.Log.Sync()

	enc, err := crypto.NewEncryptor(cfg.PHIEncryptionKey)
	if err != nil {
		return err
	}

	db, err := store.New(cfg.DBDSN, enc, logger.Log)
	if err != nil {
		return err
	}
	defer db.Close()

	auditor, err := audit.New(
		cfg.AuditLogPath,
		db.WriteAuditLog,
	)
	if err != nil {
		return err
	}
	defer auditor.Close()

	server := api.NewAPIServer(
		cfg,
		db,
		auditor,
		logger.Log,
	)

	return server.Run()
}
