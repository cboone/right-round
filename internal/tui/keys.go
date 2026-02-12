package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Enter    key.Binding
	Back     key.Binding
	Tab      key.Binding
	Search   key.Binding
	Copy     key.Binding
	Help     key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "go to top"),
	),
	End: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "go to bottom"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "l"),
		key.WithHelp("enter/l", "expand"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "h"),
		key.WithHelp("esc/h", "back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch type"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy JSON"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) shortHelp() string {
	return helpStyle.Render(
		key.NewBinding(key.WithHelp("up/down", "navigate")).Help().Key + " " + key.NewBinding(key.WithHelp("up/down", "navigate")).Help().Desc +
			"  " + k.Enter.Help().Key + " " + k.Enter.Help().Desc +
			"  " + k.Tab.Help().Key + " " + k.Tab.Help().Desc +
			"  " + k.Search.Help().Key + " " + k.Search.Help().Desc +
			"  " + k.Copy.Help().Key + " " + k.Copy.Help().Desc +
			"  " + k.Help.Help().Key + " " + k.Help.Help().Desc +
			"  " + k.Quit.Help().Key + " " + k.Quit.Help().Desc,
	)
}

func (k keyMap) fullHelp() string {
	return helpStyle.Render(
		k.Up.Help().Key + " " + k.Up.Help().Desc +
			"  " + k.Down.Help().Key + " " + k.Down.Help().Desc +
			"  " + k.PageUp.Help().Key + " " + k.PageUp.Help().Desc +
			"  " + k.PageDown.Help().Key + " " + k.PageDown.Help().Desc +
			"  " + k.Home.Help().Key + " " + k.Home.Help().Desc +
			"  " + k.End.Help().Key + " " + k.End.Help().Desc +
			"\n" +
			k.Enter.Help().Key + " " + k.Enter.Help().Desc +
			"  " + k.Back.Help().Key + " " + k.Back.Help().Desc +
			"  " + k.Tab.Help().Key + " " + k.Tab.Help().Desc +
			"  " + k.Search.Help().Key + " " + k.Search.Help().Desc +
			"  " + k.Copy.Help().Key + " " + k.Copy.Help().Desc +
			"  " + k.Quit.Help().Key + " " + k.Quit.Help().Desc,
	)
}
