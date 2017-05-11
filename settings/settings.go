// Package settings provides a way to save and load runtime settings
package settings

import (
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type settings struct {
	sync.RWMutex

	Discord struct {
		BotPrefix    string
		BotAdmin     string
		DefaultRoles map[string]string
	}

	IRE struct {
		DeathsightEnabled bool
		LastID            int
	}
}

var settingsPath string

// Settings - runtime settings, ensure Lock() -> defer Unlock() on write
var Settings settings

// Init is called by main() once config.yaml is read/processed
func Init(path string) error {
	Settings.Discord.DefaultRoles = make(map[string]string) // ??

	settingsPath = path
	return Settings.Load()
}

func (s *settings) Load() error {
	if _, err := os.Stat(settingsPath); err == nil { // settings file already exists
		if data, err := ioutil.ReadFile(settingsPath); err == nil {
			if err := yaml.Unmarshal([]byte(data), &Settings); err != nil {
				return err // Error at yaml.Unmarshal() call
			}
		} else {
			return err // Error at ioutil.ReadFile() call
		}
	} else { // settings file does not exist, let's create a new one
		/* ---------- BEGIN DEFAULT SETTINGS ----------*/
		Settings.Discord.BotPrefix = "!"
		Settings.Discord.BotAdmin = "0"
		Settings.IRE.DeathsightEnabled = true
		Settings.IRE.LastID = 0
		/* ----------- END DEFAULT SETTINGS -----------*/
		if err := Settings.Save(); err != nil {
			return err // Error at Settings.Save() call
		}
	}
	return nil
}

func (s *settings) Save() error {
	if ySettings, err := yaml.Marshal(Settings); err == nil {
		if err := ioutil.WriteFile(settingsPath, ySettings, 0600); err != nil {
			return err // Error at WriteFile() call
		}
	} else {
		return err // Error at yaml.Marshal() call
	}
	return nil
}
