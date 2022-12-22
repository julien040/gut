package print

import (
	"fmt"

	printColor "github.com/fatih/color"
)

type Color string

const (
	Error   Color = "red"
	Warning Color = "yellow"
	Info    Color = "blue"
	Success Color = "green"
	None    Color = "none"
)

// # Message prints a message with a color
//
// Corresponding colors:
//
//   - Error: red
//
//   - Warning: yellow
//
//   - Info: blue
//
//   - Success: green
func Message(message string, color Color) {
	switch color {
	case Error:
		printColor.Red(message)
	case Warning:
		printColor.Yellow(message)
	case Info:
		printColor.Blue(message)
	case Success:
		printColor.Green(message)
	case None:
		fmt.Println(message)
	default:
		fmt.Println(message)
	}
}
