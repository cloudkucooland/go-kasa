package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloudkucooland/go-kasa"
)

const (
	minCardWidth = 24
	maxCardWidth = 34
	cardHeight   = 5
	columnGap    = 2
)

var (
	headerStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	nameStyle           = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75"))
	selectedNameStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Background(lipgloss.Color("24"))
	onStyle             = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	offStyle            = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("244"))
	metadataStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	borderStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	selectedBorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
	mutedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	errorStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("203"))
	labelStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

func (m model) View() string {
	if m.width == 0 {
		return "Starting Kasa…"
	}
	if m.renaming && m.selected != nil {
		return m.renameView()
	}
	if m.selected != nil {
		return m.detailView(*m.selected)
	}

	cols := m.columns()
	cardWidth := m.cardWidth()
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s\n", headerStyle.Render("KASA DEVICES"), metadataStyle.Render(fmt.Sprintf("%d shown · %d found · %d columns", len(m.visible), len(m.items), cols)))
	if m.searching {
		fmt.Fprintf(&b, "%s %s\n", headerStyle.Render("/"), m.searchInput.View())
	} else if m.searchInput.Value() != "" {
		fmt.Fprintf(&b, "%s\n", metadataStyle.Render("Filter: "+m.searchInput.Value()+"   (/ edit · esc clear)"))
	} else if m.isScanning {
		b.WriteString(mutedStyle.Render("Scanning devices…") + "\n")
	} else {
		b.WriteByte('\n')
	}

	first := m.rowOffset * cols
	rows := m.visibleRows()
	last := min(first+rows*cols, len(m.visible))
	for rowStart := first; rowStart < last; rowStart += cols {
		rowEnd := min(rowStart+cols, last)
		cards := make([][]string, rowEnd-rowStart)
		for idx := rowStart; idx < rowEnd; idx++ {
			cards[idx-rowStart] = m.renderCard(m.items[m.visible[idx]], idx == m.cursor, cardWidth)
		}
		for line := 0; line < cardHeight; line++ {
			for col, card := range cards {
				if col > 0 {
					b.WriteString(strings.Repeat(" ", columnGap))
				}
				b.WriteString(card[line])
			}
			b.WriteByte('\n')
		}
	}

	if len(m.visible) == 0 && !m.isScanning {
		if m.searchInput.Value() != "" {
			b.WriteString("No devices match this filter. Press esc to clear it.\n")
		} else {
			b.WriteString("No devices found. Press r to scan again.\n")
		}
	}
	if m.err != nil {
		fmt.Fprintf(&b, "%s\n", errorStyle.Render(formatOperationError(m.err)))
	} else if m.status != "" {
		fmt.Fprintf(&b, "%s\n", onStyle.Render(m.status))
	}
	b.WriteString(mutedStyle.Render("↑↓←→ / hjkl move   / search   enter details   t toggle   r scan   q quit"))
	return strings.TrimRight(b.String(), "\n")
}

func (m model) columns() int {
	width := m.width
	if width <= 0 {
		width = 80
	}
	cols := (width + columnGap) / (minCardWidth + columnGap)
	return max(1, cols)
}

func (m model) cardWidth() int {
	available := m.width - columnGap*(m.columns()-1)
	return max(1, min(maxCardWidth, available/m.columns()))
}

func (m model) visibleRows() int {
	rows := (m.height - 5) / cardHeight
	return max(1, rows)
}

func (m model) renderCard(item deviceItem, selected bool, width int) []string {
	inner := max(0, width-2)
	title := item.title()
	if selected {
		title = "> " + title
	}
	modelLine := item.info.Model + " · " + item.addr + " " + formatRSSIBadge(item.info.RSSI)
	statusLine := item.state()
	if len(item.info.Children) > 0 {
		active := 0
		aliases := make([]string, 0, len(item.info.Children))
		for _, child := range item.info.Children {
			if child.RelayState != 0 {
				active++
			}
			aliases = append(aliases, child.Alias)
		}
		statusLine = fmt.Sprintf("%d/%d on · %s", active, len(item.info.Children), strings.Join(aliases, ", "))
	}
	border := borderStyle
	name := nameStyle
	if selected {
		border = selectedBorderStyle
		name = selectedNameStyle
	}
	stateStyle := offStyle
	if item.isOn() {
		stateStyle = onStyle
	}
	if len(item.info.Children) > 0 {
		stateStyle = metadataStyle
		for _, child := range item.info.Children {
			if child.RelayState != 0 {
				stateStyle = onStyle
				break
			}
		}
	}
	lines := []string{
		border.Render("┌" + strings.Repeat("─", inner) + "┐"),
		border.Render("│") + name.Width(inner).Render(fit(title, inner)) + border.Render("│"),
		border.Render("│") + metadataStyle.Width(inner).Render(fit(modelLine, inner)) + border.Render("│"),
		border.Render("│") + stateStyle.Width(inner).Render(fit(statusLine, inner)) + border.Render("│"),
		border.Render("└" + strings.Repeat("─", inner) + "┘"),
	}
	return lines
}

