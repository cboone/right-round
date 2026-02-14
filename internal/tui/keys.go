package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	PageUp    key.Binding
	PageDown  key.Binding
	Home      key.Binding
	End       key.Binding
	Left      key.Binding
	Right     key.Binding
	Enter     key.Binding
	Back      key.Binding
	Tab       key.Binding
	Search    key.Binding
	PrevGroup key.Binding
	NextGroup key.Binding
	Sort      key.Binding
	Verbose   key.Binding
	Copy      key.Binding
	Help      key.Binding
	Quit      key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "move down"),
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
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "focus groups"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "focus entries"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "l"),
		key.WithHelp("enter/l", "open detail"),
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
		key.WithHelp("/", "filter"),
	),
	PrevGroup: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "previous group"),
	),
	NextGroup: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next group"),
	),
	Sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "group sort"),
	),
	Verbose: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "detail mode"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy JSON"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) shortHelp() string {
	return helpStyle.Render(
		k.Up.Help().Key + " " + k.Up.Help().Desc +
			"  " + k.Left.Help().Key + "/" + k.Right.Help().Key + " panes" +
			"  " + k.PrevGroup.Help().Key + "/" + k.NextGroup.Help().Key + " groups" +
			"  " + k.Search.Help().Key + " " + k.Search.Help().Desc +
			"  " + k.Enter.Help().Key + " " + k.Enter.Help().Desc +
			"  " + k.Tab.Help().Key + " " + k.Tab.Help().Desc +
			"  " + k.Copy.Help().Key + " " + k.Copy.Help().Desc +
			"  " + k.Help.Help().Key + " more" +
			"  " + k.Quit.Help().Key + " " + k.Quit.Help().Desc,
	)
}

func (k keyMap) fullHelp() string {
	return helpStyle.Render(
		k.Up.Help().Key + " " + k.Up.Help().Desc +
			"  " + k.Down.Help().Key + " " + k.Down.Help().Desc +
			"  " + k.Left.Help().Key + " " + k.Left.Help().Desc +
			"  " + k.Right.Help().Key + " " + k.Right.Help().Desc +
			"  " + k.PrevGroup.Help().Key + " " + k.PrevGroup.Help().Desc +
			"  " + k.NextGroup.Help().Key + " " + k.NextGroup.Help().Desc +
			"  " + k.PageUp.Help().Key + " " + k.PageUp.Help().Desc +
			"  " + k.PageDown.Help().Key + " " + k.PageDown.Help().Desc +
			"  " + k.Home.Help().Key + " " + k.Home.Help().Desc +
			"  " + k.End.Help().Key + " " + k.End.Help().Desc +
			"\n" +
			k.Enter.Help().Key + " " + k.Enter.Help().Desc +
			"  " + k.Back.Help().Key + " " + k.Back.Help().Desc +
			"  " + k.Sort.Help().Key + " " + k.Sort.Help().Desc +
			"  " + k.Verbose.Help().Key + " " + k.Verbose.Help().Desc +
			"  " + k.Tab.Help().Key + " " + k.Tab.Help().Desc +
			"  " + k.Search.Help().Key + " " + k.Search.Help().Desc +
			"  " + k.Copy.Help().Key + " " + k.Copy.Help().Desc +
			"  " + k.Quit.Help().Key + " " + k.Quit.Help().Desc,
	)
}
