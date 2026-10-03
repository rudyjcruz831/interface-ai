package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/rudyjcruz831/interface-ai/model"
)

// Practice runner: these functions print actions; they do not control a real UI.
// Exercise: add the missing "read" branch in the switch below.

// Step describes one instruction in a saved workflow.
type Step struct {
	Action string
	Target string
	Value  string
}

func main() {

	// get image form chormadp and convert to base64
	ctx := context.Background()
	browser := model.NewBrowser(ctx)
	defer browser.Close()

	ai, err := model.NewAI()
	if err != nil {
		fmt.Println("Error initializing AI:", err)
		return
	}

	errNav := browser.Navigate("https://parabank.parasoft.com/parabank/index.htm")
	if errNav != nil {
		fmt.Println("Error navigating to URL:", errNav)
		return
	}

	imagesBytes, err := browser.Screenshot()
	if err != nil {
		fmt.Println("Error capturing screenshot:", err)
		return
	}

	// encodeToBase64 is a helper function to encode bytes to base64f.
	encodeToBase64 := func(b []byte) string {
		return base64.StdEncoding.EncodeToString(b)
	}

	imageURL := "data:image/png;base64," + encodeToBase64(imagesBytes)
	// Call the GetAI function to initialize the AI client and tools.

	result, err := ai.GetActions(context.Background(), imageURL, "Click inside the username input field")
	if err != nil {
		fmt.Println("Error getting AI response:", err)
		return
	}

	if len(result.Actions) == 0 {
		fmt.Println("No action returned")

		if result.Text != "" {
			fmt.Println("Model message:", result.Text)
		}
		return
	}

	action := result.Actions[0]

	switch action.FunctionName {
	case "click":
		target, targetOk := action.Arguments["target"].(string)
		x, xOk := action.Arguments["x"].(float64)
		y, yOk := action.Arguments["y"].(float64)

		if !targetOk || !xOk || !yOk {
			fmt.Println("Invalid click arguments : expected target, x, and y")
			return
		}

		if target == "" {
			fmt.Println("Click target is empty")
			return
		}

		// These bounds match the viewport configured in Navigate().
		if x < 0 || x > 1920 || y < 0 || y > 1080 {
			fmt.Println("Click coordinates are outside the viewport.")
			return
		}

		err := browser.Click(target, x, y)
		if err != nil {
			fmt.Println("Error performing click:", err)
			return
		}

	default:
		fmt.Println("Unknown action:", action.FunctionName)
	}

	afterBytes, err := browser.Screenshot()
	if err != nil {
		fmt.Println("Error capturing screenshot after action:", err)
		return
	}

	errWrite := os.WriteFile("after-click.png", afterBytes, 0600)
	if errWrite != nil {
		fmt.Println("Error saving screenshot after action:", errWrite)
		return
	}
	fmt.Println("Saved after-click.png")
}
