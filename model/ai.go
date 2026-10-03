package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

/*

field Name string `json:"name" api:"required"`
field Arguments string `json:"arguments" api:"required"`

*/

type AI struct {
	Client openai.Client
	Tools  []responses.ToolUnionParam
}

type AIResponse struct {
	FunctionName string         `json:"name" api:"required"`
	Arguments    map[string]any `json:"arguments" api:"required"`
}

type AIReslut struct {
	Actions []AIResponse
	Text    string
}

func NewAI() (*AI, error) {
	open_ai_key := os.Getenv("OPENAI_API_KEY")
	if open_ai_key == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient()

	readParameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"target": map[string]any{
				"type":        "string",
				"description": "The target to read from the current observation.",
			},
		},
		"required":             []string{"target"},
		"additionalProperties": false,
	}

	clickParameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"target": map[string]any{
				"type":        "string",
				"description": "The target to click on the current observation.",
			},
			"x": map[string]any{
				"type":        "number",
				"description": "The x coordinate to click on the target.",
			},
			"y": map[string]any{
				"type":        "number",
				"description": "The y coordinate to click on the target.",
			},
		},
		"required":             []string{"target", "x", "y"},
		"additionalProperties": false,
	}

	typeParameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"value": map[string]any{
				"type":        "string",
				"description": "The text to type into the target.",
			},
		},
		"required":             []string{"target", "value"},
		"additionalProperties": false,
	}

	readTool := responses.ToolParamOfFunction("read", readParameters, true)
	clickTool := responses.ToolParamOfFunction("click", clickParameters, true)
	typeTool := responses.ToolParamOfFunction("type", typeParameters, true)

	readTool.OfFunction.Description = openai.String(
		"Request a read of a field in the current observation.",
	)
	clickTool.OfFunction.Description = openai.String(
		"Request a click on a target in the current observation.",
	)
	typeTool.OfFunction.Description = openai.String(
		"Enter text into the current focused input",
	)

	tools := []responses.ToolUnionParam{readTool, clickTool, typeTool}
	return &AI{
		Client: client,
		Tools:  tools,
	}, nil
}

// func Read(target string) {
// 	fmt.Println("Reading:", target)
// }

// func TypeText(target, value string) {
// 	fmt.Printf("Target: %s, Value: %s\n", target, value)
// }

func (ai *AI) GetActions(ctx context.Context, imageURL string, goal string) (AIReslut, error) {

	prompt := fmt.Sprintf(`
		Goal: %s

		Examine the current screenshot.
		Choose the next action toward the goal using the provided tools.
		If the goal is already achieved, report that.
	`, goal)

	resp, err := ai.Client.Responses.New(ctx, responses.ResponseNewParams{
		Model: "gpt-6-astra",
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentParamOfInputText(prompt),
						{
							OfInputImage: &responses.ResponseInputImageParam{
								Detail:   responses.ResponseInputImageDetailAuto,
								ImageURL: openai.String(imageURL),
							},
						},
					},
					responses.EasyInputMessageRoleUser,
				),
			},
		},
		Tools: ai.Tools,
	})

	if err != nil {
		return AIReslut{}, err
	}

	result := AIReslut{
		Text: resp.OutputText(),
	}

	// var actions []AIResponse
	for _, item := range resp.Output {
		if item.Type != "function_call" {
			continue
		}

		call := item.AsFunctionCall()

		action := AIResponse{
			FunctionName: call.Name,
		}

		err := json.Unmarshal([]byte(call.Arguments), &action.Arguments)
		if err != nil {
			return AIReslut{}, fmt.Errorf("decode arguments for %s: %w", call.Name, err)
		}

		result.Actions = append(result.Actions, action)
	}

	return result, nil

}
