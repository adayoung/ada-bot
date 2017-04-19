# Bot Reactions

Hi! It's easy to implement simple bot reactions! We have a lovely dispatcher to take care of the gory details as long as your code conforms to the following example:

```Go
package bot_reactions

import (
	"github.com/bwmarrin/discordgo"
)

// A bot reaction implements the following interface
// type BotReaction interface {
// 	Help() string
// 	HelpDetail(*discordgo.Message) string
// 	Reaction(message *discordgo.Message, author *discordgo.Member) string
// }

type Ping struct {  
	Trigger string
}

// This should be a single line without \n breaks, it appears in !help
func (p *Ping) Help() string {
	return "Pong!"
}

// This is detailed help and it's supposed to be invoked when people do !help <trigger>
func (p *Ping) HelpDetail(m *discordgo.Message) string {
	return p.Help()
}

// The actual reaction, the full message content is available in m.Content
// The returned string is sent to the channel where the trigger was seen
func (p *Ping) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	return "Pong!"
}

// Here's where you add your 'trigger' and register your reaction with the dispatcher
func init() {
	ping := &Ping{
		Trigger: "ping",
	}
	addReaction(ping.Trigger, ping)
}
```

Have a look at the other reactions for more complex examples!
