package medication

import (
	"log"
	"sync"
	"time"

	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"

	"gorm.io/gorm"
)

// Store interface for reminder persistence.
type Store interface {
	CreateOrUpdateReminder(matterUID string, resp *ai.MedicationRefillResponse) error
}

// ProcessMedicationReminders runs the medication‑refill prediction for all clients that have current medications.
// It respects the ClientSyncState table to avoid re‑processing the same client more than once per day,
// unless a partial reminder is required. The function processes clients in parallel with a maximum of 10 concurrent workers.
func ProcessMedicationReminders(db *gorm.DB, refillStore Store) error {
	log.Printf("[MedicationReminder] start fetching all clients")
	// 1️⃣ Retrieve all clients (paged internally).
	clients, err := GetAllClients()
	if err != nil {
		return err
	}

	// Determine the start of today in the server's timezone.
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 2️⃣ Filter clients that need a check.
	var toProcess []Client
	for _, c := range clients {
		meds := ExtractCurrentMedications(c.Notes)
		if meds == "" {
			continue
		}
		// Load sync state for this client.
		var syncState store.ClientSyncState
		err := db.Where("matter_uid = ?", c.MatterUid).First(&syncState).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		needCheck := false
		if err == gorm.ErrRecordNotFound {
			needCheck = true
		} else {
			if syncState.PartialReminderNeeded {
				needCheck = true
			} else if syncState.LastAICheckAt == nil || syncState.LastAICheckAt.Before(startOfDay) {
				needCheck = true
			}
		}
		if needCheck {
			c.Notes = meds
			toProcess = append(toProcess, c)
		}
	}

	log.Printf("[MedicationReminder] %d clients need processing", len(toProcess))

	// 3️⃣ Worker pool – up to 10 concurrent predictions.
	const maxWorkers = 10
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for _, client := range toProcess {
		wg.Add(1)
		clientCopy := client // capture loop variable
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // acquire slot
			defer func() { <-sem }() // release slot

			log.Printf("[MedicationReminder] processing client %s", clientCopy.MatterUid)

			// Predict medication refill only for the test matterUID.
			var resp *ai.MedicationRefillResponse
			if clientCopy.MatterUid == "2y1ncrj8uh7qk8ye" {
				log.Printf("inititationg the prediction for test account.")
				var err error
				resp, err = ai.PredictMedicationRefill(clientCopy.Notes)
				if err != nil {
					log.Printf("[MedicationReminder] AI prediction failed for client %s: %v", clientCopy.MatterUid, err)
					return
				}
			} else {
				// No prediction – create empty response to keep sync state consistent.
				resp = &ai.MedicationRefillResponse{ReminderNeeded: false, PartialReminder: false, Medications: nil}
			}
			log.Printf("ai response: %v", resp)

			// Create a reminder message only if we have a valid response.
			if resp != nil && resp.Response != nil {
				if _, err := vcita.CreateMessage(clientCopy.UID, *resp.Response); err != nil {
					log.Printf("[MedicationReminder] failed to create reminder message for client %s: %v", clientCopy.MatterUid, err)
				}
			}

			// Update or create the sync state entry.
			if err := refillStore.CreateOrUpdateReminder(clientCopy.MatterUid, resp); err != nil {
				log.Printf("[MedicationReminder] failed to sync reminder state for client %s: %v", clientCopy.MatterUid, err)
			}
		}()
	}

	wg.Wait()
	log.Printf("[MedicationReminder] all tasks completed")
	return nil
}

// ScheduleDailyReminder runs the medication reminder task. Call this from the main server startup.
func ScheduleDailyReminder(db *gorm.DB, store Store) {
	log.Printf("[MedicationReminder] starting")

	if err := ProcessMedicationReminders(db, store); err != nil {
		log.Printf("[MedicationReminder] error: %v", err)
	} else {
		log.Printf("[MedicationReminder] completed successfully")
	}
}
