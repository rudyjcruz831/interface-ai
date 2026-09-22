package main



import (
	"fmt"
	"encoding/json"
	"os"
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

func typeText(target, value string){
	fmt.Printf("Target: %s, Value: %s\n", target, value)
}

func usreInput(){

}

func main() {
	// Target names are simplified for this exercise.
	// steps := []Step{
	// 	{Action: "type", Target: "Memmber id field", Value: "12344"},
	// 	{Action: "click", Target: "Search button", Value: ""},
	// 	{Action: "read", Target: "Balance field", Value: ""},
	// }

	getAI()

	// Run from this file's folder; workflow.json is two folders above it.
	data, err := os.ReadFile("../../workflow.json")
	if err != nil {
		fmt.Println("Error reading workflow.json:", err)
		return
	}

	var steps []Step
	err = json.Unmarshal(data, &steps)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Visit the saved steps in order. The underscore ignores the list index.
	for _, step := range steps {
		switch step.Action {
		case "click":
			click(step.Target)
		// TODO: Add a case for the "read" action here.
		case "read":
			read(step.Target)
		case "type":
			typeText(step.Target, step.Value)

		default:
			fmt.Println("Unsupported action:", step.Action)
			return
		}
	}
}
