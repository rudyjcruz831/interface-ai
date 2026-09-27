package main

import (
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

func click(target string) {
	fmt.Println("Clicking:", target)
}

func read(target string) {
	fmt.Println("Reading:", target)
}

func typeText(target, value string) {
	fmt.Printf("Target: %s, Value: %s\n", target, value)
}

func main() {
	// Target names are simplified for this exercise.
	// steps := []Step{
	// 	{Action: "type", Target: "Memmber id field", Value: "12344"},
	// 	{Action: "click", Target: "Search button", Value: ""},
	// 	{Action: "read", Target: "Balance field", Value: ""},
	// }

	imagesBytes, err := os.ReadFile("images/p.png")
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}

	// encodeToBase64 is a helper function to encode bytes to base64.
	encodeToBase64 := func(b []byte) string {
		return base64.StdEncoding.EncodeToString(b)
	}

	imageURL := "data:image/png;base64," + encodeToBase64(imagesBytes)
	// Call the GetAI function to initialize the AI client and tools.

	model.GetAI(imageURL)

	// Run from this file's folder; workflow.json is two folders above it.
	// data, err := os.ReadFile("workflows/workflow.json")
	// if err != nil {
	// 	fmt.Println("Error reading workflow.json:", err)
	// 	return
	// }

	// var steps []Step
	// err = json.Unmarshal(data, &steps)
	// if err != nil {
	// 	fmt.Println("Error parsing JSON:", err)
	// 	return
	// }

	// // Visit the saved steps in order. The underscore ignores the list index.
	// for _, step := range steps {
	// 	switch step.Action {
	// 	case "click":
	// 		click(step.Target)
	// 	// TODO: Add a case for the "read" action here.
	// 	case "read":
	// 		read(step.Target)
	// 	case "type":
	// 		typeText(step.Target, step.Value)

	// 	default:
	// 		fmt.Println("Unsupported action:", step.Action)
	// 		return
	// 	}
	// }

}
