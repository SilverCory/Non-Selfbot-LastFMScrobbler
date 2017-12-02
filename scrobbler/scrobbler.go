package scrobbler

import (
	"fmt"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler/sources"
	"github.com/SilverCory/go_discordrpc"
)

//go:generate go run ../../assets/generate.go ../../assets/unknown.png
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
			ret.sourcesByPriority[*sco] = v
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

	ret.UploadDefaults()

	return ret, nil
}

func (sc *Scrobbler) UploadDefaults() error {

	if err := sc.checkDefault("unknown_art", unknown_icon, 2); err != nil {
		return err
	}

	for k := range sc.sourcesByPriority {
		if err := sc.checkDefault(k.GetSourceName(), k.GetSourceIcon(), 1); err != nil {
			return err
		}
	}

	return nil
}

func (sc *Scrobbler) checkDefault(name, source string, Type int) error {
	asset, err := sc.AssetManager.GetAssetViaName(name)
	if asset == nil || err == sources.NotFoundError {
		fmt.Println("No asset or not found err")
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