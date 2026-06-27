package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
)

func HumanInterventionNeeded(
	incomingMessage string,
) (string, error) {

	apiKey := config.Envs.OpenAPIKey

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	systemPrompt := fmt.Sprintf(`
You are a classifier. You will be given a FAQ and a user question.

FAQ:
%s

Your only job is to determine whether the user's question needs a human to respond.

Rules:
- Casual/greeting messages (e.g. "hi", "hello", "how are you", "thanks") → return: {"human_intervention_needed": false}
- General non-critical questions that can be answered from the FAQ → return: {"human_intervention_needed": false, "escalation_reason":"reason in one sentence"}
- Questions outside the FAQ but non-critical and general in nature → return: {"human_intervention_needed": false, "escalation_reason":"reason in one sentence"}
- Questions outside the FAQ that are critical (medical advice, emergencies, dosage changes, serious side effects, sensitive health decisions) → return: {"human_intervention_needed": true, "escalation_reason":"reason in one sentence"}

Return ONLY valid JSON. No explanation, no extra text, nothing else.
`, DocString)

	fullInput := fmt.Sprintf(
		"System:\n%s\n\nUser:\n%s",
		systemPrompt,
		incomingMessage,
	)

	resp, err := client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT4oMini,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(fullInput),
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.OutputText(), nil
}

func ConfirmedAppointmentDate(incomingMessage string) (*AppointmentAiResponse, error) {
	apiKey := config.Envs.OpenAPIKey

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	var systemPrompt string = ConfirmedAppointmentDateSystemPrompt

	fullInput := fmt.Sprintf(
		"Today's date: %s \nSystem prompt: %s \nUser message: \n%s",
		time.Now().UTC().Format("2006-01-02"), systemPrompt, incomingMessage,
	)

	resp, err := client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT4oMini,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(fullInput),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("openai request failed: %w", err)
	}

	outputText := resp.OutputText()
	if outputText == "" {
		return nil, fmt.Errorf("empty response from AI")
	}

	var result AppointmentAiResponse
	if err := json.Unmarshal([]byte(outputText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON response: %w\nResponse was: %s", err, outputText)
	}

	return &result, nil
}

func PredictMedicationRefill(medicationNote string) (*MedicationRefillResponse, error) {
	apiKey := config.Envs.OpenAPIKey
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	currentDate := time.Now().UTC().Format("2006-01-02")

	fullInput := fmt.Sprintf(
		"CURRENT_DATE: %s\n\nMEDICATION_NOTE:\n%s",
		currentDate, medicationNote,
	)

	resp, err := client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT4oMini,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(fullInput),
			},
			Instructions: openai.String(MedicationReminderSystemPrompt),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("openai request failed: %w", err)
	}

	outputText := resp.OutputText()
	if outputText == "" {
		return nil, fmt.Errorf("empty response from AI")
	}

	var result MedicationRefillResponse
	if err := json.Unmarshal([]byte(outputText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON response: %w\nResponse was: %s", err, outputText)
	}

	// Enforce deterministic reminder logic, same as _validate_and_fix in ai.py.
	validateAndFixRefill(&result, currentDate)

	return &result, nil
}
