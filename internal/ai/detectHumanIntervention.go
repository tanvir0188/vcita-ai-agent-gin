package ai

import (
	"context"
	"fmt"
	"os"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

func HumanInterventionNeeded(
	incomingMessage string,
) (string, error) {

	apiKey := os.Getenv("OPEN_AI_KEY")

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
