package main

import (
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucas-ingemar/clergo/internal/components"
	"github.com/lucas-ingemar/clergo/internal/config"
	"github.com/lucas-ingemar/clergo/internal/io"
	"github.com/lucas-ingemar/clergo/internal/markdown"
	"github.com/lucas-ingemar/clergo/internal/shared"
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())

	appStyle = lipgloss.NewStyle().
		// Background(lipgloss.Color("#ffffff")).
		Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

func Init() {
	// TODO: Read cofig file

	// FIXME: Does not work
	os.Mkdir(config.CONFIG.LibPath, os.ModePerm)
	os.Mkdir(path.Join(config.CONFIG.LibPath, "notes"), os.ModePerm)
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	scrollPaperUp    key.Binding
	taEnter          key.Binding
	taExit           key.Binding
	deleteItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		scrollPaperUp: key.NewBinding(
			key.WithKeys("J"),
			key.WithHelp("J", "Scroll paper up"),
		),
		taEnter: key.NewBinding(
			key.WithKeys("enter"),
		),
		taExit: key.NewBinding(
			key.WithKeys("esc"),
		),
		deleteItem: key.NewBinding(
			key.WithKeys("x"),
		),
	}
}

type model struct {
	list list.Model
	// itemGenerator *randomItemGenerator
	keys         *listKeyMap
	delegateKeys *components.DelegateKeyMap
	textarea     textarea.Model
	textareaMode bool
	// viewportReady bool
}

func newModel() model {
	var (
		delegateKeys = components.NewDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// items := []list.Item{}
	items, err := io.ReadFiles()
	if err != nil {
		// FIXME: LOG LOG
		fmt.Println(err)
	}

	// Setup list
	delegate := components.NewItemDelegate(delegateKeys)
	itemsList := list.New(items, delegate, 0, 0)
	itemsList.Title = "Clergo"
	itemsList.Styles.Title = titleStyle
	itemsList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
			listKeys.scrollPaperUp,
		}
	}

	var initItem *shared.Item
	if len(items) > 0 {
		tmpItem := items[0].(shared.Item)
		initItem = &tmpItem
	}

	return model{
		list:         itemsList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		textarea:     components.NewTextarea(initItem),
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		listWidth := msg.Width/3 - h
		m.list.SetSize(listWidth, msg.Height-v)

		m.textarea.SetHeight(msg.Height - v)
		m.textarea.SetWidth(msg.Width - listWidth)

	case tea.KeyMsg:
		if m.textarea.Focused() {
			if key.Matches(msg, m.keys.taExit) {
				item := markdown.Parse(m.textarea.Value())
				m.list.SetItem(m.list.Index(), item)
				m.textarea.SetValue(components.BodyText(item))
				m.textarea.Blur()
				// FIXME: Error handler
				err := io.WriteFile(&item)
				if err != nil {
					return m, m.list.NewStatusMessage(statusMessageStyle(err.Error()))
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)

			return m, tea.Batch(cmd)
		}
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.deleteItem):
			err := io.DeleteFile(m.list.SelectedItem().(shared.Item))
			// FIXME: Error handler
			if err != nil {
				return m, m.list.NewStatusMessage(statusMessageStyle(err.Error()))
			}

		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.Remove.SetEnabled(true)
			newItem := shared.Item{
				TitleVar: "New",
				TagsVar:  []string{"tag1"},
				BodyText: "",
			}
			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added new item"))

			cmds = append(cmds, tea.Batch(insCmd, statusCmd))
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.list.KeyMap.CursorUp) || key.Matches(msg, m.list.KeyMap.CursorDown):
			m.textarea.SetValue(m.list.SelectedItem().(shared.Item).Body())

		case key.Matches(msg, m.keys.taEnter):
			m.textarea.SetValue(markdown.Generate(m.list.SelectedItem().(shared.Item)))
			m.textarea.Focus()
			return m, tea.Batch(cmds...)

		case key.Matches(msg, m.keys.insertItem):
			if m.list.FilterState() == list.Filtering {
				break
			}
			m.textarea.SetValue(markdown.Generate(m.list.SelectedItem().(shared.Item)))
			m.textarea.Focus()
			return m, tea.Batch(cmds...)
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.textarea.View()) // + "\n\n" + help
	// return appStyle.Render(m.list.View(), m.textarea.View())
}

func main() {
	// items, err := io.ReadFiles()
	// fmt.Println(err)
	// fmt.Println(items[0].(shared.Item).Title())
	// return

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
