package ai

import (
	"fmt"
	"testing"
)

func TestHumanIntervention(t *testing.T) {

	result, err := HumanInterventionNeeded(
		"I want to schedule a meeting",
	)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}
