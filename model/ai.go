package model

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

/*

field Name string `json:"name" api:"required"`
field Arguments string `json:"arguments" api:"required"`

*/

type AIResponse struct {
	FunctionName string         `json:"name" api:"required"`
	Arguments    map[string]any `json:"arguments" api:"required"`
	// Text         string
}

func GetAI(imageURL string, goal string) {

	open_ai_key := os.Getenv("OPENAI_API_KEY")
	if open_ai_key == "" {
		panic("OPENAI_API_KEY environment variable is not set")
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
		},
		"required":             []string{"target"},
		"additionalProperties": false,
	}

	typeParameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"target": map[string]any{
				"type":        "string",
				"description": "The target to type text into in the current observation.",
			},
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
		"Request to type text into a target in the current observation.",
	)

	// resp, err := client.Responses.New(context.TODO(), responses.ResponseNewParams{
	// 	Model: "gpt-6-astra",
	// 	Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(`
	// 	Task: Read the available balance for account 17229.

	// 	Current screen:
	// 	Account Details for account 17229.
	// 	Visible fields: Account type, Balance, Available balance.

	// 	Use the read tool to request the field you need.
	// 	`)},
	// 	Tools: []responses.ToolUnionParam{readTool, clickTool, typeTool},
	// })

	prompt := fmt.Sprintf(`
		Goal: %s

		Examine the current screenshot.
		Choose the next action toward the goal using the provided tools.
		If required information is missing, explain what is missing.
		If the goal is already achieved, report that.
	`, goal)

	resp, err := client.Responses.New(context.TODO(), responses.ResponseNewParams{
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
		Tools: []responses.ToolUnionParam{readTool, clickTool, typeTool},
	})

	if err != nil {
		panic(err.Error())
	}

	// var aiResponseList []AIResponse

	for _, item := range resp.Output {
		if item.Type == "function_call" {
			call := item.AsFunctionCall()
			fmt.Println("Function:", call.Name)
			fmt.Println("Arguments:", call.Arguments)
		}
	}
	fmt.Println("----------------------------------------------------")
	fmt.Println("Text response:", resp.OutputText())
}
