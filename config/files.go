package config

import (
	"encoding/json"
	"fmt"
	"github.com/shibukawa/configdir"
	"log"
	"os"
)

var confDir = configdir.New("discord-scrobbler", "daemon")

func New() *Config {
	conf := GetDefaultConfig()
	conf.Load()
	return conf
}

func (c *Config) Load() {
	confFile := confDir.QueryFolderContainsFile("settings.json")
	if confFile == nil {
		c.Save()
		// TODO do we want to do this?
		fmt.Println("The default configuration has been saved. Please edit this and restart!")
		os.Exit(0)
		return
	} else {
		data, err := confFile.ReadFile("settings.json")
		if err != nil {
			log.Fatal("There was an error loading the config!", err)
			return
		}

		err = json.Unmarshal(data, c)
		if err != nil {
			log.Fatal("There was an error loading the config!", err)
			return
		}
	}
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "\t")

	if err != nil {
		fmt.Println("There was an error saving the config!", err)
		return err
	}

	confFile := confDir.QueryFolders(configdir.Global)
	err = confFile[0].WriteFile("./settings.json", data)
	if err != nil {
		fmt.Println("There was an error saving the config!", err)
		return err
	}

	return nil

}
