package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudkucooland/go-kasa"
)

type discoveryMsg map[string]*kasa.Sysinfo
type actionMsg string
type renameMsg struct {
	addr    string
	childID string
	alias   string
}
type emeterMsg struct {
	addr     string
	childID  string
	reading  *kasa.EmeterRealtime
	month    *kasa.EmeterDaystat
	err      error
	monthErr error
}

func scanDevices() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		devices, err := kasa.BroadcastDiscovery(ctx, 3)
		if err != nil {
			return errMsg(err)
		}
		return discoveryMsg(devices)
	}
}

func toggleDevice(item deviceItem) tea.Cmd {
	return func() tea.Msg {
		d, err := kasa.NewDevice(item.addr)
		if err != nil {
			return errMsg(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s, err := d.GetSettingsCtx(ctx)
		if err != nil {
			return errMsg(err)
		}

		var state bool
		if item.childID != "" {
			childID, child, err := resolveChild(s.Children, item.childID, item.title())
			if err != nil {
				return errMsg(err)
			}
			state = child.RelayState == 0
			err = d.SetRelayStateChildCtx(ctx, childID, state)
		} else {
			state = s.RelayState == 0
			err = d.SetRelayStateCtx(ctx, state)
		}
		if err != nil {
			return errMsg(err)
		}
		return actionMsg("Device toggled")
	}
}

func setBrightness(addr string, brightness int) tea.Cmd {
	return func() tea.Msg {
		if brightness < 1 {
			brightness = 1
		}
		if brightness > 100 {
			brightness = 100
		}
		d, err := kasa.NewDevice(addr)
		if err != nil {
			return errMsg(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := d.SetBrightnessCtx(ctx, brightness); err != nil {
			return errMsg(err)
		}
		return actionMsg(fmt.Sprintf("Brightness set to %d%%", brightness))
	}
}

func setAlias(item deviceItem, alias string) tea.Cmd {
	return func() tea.Msg {
		d, err := kasa.NewDevice(item.addr)
		if err != nil {
			return errMsg(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if item.childID != "" {
			settings, err := d.GetSettingsCtx(ctx)
			if err != nil {
				return errMsg(err)
			}
			childID, _, err := resolveChild(settings.Children, item.childID, item.title())
			if err != nil {
				return errMsg(err)
			}
			if err := d.SetChildAliasCtx(ctx, childID, alias); err != nil {
				return errMsg(err)
			}
		} else if err := d.SetAliasCtx(ctx, alias); err != nil {
			return errMsg(err)
		}
		return renameMsg{addr: item.addr, childID: item.childID, alias: alias}
	}
}

func fetchEmeter(item deviceItem) tea.Cmd {
	return func() tea.Msg {
		d, err := kasa.NewDevice(item.addr)
		if err != nil {
			return emeterMsg{addr: item.addr, childID: item.childID, err: err}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var reading *kasa.EmeterRealtime
		var month *kasa.EmeterDaystat
		var monthErr error
		now := time.Now()
		if item.childID != "" {
			settings, settingsErr := d.GetSettingsCtx(ctx)
			if settingsErr != nil {
				return emeterMsg{addr: item.addr, childID: item.childID, err: settingsErr}
			}
			childID, _, resolveErr := resolveChild(settings.Children, item.childID, item.title())
			if resolveErr != nil {
				return emeterMsg{addr: item.addr, childID: item.childID, err: resolveErr}
			}
			reading, err = d.GetEmeterChildCtx(ctx, childID)
			if err == nil {
				month, monthErr = d.GetEmeterChildMonthCtx(ctx, int(now.Month()), now.Year(), childID)
			}
		} else {
			reading, err = d.GetEmeterCtx(ctx)
			if err == nil {
				month, monthErr = d.GetEmeterMonthCtx(ctx, int(now.Month()), now.Year())
			}
		}
		return emeterMsg{addr: item.addr, childID: item.childID, reading: reading, month: month, err: err, monthErr: monthErr}
	}
}

func (m model) Init() tea.Cmd { return scanDevices() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.searchInput.SetValue("")
				m.searchInput.Blur()
				m.searching = false
				m.applyFilter()
				return m, nil
			case "enter":
				m.searchInput.Blur()
				m.searching = false
				return m, nil
			default:
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.applyFilter()
				return m, cmd
			}
		}
		if m.renaming {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.renaming = false
				m.renameInput.Blur()
				return m, nil
			case "enter":
				alias := strings.TrimSpace(m.renameInput.Value())
				m.renaming = false
				m.renameInput.Blur()
				if alias == "" {
					m.err = fmt.Errorf("alias cannot be empty")
					return m, nil
				}
				m.status = "Contacting device to save alias…"
				return m, setAlias(*m.selected, alias)
			default:
				var cmd tea.Cmd
				m.renameInput, cmd = m.renameInput.Update(msg)
				return m, cmd
			}
		}
		m.err = nil
		m.status = ""
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			if m.selected != nil {
				if m.selected.child == nil && len(m.selected.info.Children) > 0 {
					return m, m.focusChild(1)
				}
				m.selected.emeterLoading = true
				m.selected.emeterError = nil
				m.status = "Contacting device for energy data…"
				return m, fetchEmeter(*m.selected)
			}
			m.isScanning = true
			m.status = "Scanning for devices…"
			return m, scanDevices()
		case "n":
			if m.selected != nil {
				m.renameInput.SetValue(m.selected.title())
				m.renameInput.CursorEnd()
				m.renaming = true
				return m, m.renameInput.Focus()
			}
		case "/":
			if m.selected == nil {
				m.searchInput.CursorEnd()
				m.searching = true
				return m, m.searchInput.Focus()
			}
		case "esc":
			if m.selected != nil {
				m.selected = nil
			} else if m.searchInput.Value() != "" {
				m.searchInput.SetValue("")
				m.applyFilter()
			}
		case "enter":
			if m.selected == nil && len(m.visible) > 0 {
				item := m.items[m.visible[m.cursor]]
				m.selected = &item
				if len(item.info.Children) == 0 && !item.emeterLoaded {
					m.selected.emeterLoading = true
					m.status = "Contacting device for live and historical energy…"
					return m, fetchEmeter(item)
				}
			}
		case "t":
			if m.selected != nil {
				m.status = "Contacting device to change power…"
				return m, toggleDevice(*m.selected)
			}
			if len(m.visible) > 0 {
				item := m.items[m.visible[m.cursor]]
				m.status = "Contacting " + item.title() + "…"
				return m, toggleDevice(item)
			}
		case "+", "=":
			if m.selected != nil && m.selected.supportsBrightness() {
				m.status = "Contacting device to set brightness…"
				return m, setBrightness(m.selected.addr, int(m.selected.info.Brightness)+5)
			}
		case "-":
			if m.selected != nil && m.selected.supportsBrightness() {
				m.status = "Contacting device to set brightness…"
				return m, setBrightness(m.selected.addr, int(m.selected.info.Brightness)-5)
			}
		case "left", "h":
			if m.selected == nil {
				m.moveCursor(-1)
			} else if m.selected.supportsBrightness() {
				m.status = "Contacting device to set brightness…"
				return m, setBrightness(m.selected.addr, int(m.selected.info.Brightness)-5)
			}
		case "right", "l":
			if m.selected == nil {
				m.moveCursor(1)
			} else if m.selected.supportsBrightness() {
				m.status = "Contacting device to set brightness…"
				return m, setBrightness(m.selected.addr, int(m.selected.info.Brightness)+5)
			}
		default:
			switch msg.String() {
			case "up", "k":
				if m.selected != nil {
					return m, m.focusChild(-1)
				}
				m.moveCursor(-m.columns())
			case "down", "j":
				if m.selected != nil {
					return m, m.focusChild(1)
				}
				m.moveCursor(m.columns())
			}
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.searchInput.Width = max(10, min(40, msg.Width-4))
		m.keepCursorVisible()
	case discoveryMsg:
		previousSelection := m.selected
		m.items = make([]deviceItem, 0, len(msg))
		for addr, info := range msg {
			m.items = append(m.items, deviceItem{addr: addr, info: info})
		}
		sort.Slice(m.items, func(i, j int) bool {
			if m.items[i].info.Alias != m.items[j].info.Alias {
				return m.items[i].info.Alias < m.items[j].info.Alias
			}
			if m.items[i].title() != m.items[j].title() {
				return m.items[i].title() < m.items[j].title()
			}
			if m.items[i].addr != m.items[j].addr {
				return m.items[i].addr < m.items[j].addr
			}
			return m.items[i].childID < m.items[j].childID
		})
		m.selected = nil
		if previousSelection != nil {
			for i := range m.items {
				root := m.items[i]
				if root.addr != previousSelection.addr {
					continue
				}
				if previousSelection.childID != "" {
					childID, child, err := resolveChild(root.info.Children, previousSelection.childID, previousSelection.title())
					if err == nil {
						root.child = &child
						root.childID = childID
					}
				}
				root.emeter = previousSelection.emeter
				root.emeterMonth = previousSelection.emeterMonth
				root.emeterLoaded = previousSelection.emeterLoaded
				root.emeterLoading = previousSelection.emeterLoading
				root.emeterError = previousSelection.emeterError
				root.monthError = previousSelection.monthError
				m.selected = &root
				break
			}
		}
		if len(m.items) == 0 {
			m.cursor = 0
		} else if m.cursor >= len(m.items) {
			m.cursor = len(m.items) - 1
		}
		m.visible = nil
		m.cursor = 0
		m.applyFilter()
		m.isScanning = false
		if m.status == "Scanning for devices…" {
			m.status = ""
		} else if strings.HasSuffix(m.status, "; refreshing…") {
			m.status = strings.TrimSuffix(m.status, "; refreshing…")
		}
		if m.pendingRename != nil {
			pending := m.pendingRename
			m.pendingRename = nil
			confirmed, reported := false, "device not found"
			for _, item := range m.items {
				if item.addr != pending.addr {
					continue
				}
				if pending.childID == "" {
					reported = item.info.Alias
					confirmed = reported == pending.alias
				} else if _, child, err := resolveChild(item.info.Children, pending.childID, ""); err == nil {
					reported = child.Alias
					confirmed = reported == pending.alias
				}
				break
			}
			if confirmed {
				m.status = fmt.Sprintf("Rename confirmed: %s", pending.alias)
			} else {
				m.status = fmt.Sprintf("Rename not confirmed: device reports %q", reported)
			}
		}
	case emeterMsg:
		for i := range m.items {
			if m.items[i].addr == msg.addr && m.items[i].childID == msg.childID {
				m.items[i].emeter = msg.reading
				m.items[i].emeterMonth = msg.month
				m.items[i].emeterLoaded = true
				m.items[i].emeterLoading = false
				m.items[i].emeterError = msg.err
				m.items[i].monthError = msg.monthErr
				break
			}
		}
		if m.selected != nil && m.selected.addr == msg.addr && m.selected.childID == msg.childID {
			m.selected.emeter = msg.reading
			m.selected.emeterMonth = msg.month
			m.selected.emeterLoaded = true
			m.selected.emeterLoading = false
			m.selected.emeterError = msg.err
			m.selected.monthError = msg.monthErr
			m.status = ""
		}
	case actionMsg:
		m.status = string(msg) + "; refreshing…"
		m.isScanning = true
		return m, scanDevices()
	case renameMsg:
		m.pendingRename = &renameRequest{addr: msg.addr, childID: msg.childID, alias: msg.alias}
		m.status = "Alias saved; confirming with device…"
		m.isScanning = true
		return m, scanDevices()
	case errMsg:
		m.err = msg
		m.isScanning = false
		if m.pendingRename != nil {
			m.pendingRename = nil
			m.status = "Could not confirm rename from refreshed device list"
		}
	}
	return m, nil
}

func (m *model) focusChild(delta int) tea.Cmd {
	if m.selected == nil || len(m.selected.info.Children) == 0 {
		return nil
	}
	current := -1
	if m.selected.childID != "" {
		for i, child := range m.selected.info.Children {
			if child.ID == m.selected.childID || child.Alias == m.selected.title() {
				current = i
				break
			}
		}
	}
	next := current + delta
	if next < -1 {
		return nil
	}
	if next >= len(m.selected.info.Children) {
		return nil
	}
	if next == -1 {
		root := deviceItem{addr: m.selected.addr, info: m.selected.info}
		m.selected = &root
		return nil
	}
	child := m.selected.info.Children[next]
	item := deviceItem{addr: m.selected.addr, info: m.selected.info, child: &child, childID: child.ID, emeterLoading: true}
	m.selected = &item
	m.status = "Contacting device for live and historical energy…"
	return fetchEmeter(item)
}

func resolveChild(children []kasa.Child, requestedID, alias string) (string, kasa.Child, error) {
	for _, child := range children {
		if child.ID == requestedID {
			return child.ID, child, nil
		}
	}

	var idMatches []kasa.Child
	for _, child := range children {
		if child.ID != "" && requestedID != "" && (strings.HasSuffix(child.ID, requestedID) || strings.HasSuffix(requestedID, child.ID)) {
			idMatches = append(idMatches, child)
		}
	}
	if len(idMatches) == 1 {
		return idMatches[0].ID, idMatches[0], nil
	}

	var aliasMatches []kasa.Child
	for _, child := range children {
		if child.Alias == alias {
			aliasMatches = append(aliasMatches, child)
		}
	}
	if len(aliasMatches) == 1 {
		return aliasMatches[0].ID, aliasMatches[0], nil
	}
	if len(aliasMatches) > 1 || len(idMatches) > 1 {
		return "", kasa.Child{}, fmt.Errorf("outlet %q matches multiple device children", alias)
	}
	return "", kasa.Child{}, fmt.Errorf("device did not report outlet %q (ID %q)", alias, requestedID)
}

func (m *model) moveCursor(delta int) {
	if len(m.visible) == 0 {
		return
	}
	next := m.cursor + delta
	if next < 0 || next >= len(m.visible) {
		return
	}
	m.cursor = next
	m.keepCursorVisible()
}

func (m *model) applyFilter() {
	var current deviceItem
	hasCurrent := m.cursor >= 0 && m.cursor < len(m.visible)
	if hasCurrent {
		current = m.items[m.visible[m.cursor]]
	}

	query := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	m.visible = make([]int, 0, len(m.items))
	for i, item := range m.items {
		if item.matches(query) {
			m.visible = append(m.visible, i)
		}
	}

	m.cursor = 0
	if hasCurrent {
		for i, index := range m.visible {
			item := m.items[index]
			if item.addr == current.addr && item.childID == current.childID {
				m.cursor = i
				break
			}
		}
	}
	m.rowOffset = 0
	m.keepCursorVisible()
}

func (m *model) keepCursorVisible() {
	if len(m.visible) == 0 {
		m.rowOffset = 0
		return
	}
	cols := m.columns()
	visibleRows := m.visibleRows()
	row := m.cursor / cols
	if row < m.rowOffset {
		m.rowOffset = row
	} else if row >= m.rowOffset+visibleRows {
		m.rowOffset = row - visibleRows + 1
	}
	if m.rowOffset < 0 {
		m.rowOffset = 0
	}
}
