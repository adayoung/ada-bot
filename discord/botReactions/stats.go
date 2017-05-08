// Parts of this file is totally copied from https://github.com/iopred/bruxism/blob/master/statsplugin/statsplugin.go
// This applies: https://github.com/iopred/bruxism/blob/master/LICENSE
// Thanks to Mister Christopher Rhodes for this! :D

package botReactions

import (
	"bytes"
	"fmt"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var statsStartTime = time.Now()

func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

type stats struct {
	Trigger string
}

func (s *stats) Help() string {
	return "I have statistics!"
}

func (s *stats) HelpDetail() string {
	return s.Help()
}

func (s *stats) Reaction(m *discordgo.Message, a *discordgo.Member, mType string) Reaction {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)

	fmt.Fprintf(w, "```\n")
	fmt.Fprintf(w, "Discordgo: \t%s\n", discordgo.VERSION)
	fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
	fmt.Fprintf(w, "Uptime: \t%s\n", getDurationString(time.Now().Sub(statsStartTime)))
	fmt.Fprintf(w, "Memory used: \t%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Concurrent tasks: \t%d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "\n```")

	w.Flush()
	out := buf.String()
	return Reaction{Text: out}
}

func init() {
	_stats := &stats{
		Trigger: "stats",
	}
	addReaction(_stats.Trigger, "CREATE", _stats)
}
