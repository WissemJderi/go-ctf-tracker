package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormFieldType int

const (
	FieldText FormFieldType = iota
	FieldSelect
	FieldBool
	FieldTextArea
)

type FormField struct {
	Key         string
	Label       string
	Type        FormFieldType
	Placeholder string
	Options     []string // For FieldSelect
	SelectIdx   int      // Active index for FieldSelect
	BoolValue   bool     // Active value for FieldBool
	TextInput   textinput.Model
	TextArea    textarea.Model
}

type FormModel struct {
	fields       []FormField
	focusIndex   int
	width        int
	height       int
	onSubmit     func(map[string]string, map[string]bool) tea.Cmd
	onCancel     func() tea.Cmd
	errorMessage string
}

func NewFormModel(fields []FormField, onSubmit func(map[string]string, map[string]bool) tea.Cmd, onCancel func() tea.Cmd) FormModel {
	for i := range fields {
		switch fields[i].Type {
		case FieldText:
			ti := textinput.New()
			ti.Placeholder = fields[i].Placeholder
			ti.CharLimit = 156
			ti.Width = 40
			if fields[i].Key == "points" {
				ti.CharLimit = 10
			}
			fields[i].TextInput = ti
		case FieldTextArea:
			ta := textarea.New()
			ta.Placeholder = fields[i].Placeholder
			ta.SetWidth(50)
			ta.SetHeight(5)
			ta.ShowLineNumbers = false
			fields[i].TextArea = ta
		}
	}

	// Focus first field
	if len(fields) > 0 {
		fields[0].focus(true)
	}

	return FormModel{
		fields:     fields,
		focusIndex: 0,
		onSubmit:   onSubmit,
		onCancel:   onCancel,
	}
}

func (f *FormField) focus(focused bool) tea.Cmd {
	var cmd tea.Cmd
	if focused {
		switch f.Type {
		case FieldText:
			cmd = f.TextInput.Focus()
		case FieldTextArea:
			cmd = f.TextArea.Focus()
		}
	} else {
		switch f.Type {
		case FieldText:
			f.TextInput.Blur()
		case FieldTextArea:
			f.TextArea.Blur()
		}
	}
	return cmd
}

func (fm *FormModel) SetValues(textVals map[string]string, boolVals map[string]bool) {
	for i := range fm.fields {
		f := &fm.fields[i]
		if f.Type == FieldBool {
			f.BoolValue = boolVals[f.Key]
		} else if f.Type == FieldSelect {
			f.SelectIdx = 0
			val := textVals[f.Key]
			for idx, opt := range f.Options {
				if strings.EqualFold(opt, val) {
					f.SelectIdx = idx
					break
				}
			}
		} else if f.Type == FieldText {
			f.TextInput.SetValue(textVals[f.Key])
		} else if f.Type == FieldTextArea {
			f.TextArea.SetValue(textVals[f.Key])
		}
	}
}

