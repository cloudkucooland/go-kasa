package main

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/cloudkucooland/go-kasa"
)

type deviceItem struct {
	addr          string
	info          *kasa.Sysinfo
	child         *kasa.Child
	childID       string
	childIndex    int
	emeter        *kasa.EmeterRealtime
	emeterMonth   *kasa.EmeterDaystat
	emeterLoaded  bool
	emeterLoading bool
	emeterError   error
	monthError    error
}

func (d deviceItem) matches(query string) bool {
	if query == "" {
		return true
	}
	searchable := strings.ToLower(strings.Join([]string{
		d.info.Alias,
		d.title(),
		d.info.Model,
		d.addr,
		d.childID,
	}, " "))
	if strings.Contains(searchable, query) {
		return true
	}
	for _, child := range d.info.Children {
		if strings.Contains(strings.ToLower(child.Alias+" "+child.ID), query) {
			return true
		}
	}
	return false
}

func (d deviceItem) title() string {
	if d.child != nil {
		return d.child.Alias
	}
	return d.info.Alias
}

func (d deviceItem) state() string {
	if d.isOn() {
		if d.child == nil && d.info.Brightness > 0 {
			return "On · " + strconv.FormatUint(uint64(d.info.Brightness), 10) + "%"
		}
		return "On"
	}
	return "Off"
}

func (d deviceItem) isOn() bool {
	if d.child != nil {
		return d.child.RelayState != 0
	}
	return d.info.RelayState != 0
}

func (d deviceItem) supportsBrightness() bool {
	return d.child == nil && strings.HasPrefix(strings.ToUpper(d.info.Model), "HS220")
}

type model struct {
	items         []deviceItem
	visible       []int
	cursor        int
	selected      *deviceItem
	renameInput   textinput.Model
	renaming      bool
	searchInput   textinput.Model
	searching     bool
	width         int
	height        int
	rowOffset     int
	err           error
	status        string
	isScanning    bool
	pendingRename *renameRequest
	brightnessSeq int
}

type renameRequest struct {
	addr    string
	childID string
	alias   string
}

type errMsg error
