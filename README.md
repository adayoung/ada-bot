# ada-bot
Hi! This is ada-bot, a robot for Achaea's Discord channel. It's very basic at the moment and usually sticks around to look cute and call out deathsightings over a dedicated channel.

## Prerequisites

 * [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)
 * [DiscordGo](https://github.com/bwmarrin/discordgo)

## How to use

 * Create an App on [Discord API](https://discordapp.com/developers/applications/me) and add a Bot user to it.
 * Install the package with `go get github.com/adayoung/ada-bot`
 * Copy the file `env.yaml` from the src somewhere
 * Copy the Bot's Authorization token and put it in env.yaml
 * Copy your dedicated channel's ID and put it in env.yaml
 * Put your IRE game of choice's API URL in env.yaml
 * Build and install with `go install github.com/adayoung/ada-bot`
 * Launch with `$GOPATH/bin/ada-bot`
 * There's a [systemd service](https://www.freedesktop.org/software/systemd/man/systemd.service.html) file included as well to let systemd manage the bot!