func (fm *FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if fm.onCancel != nil {
				return *fm, fm.onCancel()
			}
			return *fm, nil

		case "tab", "down":
			// Unfocus current
			cmd := fm.fields[fm.focusIndex].focus(false)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			// Move to next (including Submit button as the last index)
			fm.focusIndex = (fm.focusIndex + 1) % (len(fm.fields) + 1)

			// Focus next
			if fm.focusIndex < len(fm.fields) {
				cmd = fm.fields[fm.focusIndex].focus(true)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}

		case "shift+tab", "up":
			// Unfocus current
			if fm.focusIndex < len(fm.fields) {
				cmd := fm.fields[fm.focusIndex].focus(false)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}

			// Move to prev
			fm.focusIndex--
			if fm.focusIndex < 0 {
				fm.focusIndex = len(fm.fields) // Submit button
			}

			// Focus prev
			if fm.focusIndex < len(fm.fields) {
				cmd := fm.fields[fm.focusIndex].focus(true)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}

		case "enter":
			if fm.focusIndex == len(fm.fields) {
				// Submit button clicked
				return *fm, fm.submit()
			}

			// For fields:
			currentField := &fm.fields[fm.focusIndex]
			if currentField.Type == FieldBool {
				currentField.BoolValue = !currentField.BoolValue
			} else if currentField.Type == FieldSelect {
				currentField.SelectIdx = (currentField.SelectIdx + 1) % len(currentField.Options)
			} else if currentField.Type == FieldTextArea {
				// Allow enter in text area, pass along
			} else {
				// For text, enter moves to next field
				cmd := currentField.focus(false)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				fm.focusIndex = (fm.focusIndex + 1) % (len(fm.fields) + 1)
				if fm.focusIndex < len(fm.fields) {
					cmd = fm.fields[fm.focusIndex].focus(true)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			}

		case "left", "right":
			currentField := &fm.fields[fm.focusIndex]
			if currentField.Type == FieldBool {
				currentField.BoolValue = !currentField.BoolValue
			} else if currentField.Type == FieldSelect {
				if msg.String() == "left" {
					currentField.SelectIdx--
					if currentField.SelectIdx < 0 {
						currentField.SelectIdx = len(currentField.Options) - 1
					}
				} else {
					currentField.SelectIdx = (currentField.SelectIdx + 1) % len(currentField.Options)
				}
			}
		}
	}

	// Update active component
	if fm.focusIndex < len(fm.fields) {
		f := &fm.fields[fm.focusIndex]
		if f.Type == FieldText {
			var cmd tea.Cmd
			f.TextInput, cmd = f.TextInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if f.Type == FieldTextArea {
			var cmd tea.Cmd
			f.TextArea, cmd = f.TextArea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return *fm, tea.Batch(cmds...)
}

func (fm *FormModel) submit() tea.Cmd {
	// Validate points if present
	fm.errorMessage = ""
	textVals := make(map[string]string)
	boolVals := make(map[string]bool)

	for _, f := range fm.fields {
		switch f.Type {
		case FieldText:
			val := strings.TrimSpace(f.TextInput.Value())
			if f.Key == "points" && val != "" {
				if _, err := strconv.Atoi(val); err != nil {
					fm.errorMessage = "Points must be a valid number!"
					return nil
				}
			}
			textVals[f.Key] = sanitizeString(val)
		case FieldTextArea:
			textVals[f.Key] = sanitizeString(f.TextArea.Value())
		case FieldSelect:
			textVals[f.Key] = f.Options[f.SelectIdx]
		case FieldBool:
			boolVals[f.Key] = f.BoolValue
		}
	}

	// Basic validation
	if textVals["ctf_name"] == "" {
		fm.errorMessage = "CTF Name is required!"
		return nil
	}
	if textVals["name"] == "" {
		fm.errorMessage = "Challenge Name is required!"
		return nil
	}

	if fm.onSubmit != nil {
		return fm.onSubmit(textVals, boolVals)
	}
	return nil
}

func (fm FormModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	for i, f := range fm.fields {
		isFocused := i == fm.focusIndex
		var labelStyle lipgloss.Style
		if isFocused {
			labelStyle = lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)
		} else {
			labelStyle = lipgloss.NewStyle().Foreground(ColorHighlight)
		}

		s.WriteString(labelStyle.Render(fmt.Sprintf("%-16s ", f.Label)))

		switch f.Type {
		case FieldText:
			s.WriteString(f.TextInput.View())
		case FieldTextArea:
			s.WriteString("\n" + f.TextArea.View())
		case FieldSelect:
			var optViews []string
			for idx, opt := range f.Options {
				if idx == f.SelectIdx {
					optViews = append(optViews, lipgloss.NewStyle().
						Background(ColorPrimary).
						Foreground(ColorHighlight).
						Padding(0, 1).
						Render(sanitizeString(opt)))
				} else {
					optViews = append(optViews, lipgloss.NewStyle().
						Foreground(ColorMuted).
						Padding(0, 1).
						Render(sanitizeString(opt)))
				}
			}
			s.WriteString(strings.Join(optViews, " | "))
			if isFocused {
				s.WriteString(" " + lipgloss.NewStyle().Foreground(ColorMuted).Render("(Left/Right to select)"))
			}
		case FieldBool:
			var valStr string
			if f.BoolValue {
				valStr = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true).Render("[🔥 YES]")
			} else {
				valStr = lipgloss.NewStyle().Foreground(ColorMuted).Render("[  NO ]")
			}
			s.WriteString(valStr)
			if isFocused {
				s.WriteString(" " + lipgloss.NewStyle().Foreground(ColorMuted).Render("(Enter/Left/Right to toggle)"))
			}
		}
		s.WriteString("\n\n")
	}

	// Submit Button
	isSubmitFocused := fm.focusIndex == len(fm.fields)
	var btnStyle lipgloss.Style
	if isSubmitFocused {
		btnStyle = lipgloss.NewStyle().
			Background(ColorSuccess).
			Foreground(lipgloss.Color("#111")).
			Bold(true).
			Padding(0, 3)
	} else {
		btnStyle = lipgloss.NewStyle().
			Background(ColorMuted).
			Foreground(ColorHighlight).
			Padding(0, 3)
	}
	s.WriteString(btnStyle.Render("SUBMIT") + "   ")

	// Cancel hint
	s.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render("Press ESC to cancel / Tab to navigate"))

	if fm.errorMessage != "" {
		s.WriteString("\n\n" + lipgloss.NewStyle().Foreground(ColorDanger).Bold(true).Render("Error: "+fm.errorMessage))
	}

	return s.String()
}

// sanitizeString removes potentially problematic characters from user input
func sanitizeString(s string) string {
	// Remove null characters and ensure no control characters
	runes := make([]rune, 0, len(s))
	for _, r := range s {
		if r != 0 && (r < 32 || (r >= 127 && r <= 159)) {
			continue // Skip control characters
		}
		runes = append(runes, r)
	}
	return string(runes)
}
