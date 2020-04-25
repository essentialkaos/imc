package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v11/fmtutil"
	"pkg.re/essentialkaos/ek.v11/options"
	"pkg.re/essentialkaos/ek.v11/sortutil"
	"pkg.re/essentialkaos/ek.v11/timeutil"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	ic "pkg.re/essentialkaos/go-icecast.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type GUI struct {
	banner      *widgets.Paragraph
	serverInfo  *widgets.Table
	streamsInfo *widgets.Table
	sourcesInfo *widgets.Table

	tw int
	th int

	sourceCount  int
	listenerPeak int

	lastError error
}

var (
	greyStyle  = ui.Style{8, -1, ui.ModifierClear}
	boldStyle  = ui.Style{-1, -1, ui.ModifierBold}
	onAirStyle = ui.Style{10, -1, ui.ModifierClear}
	alertStyle = ui.Style{1, -1, ui.ModifierBold}
)

var (
	serverHeader  = []string{"ID", "STARTED", "SOURCES", "LISTENERS", "LISTENER PEAK", "BANNED IP", "OUT NOW", "IN TOTAL", "OUT TOTAL"}
	streamsHeader = []string{"MOUNT", "LISTENERS", "LISTENER PEAK", "SLOW LISTENERS", "CONNECTIONS", "IN NOW", "OUT NOW", "IN TOTAL", "OUT TOTAL"}
	sourcesHeader = []string{"MOUNT", "SOURCE IP", "AIR TIME", "BITRATE", "SAMPLE RATE", "CHANNELS", "TRACK"}
)

var stats *ic.Stats

var gui *GUI

// ////////////////////////////////////////////////////////////////////////////////// //

func renderGUI() error {
	err := ui.Init()

	if err != nil {
		return err
	}

	defer ui.Close()

	gui = NewGUI()

	uiEvents := ui.PollEvents()
	interval := time.Duration(options.GetI(OPT_INTERVAL)) * time.Second
	ticker := time.NewTicker(interval).C

	fetchStats()
	gui.Render()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "u":
				fetchStats()
				gui.Render()
			case "<Resize>":
				ui.Clear()
				gui.Render()
			}
		case <-ticker:
			fetchStats()
			gui.Render()
		}
	}

	return nil
}

func fetchStats() {
	newStats, err := icecast.GetStats()

	if err == nil {
		stats = newStats
	}

	gui.lastError = err
}

func formatServerStats() []string {
	var result []string

	result = append(result, stats.Info.ID)
	result = append(result, timeutil.Format(stats.Started, "%Y/%m/%d %H:%M"))
	result = append(result, fmtutil.PrettyNum(stats.Stats.Sources))
	result = append(result, fmtutil.PrettyNum(stats.Stats.Listeners))
	result = append(result, fmtutil.PrettyNum(countListenerPeak(stats)))
	result = append(result, fmtutil.PrettyNum(stats.Stats.BannedIPs))
	result = append(result, fmt.Sprintf("%s/s", fmtutil.PrettySize(stats.Stats.OutgoingBitrate)))
	result = append(result, fmtutil.PrettySize(stats.Stats.StreamBytesRead))
	result = append(result, fmtutil.PrettySize(stats.Stats.StreamBytesSent))

	return result
}

func formatStreamsStats() [][]string {
	var result [][]string

	result = append(result, streamsHeader)

	for _, mount := range getSources(stats) {
		var sourceInfo []string

		source := stats.Sources[mount]

		sourceInfo = append(sourceInfo, mount)
		sourceInfo = append(sourceInfo, fmtutil.PrettyNum(source.Stats.Listeners))
		sourceInfo = append(sourceInfo, fmtutil.PrettyNum(source.Stats.ListenerPeak))
		sourceInfo = append(sourceInfo, fmtutil.PrettyNum(source.Stats.SlowListeners))
		sourceInfo = append(sourceInfo, fmtutil.PrettyNum(source.Stats.ListenerConnections))
		sourceInfo = append(sourceInfo, fmt.Sprintf("%s/s", fmtutil.PrettySize(source.Stats.IncomingBitrate)))
		sourceInfo = append(sourceInfo, fmt.Sprintf("%s/s", fmtutil.PrettySize(source.Stats.OutgoingBitrate)))
		sourceInfo = append(sourceInfo, fmtutil.PrettySize(source.Stats.TotalBytesRead))
		sourceInfo = append(sourceInfo, fmtutil.PrettySize(source.Stats.TotalBytesSent))

		result = append(result, formatRowData(sourceInfo))
	}

	return result
}

func formatSourcesStats() [][]string {
	var result [][]string

	result = append(result, sourcesHeader)

	for _, mount := range getSources(stats) {
		var sourceInfo []string

		source := stats.Sources[mount]

		sourceInfo = append(sourceInfo, mount)
		sourceInfo = append(sourceInfo, source.SourceIP)
		sourceInfo = append(sourceInfo, formatDuration(time.Since(source.StreamStarted)))
		sourceInfo = append(sourceInfo, fmt.Sprintf("%s/s", fmtutil.PrettySize(source.AudioInfo.Bitrate)))
		sourceInfo = append(sourceInfo, fmtutil.PrettyNum(source.AudioInfo.SampleRate))
		sourceInfo = append(sourceInfo, fmtutil.PrettyNum(source.AudioInfo.Channels))

		if source.Track.RawInfo != "" {
			sourceInfo = append(sourceInfo, source.Track.RawInfo)
		} else {
			sourceInfo = append(sourceInfo, "—")
		}

		result = append(result, formatRowData(sourceInfo))
	}

	return result
}