func fit(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) > width {
		if width == 1 {
			runes = []rune("…")
		} else {
			runes = append(runes[:width-1], '…')
		}
	}
	return string(runes) + strings.Repeat(" ", width-len(runes))
}

func (m model) detailView(item deviceItem) string {
	title := item.title()
	state := item.state()
	stateStyle := offStyle
	if item.isOn() {
		stateStyle = onStyle
	}
	var detail string
	signalStr := renderRSSI(item.info.RSSI)
	if item.child != nil {
		detail = fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s", nameStyle.Render(title), labelStyle.Render("Device:"), item.info.Alias, labelStyle.Render("Outlet ID:"), item.childID, labelStyle.Render("Address:"), item.addr, labelStyle.Render("Model:"), item.info.Model, labelStyle.Render("Signal:"), signalStr, labelStyle.Render("State:"), stateStyle.Render(state))
	} else {
		detail = fmt.Sprintf("%s\n\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s", nameStyle.Render(title), labelStyle.Render("Address:"), item.addr, labelStyle.Render("Model:"), item.info.Model, labelStyle.Render("State:"), stateStyle.Render(state), labelStyle.Render("Firmware:"), item.info.SWVersion, labelStyle.Render("Signal:"), signalStr)
		if item.supportsBrightness() {
			bar := brightnessBar(item.info.Brightness)
			detail += fmt.Sprintf("\n%s %s", labelStyle.Render("Brightness:"), onStyle.Render(fmt.Sprintf("%s %d%%", bar, item.info.Brightness)))
		}
	}
	if len(item.info.Children) > 0 {
		detail += "\n\n" + headerStyle.Render("OUTLETS")
		for idx, child := range item.info.Children {
			prefix := "  "
			if idx == item.childIndex {
				prefix = "> "
			}
			childStyle := offStyle
			if child.RelayState != 0 {
				childStyle = onStyle
			}
			detail += fmt.Sprintf("\n%s%s  %s", prefix, nameStyle.Render(child.Alias), childStyle.Render(relayLabel(child.RelayState)))
		}
	}
	detail += "\n\n" + headerStyle.Render("LIVE ENERGY")
	switch {
	case item.child == nil && len(item.info.Children) > 0:
		detail += "\n" + mutedStyle.Render("Select an outlet to read its meter.")
	case item.emeterLoading:
		detail += "\n" + mutedStyle.Render("Reading meter…")
	case item.emeterError != nil:
		detail += "\n" + mutedStyle.Render("Energy data unavailable: "+formatOperationError(item.emeterError))
	case item.emeter != nil:
		volts := float64(item.emeter.VoltageMV) / 1000
		watts := float64(item.emeter.PowerMW) / 1000
		detail += fmt.Sprintf("\n%s %s    %s %s", labelStyle.Render("Voltage:"), nameStyle.Render(fmt.Sprintf("%.2f V", volts)), labelStyle.Render("Power:"), onStyle.Render(fmt.Sprintf("%.2f W", watts)))
	default:
		detail += "\n" + mutedStyle.Render("No meter reading")
	}
	detail += "\n\n" + headerStyle.Render(time.Now().Format("January 2006")+" ENERGY")
	if item.emeterLoading {
		detail += "\n" + mutedStyle.Render("Loading monthly history…")
	} else if item.monthError != nil {
		detail += "\n" + mutedStyle.Render("Monthly history unavailable: "+formatOperationError(item.monthError))
	} else if item.emeterMonth != nil {
		totalWH := uint64(0)
		for _, day := range item.emeterMonth.List {
			totalWH += uint64(day.WH)
		}
		detail += fmt.Sprintf("\n%s %s    %s %d days", labelStyle.Render("Month to date:"), onStyle.Render(fmt.Sprintf("%.2f kWh", float64(totalWH)/1000)), labelStyle.Render("Recorded:"), len(item.emeterMonth.List))
		if len(item.emeterMonth.List) > 0 {
			detail += "\n" + mutedStyle.Render("Daily usage (Wh)")
			for i := 0; i < len(item.emeterMonth.List); i += 2 {
				left := formatEnergyDay(item.emeterMonth.List[i])
				right := ""
				if i+1 < len(item.emeterMonth.List) {
					right = formatEnergyDay(item.emeterMonth.List[i+1])
				}
				detail += fmt.Sprintf("\n  %-19s %s", left, right)
			}
		}
	} else {
		detail += "\n" + mutedStyle.Render("No monthly history loaded; press r to refresh.")
	}
	detail += "\n\n" + headerStyle.Render("COUNTDOWN TIMERS")
	switch {
	case item.rulesLoading:
		detail += "\n" + mutedStyle.Render("Loading countdown timers…")
	case item.rulesError != nil:
		detail += "\n" + mutedStyle.Render("Countdown rules unavailable: "+formatOperationError(item.rulesError))
	case len(item.countdownRules) > 0:
		for _, r := range item.countdownRules {
			status := "Disabled"
			if r.Enable != 0 {
				status = fmt.Sprintf("Enabled (delay: %ds, active: %d)", r.Delay, r.Active)
			}
			detail += fmt.Sprintf("\n  %s: %s", nameStyle.Render(r.Name), labelStyle.Render(status))
		}
	default:
		detail += "\n" + mutedStyle.Render("No countdown timers")
	}

	detail += "\n\n" + headerStyle.Render("SCHEDULE RULES")
	switch {
	case item.rulesLoading:
		detail += "\n" + mutedStyle.Render("Loading schedule rules…")
	case item.rulesError != nil:
		detail += "\n" + mutedStyle.Render("Schedule rules unavailable: "+formatOperationError(item.rulesError))
	case len(item.scheduleRules) > 0:
		for _, s := range item.scheduleRules {
			status := "Disabled"
			if s.Enable != 0 {
				status = "Enabled"
			}
			timeStr := fmt.Sprintf("%02d:%02d", s.SMin/60, s.SMin%60)
			detail += fmt.Sprintf("\n  %s [%s]: %s", nameStyle.Render(s.Name), timeStr, labelStyle.Render(status))
		}
	default:
		detail += "\n" + mutedStyle.Render("No schedule rules")
	}
	if m.err != nil {
		detail += "\n\n" + errorStyle.Render(formatOperationError(m.err))
	} else if m.status != "" {
		detail += "\n\n" + onStyle.Render(m.status)
	}
	controls := "↑↓/j/k or pgup/pgdown scroll   t toggle   c 30m timer   x clear timers   n rename   r refresh   esc back   q quit"
	if len(item.info.Children) > 0 {
		controls = "↑↓ / j/k choose outlet   pgup/pgdown scroll   t toggle   c 30m timer   x clear timers   n rename   r meter   esc back   q quit"
	} else if item.supportsBrightness() {
		controls = "← / - dimmer   + / → brighter   t toggle   c 30m timer   x clear timers   n rename   r refresh   esc back   q quit"
	}

	lines := strings.Split(detail, "\n")
	maxHeight := max(10, m.height-3)
	if maxHeight <= 0 {
		maxHeight = 20
	}

	totalLines := len(lines)
	maxScroll := max(0, totalLines-maxHeight)
	offset := min(m.detailScrollOffset, maxScroll)
	if offset < 0 {
		offset = 0
	}

	end := min(offset+maxHeight, totalLines)
	visibleLines := lines[offset:end]

	var b strings.Builder
	if offset > 0 {
		b.WriteString(mutedStyle.Render(fmt.Sprintf("  [↑ %d more lines above]", offset)) + "\n")
	}
	b.WriteString(strings.Join(visibleLines, "\n"))
	if end < totalLines {
		b.WriteString("\n" + mutedStyle.Render(fmt.Sprintf("  [↓ %d more lines below]", totalLines-end)))
	}
	b.WriteString("\n\n" + mutedStyle.Render(controls))
	return b.String()
}

