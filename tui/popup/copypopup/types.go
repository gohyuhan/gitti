package copypopup

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/utils"
)

const (
	ShortHash = "short_hash"
	Upstream  = "upstream"
)

type Option struct {
	Label  string
	Value  string
	Ready  bool
	Failed bool
}

func (o Option) FilterValue() string { return o.Label }

type optionDelegate struct{}

func (optionDelegate) Height() int                             { return 1 }
func (optionDelegate) Spacing() int                            { return 0 }
func (optionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (optionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	option, ok := item.(Option)
	if !ok {
		return
	}
	value := utils.EscapeControlCharacters(option.Value)
	if !option.Ready {
		value = i18n.LANGUAGEMAPPING.CopyLoading
		if option.Failed {
			value = i18n.LANGUAGEMAPPING.CopyUnavailable
		}
	}
	line := ansi.Truncate(option.Label+": "+value, max(0, m.Width()-constant.ListItemOrTitleWidthPad), "...")
	if index == m.Index() {
		fmt.Fprint(w, style.SelectedItemStyle.Render("❯ "+line))
	} else {
		fmt.Fprint(w, style.ItemStyle.Render("  "+line))
	}
}

type Model struct {
	ID      uint64
	Options list.Model
}

func (m *Model) Resolve(kind, value string, err error) {
	index := 1 // Both deferred values occupy the second row.
	if len(m.Options.Items()) <= index {
		return
	}
	if kind == Upstream && err == nil && value == "" {
		if m.Options.Index() == index {
			m.Options.Select(0)
		}
		m.Options.RemoveItem(index)
		return
	}
	option, ok := m.Options.Items()[index].(Option)
	if !ok {
		return
	}
	option.Value = value
	option.Ready = err == nil && value != ""
	option.Failed = !option.Ready
	_ = m.Options.SetItem(index, option)
}
