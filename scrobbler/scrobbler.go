package scrobbler

import (
	"fmt"
	"github.com/SilverCory/LastFMScrobbler/scrobbler"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler/sources"
	"github.com/SilverCory/go_discordrpc"
)

//go:generate go run ../assets/generate.go ../assets/unknown.png
const unknown_icon = ``

type Scrobbler struct {
	sourcesByPriority map[ScrobbleSource]int
	AssetManager      *sources.AssetManager
	DiscordRPC        *go_discordrpc.API
}

func New(config *config.Config) (*Scrobbler, error) {

	ret := &Scrobbler{
		sourcesByPriority: make(map[ScrobbleSource]int),
	}

	for k, v := range config.ModulePriorities {
		sco, ok := GetScrobbler(k)
		if ok {
			ret.sourcesByPriority[sco] = v
		}
	}

	api, err := go_discordrpc.New(config.ApplicationID)
	if err != nil {
		return nil, err
	}

	ret.DiscordRPC = &*api
	fmt.Println("Error: ", ret.DiscordRPC.Open())
	fmt.Println("Open: ", ret.DiscordRPC.IsOpen())

	assetManager, err := sources.NewAssetManager(config.AuthorisationToken, config.ApplicationID)
	if err != nil {
		return nil, err
	}
	ret.AssetManager = assetManager

	ret.uploadDefaults()

	for configName, moduleConfig := range config.ModuleConfigs {
		for module, _ := range ret.sourcesByPriority {
			if module.GetSourceName() == configName {
				module.New(ret, ret.newSong, moduleConfig)
			}
		}
	}

	return ret, nil
}

func (sc *Scrobbler) Start() error {

}

func (sc *Scrobbler) newSong(song *Song, source ScrobbleSource) {
	// TODO newSong
}

func (sc *Scrobbler) uploadDefaults() error {

	if err := sc.checkDefault("unknown_art", unknown_icon, 2); err != nil {
		return err
	}

	for k := range sc.sourcesByPriority {
		if err := sc.checkDefault(k.GetSourceName(), k.GetSourceIcon(), 1); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (sc *Scrobbler) checkDefault(name, source string, Type int) error {
	asset, err := sc.AssetManager.GetAssetViaName(name)
	if asset == nil || err == sources.NotFoundError {
		_, err := sc.AssetManager.AddAsset(name, source, Type)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		return err
	}
	return nil
}
