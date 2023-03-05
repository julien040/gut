package print

import (
	"github.com/charmbracelet/lipgloss"
	printColor "github.com/fatih/color"
)

type Color string

const (
	Error    Color = "red"
	Warning  Color = "yellow"
	Info     Color = "blue"
	Success  Color = "green"
	Optional Color = "optional"
	None     Color = "none"
)

var optionalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("248"))

var infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

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
func Message(message string, color Color, a ...interface{}) {
	switch color {
	case Error:
		printColor.Red(message, a...)
	case Warning:
		printColor.Yellow(message, a...)
	case Info:
		printColor.Blue(infoStyle.Render(message), a...)
	case Success:
		printColor.Green(message, a...)
	case Optional:
		printColor.White(optionalStyle.Render(message), a...)
	case None:
		printColor.White(message, a...)
	default:
		printColor.White(message, a...)
	}
}