func getSources(s *ic.Stats) []string {
	var result []string

	for mount := range s.Sources {
		result = append(result, mount)
	}

	sortutil.StringsNatural(result)

	return result
}

func formatDuration(d time.Duration) string {
	dur := int(d.Seconds())

	var hours, minutes, seconds int

	for i := 0; i < 3; i++ {
		switch {
		case dur > 3600:
			hours = dur / 3600
			dur = dur % 3600
		case dur > 60:
			minutes = dur / 60
			dur = dur % 60
		default:
			seconds = dur
		}
	}

	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

func countListenerPeak(stats *ic.Stats) int {
	var result int

	for _, source := range stats.Sources {
		result += source.Stats.ListenerPeak
	}

	return result
}

func formatRowData(data []string) []string {
	for i, v := range data {
		data[i] = " " + v
	}

	return data
}

// ////////////////////////////////////////////////////////////////////////////////// //

func NewGUI() *GUI {
	g := &GUI{}

	g.tw, g.th = ui.TerminalDimensions()

	g.banner = widgets.NewParagraph()
	g.serverInfo = g.createTable("Server")
	g.streamsInfo = g.createTable("Streams")
	g.sourcesInfo = g.createTable("Sources")

	g.serverInfo.Rows = [][]string{
		formatRowData(serverHeader),
		formatRowData([]string{"—", "—", "—", "—", "—", "—", "—", "—"}),
	}

	g.streamsInfo.Rows = [][]string{
		formatRowData(streamsHeader),
		formatRowData([]string{"—", "—", "—", "—", "—", "—", "—", "—", "—"}),
	}

	g.sourcesInfo.Rows = [][]string{
		formatRowData(sourcesHeader),
		formatRowData([]string{"—", "—", "—", "—", "—", "—", "—"}),
	}
	g.sourcesInfo.ColumnWidths = []int{20, 18, 12, 16, 14, 10, -1}

	return g
}

func (g *GUI) Render() {
	g.Update()

	if stats != nil {
		if g.sourceCount != stats.Stats.Sources {
			g.sourceCount = stats.Stats.Sources
			ui.Clear()
		}

		if stats.Stats.Sources != 0 {
			ui.Render(g.banner, g.serverInfo, g.streamsInfo, g.sourcesInfo)
		} else {
			g.streamsInfo.Rows, g.sourcesInfo.Rows = nil, nil
			ui.Render(g.banner, g.serverInfo)
		}
	} else {
		ui.Render(g.banner)
	}
}

func (g *GUI) Update() {
	g.tw, g.th = ui.TerminalDimensions()

	g.updateBanner()

	if stats != nil {
		g.updateServerInfo()
		g.updateStreamsInfo()
		g.updateSourcesInfo()
	}
}

func (g *GUI) createTable(title string) *widgets.Table {
	t := widgets.NewTable()
	t.BorderStyle = greyStyle

	t.Title = " " + title + " "
	t.TitleStyle = greyStyle
	t.FillRow = true
	t.RowStyles[0] = boldStyle

	return t
}

func (g *GUI) updateServerInfo() {
	g.serverInfo.SetRect(0, 3, g.tw, 8)
	g.serverInfo.Rows[1] = formatRowData(formatServerStats())
}

func (g *GUI) updateStreamsInfo() {
	y := 9 + (stats.Stats.Sources * 2) + 3
	g.streamsInfo.SetRect(0, 9, g.tw, y)
	g.streamsInfo.Rows = formatStreamsStats()
}

func (g *GUI) updateSourcesInfo() {
	y1 := 10 + (stats.Stats.Sources * 2) + 3
	y2 := y1 + (stats.Stats.Sources * 2) + 3
	g.sourcesInfo.SetRect(0, y1, g.tw, y2)
	g.sourcesInfo.Rows = formatSourcesStats()
}

func (g *GUI) updateBanner() {
	var text string

	if g.lastError != nil {
		g.banner.TextStyle = alertStyle
		g.banner.BorderStyle = alertStyle
		text = g.lastError.Error()
	} else {
		g.banner.TextStyle = boldStyle
		g.banner.BorderStyle = greyStyle
		text = fmt.Sprintf("Connected to Icecast on %s", host)
	}

	prefixSize := (g.tw - len(text)) / 2

	g.banner.SetRect(0, 0, g.tw, 3)

	if prefixSize > 1 {
		g.banner.Text = strings.Repeat(" ", prefixSize) + text
	} else {
		g.banner.Text = text
	}

	if stats == nil || len(stats.Sources) == 0 {
		g.banner.BorderStyle = greyStyle
	} else {
		g.banner.BorderStyle = onAirStyle
	}
}
