[![Build Status](https://travis-ci.org/adayoung/ada-bot.svg?branch=master)](https://travis-ci.org/adayoung/ada-bot)
[![Go report](https://goreportcard.com/badge/adayoung/ada-bot)](https://goreportcard.com/report/adayoung/ada-bot)

# ada-bot
Hi! This is ada-bot, a robot for Achaea's Discord channel. It's very basic at the moment and usually sticks around to look cute and call out deathsightings over a dedicated channel.

## How to use

 * Create an App on [Discord API](https://discordapp.com/developers/applications/me) and add a Bot user to it.
 * Make sure to enable Privileged Gateway Intents for your bot, we need both, Presence and Server Membership.
 * Install the package with `go get github.com/adayoung/ada-bot`
 * Copy the file `config.yaml` from the src somewhere
 * Copy the Bot's Authorization token and put it in config.yaml
 * Copy your dedicated channel's ID and put it in config.yaml
 * Put your IRE game of choice's API URL in config.yaml
 * Build and install with `go install github.com/adayoung/ada-bot`
 * Navigate to wherever you have your config.yaml and launch the bot with `$GOPATH/bin/ada-bot`
 * There's a [systemd service](https://www.freedesktop.org/software/systemd/man/systemd.service.html) file included as well to let systemd manage the bot!

## Adding it to your server

 * Copy the client ID of your bot from [Discord API](https://discordapp.com/developers/applications/me) and construct an invitation URL in the form:  
    `https://discordapp.com/api/oauth2/authorize?client_id=<BOT CLIENT ID>&scope=bot&permissions=3072`  
    You can read more about the above [here](https://discordapp.com/developers/docs/topics/oauth2#adding-bots-to-guilds)
 * The bot would need at least 'Read Messages' and 'Send Messages' permissions to be configured for it. You can configure it in its role or setup permissions via the integration by modifying the 'permission' parameter in the above URL. Here's an unofficial [Permissions Calculator](https://discordapi.com/permissions.html).

## IRE Event Logging

 * The bot is capable of logging events received from IRE's gamefeed, it'll require an account with a PostgreSQL server with credentials noted in config.yaml

## Building with Docker
Use the following command to build with the latest version of Go (=>1.15):  
 * `cd <path to repository>`
 * `docker run --rm -v $PWD:/usr/src/ada-bot -w /usr/src/ada-bot -v $GOPATH:/go -e "CGO_ENABLED=0" golang:1.15 go build -v`
