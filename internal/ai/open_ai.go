package ai

import (
	"context"
	"encoding/json"
	"fmt"

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

	systemPrompt := `You are an appointment scheduling extractor.

	Analyze the message and determine if the platform AI is in the process of checking
	calendar availability to lock in an appointment time.

	This is TRUE when the message:
	- Mentions checking the calendar for a specific service
	- Echoes back a preferred time and/or backup time from the client
	- Has NOT yet confirmed a booked slot (that comes in a later message)

	Example of a message that needs scheduling:
	"I'm checking the calendar now for your LifePact: Initial Consultation on June 27.
	To confirm, you'd prefer 9:00 AM (Central) and your backup is 5:00 PM (Central)—correct?"

	Example of a message that does NOT need scheduling (already done):
	"You're all set for June 27 at 9:00 AM (Central)."

	Example of a message that does NOT need scheduling (unrelated):
	"Thank you for reaching out! How can I help you today?"

	Extract:
	- needs_scheduling: boolean.(set's to true if it needs scheduling, false otherwise)
	- service_name: strip duration/channel suffix. "LifePact: Initial Consultation (15 min, phone)" → "LifePact: Initial Consultation". "LifePact: Lab and Protocol Review (30 min)"→"LifePact: Lab and Protocol Review"
	- preferred_time: the client's first choice, ISO 8601 with timezone offset
	- backup_time: the client's second choice if mentioned, ISO 8601 with timezone offset. null if not mentioned.
	- start_time: the client's first choice, ISO 8601 with timezone offset.
	- end_time: based on the service name, you can get the duration. Add the duration and get the end_time in ISO 8601 with timezone offset.

	Return ONLY valid JSON, no markdown:
	{
		"needs_scheduling": true,
		"service_name": "LifePact: Initial Consultation",
		"preferred_time": "2026-06-27T09:00:00-05:00",
		"backup_time": "2026-06-27T17:00:00-05:00",
		"start_time":"2026-06-27T09:00:00-05:00",
		"end_time":"2026-06-27T09:00:00-05:00"
	}

	If needs_scheduling is false, return all other fields as null/empty.`
	
	fullInput := fmt.Sprintf(
		"System prompt:%s\nUser message:\n%s",
		systemPrompt, incomingMessage,
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
