// The settings package should be importable by everything under the project
package settings

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type settings struct {
	Discord struct {
		BotPrefix string
	}

	IRE struct {
		DeathsightEnabled bool
		LastID            int
	}
}

var settingsPath string
var Settings settings

func Init(path string) error {
	settingsPath = path
	if err := Settings.Load(); err != nil {
		return err // Error at Settings.Load() call
	}
	return nil
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
	if y_settings, err := yaml.Marshal(Settings); err == nil {
		if err := ioutil.WriteFile(settingsPath, y_settings, 0600); err != nil {
			return err // Error at WriteFile() call
		}
	} else {
		return err // Error at yaml.Marshal() call
	}
	return nil
}
