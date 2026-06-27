package ai

import (
	"fmt"
	"time"
)

// validateAndFixRefill overrides the LLM's reminder/partial_reminder/response
// fields to guarantee correctness per the specification rules (mirrors _validate_and_fix in ai.py).
func validateAndFixRefill(result *MedicationRefillResponse, currentDate string) {
	current, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return
	}

	// --- Recalculate needs_refill for each medication ---
	for i := range result.Medications {
		med := &result.Medications[i]
		confidence := med.Confidence
		refillDateStr := med.NextRefillDate

		if refillDateStr == nil || confidence == "unknown" {
			med.NeedsRefill = false
			continue
		}

		refillDate, err := time.Parse("2006-01-02", *refillDateStr)
		if err != nil {
			med.NeedsRefill = false
			med.Confidence = "unknown"
			continue
		}

		daysRemaining := int(refillDate.Sub(current).Hours() / 24)

		switch confidence {
		case "high":
			med.NeedsRefill = daysRemaining <= 7
		case "estimated":
			med.NeedsRefill = daysRemaining <= 5
		default:
			med.NeedsRefill = false
		}
	}

	// --- Recalculate reminder logic ---
	medicationCount := len(result.Medications)
	alertCount := 0
	for _, med := range result.Medications {
		if med.NeedsRefill {
			alertCount++
		}
	}

	result.ReminderNeeded = alertCount > 0
	result.PartialReminder = alertCount > 0 && alertCount < medicationCount

	// --- Recalculate response message ---
	if !result.ReminderNeeded {
		result.Response = nil
	} else if alertCount == 1 {
		var refillMedName string
		for _, med := range result.Medications {
			if med.NeedsRefill {
				refillMedName = med.Name
				break
			}
		}
		if refillMedName == "" {
			refillMedName = "A medication"
		}
		msg := fmt.Sprintf("%s may require a refill soon.", refillMedName)
		result.Response = &msg
	} else {
		msg := "One or more medications may require refill soon. Please review current medications."
		result.Response = &msg
	}

	// --- Ensure required root fields ---
	if result.SourceType == "" {
		result.SourceType = "current_medications"
	}
}
