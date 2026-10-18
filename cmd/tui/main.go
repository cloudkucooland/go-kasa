package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	input := textinput.New()
	input.Prompt = ""
	input.CharLimit = 64
	input.Width = 40
	search := textinput.New()
	search.Prompt = ""
	search.Placeholder = "name, model, IP, or outlet ID"
	search.CharLimit = 128
	search.Width = 40
	m := model{isScanning: true, renameInput: input, searchInput: search}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
