package bot_reactions

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/adayoung/ada-bot/settings"
)

type dice struct {
	Trigger string
}

func (d *dice) Help() string {
	return fmt.Sprintf("Roll a dice! DnD style, %sdice xdy+z", settings.Settings.Discord.BotPrefix)
}

func (d *dice) HelpDetail(m *discordgo.Message) string {
	return d.Help()
}

var diceRegexp = regexp.MustCompile(`(?i)([0-9]+)d([0-9]+)(?:\+([0-9]+))?`)

func (d *dice) Reaction(m *discordgo.Message, a *discordgo.Member) string {
	request := strings.TrimSpace(m.Content[len(settings.Settings.Discord.BotPrefix)+len(d.Trigger):])
	if !(len(request) > 0) {
		request = "6d6"
	}
	diceRoll := ""
	total := 0
	dMatch := diceRegexp.FindStringSubmatch(request)
	if len(dMatch) > 0 {
		numDice, numSides, addNum, roll := dMatch[1], dMatch[2], dMatch[3], 0
		if _numDice, err := strconv.Atoi(numDice); err == nil {
			if _numDice > 20 {
				return "But I have small hands, I can't hold that many dice :frowning:"
			}
			if _numSides, err := strconv.Atoi(numSides); err == nil {
				if _numSides > 32 {
					return "Wow those are strange die, I don't even know how to roll 'em :confused:"
				}
				for dice := 0; dice < _numDice; dice++ {
					if _numSides > 0 {
						roll = rand.Intn(_numSides) + 1
					} else {
						roll = 0
					}
					diceRoll = fmt.Sprintf("%s %d", diceRoll, roll)
					total += roll
				}
			} else {
				log.Printf("error: %v", err) // Non fatal error at strconv.Atoi() call
			}
		} else {
			log.Printf("error: %v", err) // Non fatal error at strconv.Atoi() call
		}

		if len(addNum) > 0 {
			if _addNum, err := strconv.Atoi(addNum); err == nil {
				total += _addNum
				diceRoll = fmt.Sprintf("%s %d", diceRoll, _addNum)
			} else {
				log.Printf("error: %v", err) // Non fatal error at strconv.Atoi() call
			}
		}
	}

	if len(diceRoll) > 0 {
		return fmt.Sprintf("```Dice roll: %s\tTotal: %d```", diceRoll, total)
	}
	return ""
}

func init() {
	rand.Seed(time.Now().Unix())

	_dice := &dice{
		Trigger: "dice",
	}
	addReaction(_dice.Trigger, _dice)
}