func formatEnergyDay(day kasa.EmeterDay) string {
	return fmt.Sprintf("%02d/%02d: %d Wh", day.Month, day.Day, day.WH)
}

func formatOperationError(err error) string {
	message := err.Error()
	lower := strings.ToLower(message)
	if strings.Contains(lower, "not support") || strings.Contains(lower, "unsupported") || strings.Contains(lower, "not exist") || strings.Contains(lower, "unknown method") {
		return "Unsupported by device: " + message
	}
	var networkErr net.Error
	if errors.Is(err, context.DeadlineExceeded) || errors.As(err, &networkErr) {
		return "Network error contacting device: " + message
	}
	return "Device request failed: " + message
}

func relayLabel(state uint) string {
	if state != 0 {
		return "On"
	}
	return "Off"
}

func (m model) renameView() string {
	target := "device"
	if m.selected.child != nil {
		target = "outlet"
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", headerStyle.Render("Rename "+target), metadataStyle.Render("Current name: "+m.selected.title()), m.renameInput.View(), mutedStyle.Render("enter save   esc cancel"))
}

func brightnessBar(value uint) string {
	filled := min(10, int(value+9)/10)
	return strings.Repeat("▰", filled) + strings.Repeat("▱", 10-filled)
}

var (
	signalStrong = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	signalGood   = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	signalWeak   = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

func renderRSSI(rssi int) string {
	var style lipgloss.Style
	var label string
	switch {
	case rssi >= -50:
		style = signalStrong
		label = "Strong"
	case rssi >= -70:
		style = signalGood
		label = "Good"
	default:
		style = signalWeak
		label = "Weak"
	}
	return style.Render(fmt.Sprintf("%d dBm (%s)", rssi, label))
}

func formatRSSIBadge(rssi int) string {
	switch {
	case rssi >= -50:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("●")
	case rssi >= -70:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("●")
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("●")
	}
}
