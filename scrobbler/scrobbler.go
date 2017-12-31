package scrobbler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/SilverCory/EzVote/pkg/dep/sources/https---github.com-kataras-iris/core/errors"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler/sources"
	"github.com/SilverCory/go_discordrpc"
	"image"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"time"
)

//go:generate go run ../assets/generate.go ../assets/unknown.png
const unknown_icon = ``

type Scrobbler struct {
	nowPlaying        *Song
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
		for module := range ret.sourcesByPriority {
			if module.GetSourceName() == configName {
				module.UpdateSource(ret, ret.newSong, moduleConfig)
			}
		}
	}

	return ret, nil
}

func (sc *Scrobbler) Start() error {
	return nil
}

func (sc *Scrobbler) newSong(song *Song, source ScrobbleSource) {
	sc.nowPlaying = song
	// TODO newSong send to discord.
}

func (sc *Scrobbler) GetNowPlaying() *Song {
	return sc.nowPlaying
}

func (sc *Scrobbler) UploadCoverViaURL(url string) (*sources.DiscordAsset, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 209 {
		return nil, errors.New("Status code not 200, instead : " + strconv.Itoa(resp.StatusCode) + ", status: " + resp.Status)
	}

	defer resp.Body.Close()
	imageBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_, format, err := image.Decode(bytes.NewBuffer(imageBody))
	if err != nil {
		return nil, err
	}

	data := mime.TypeByExtension("." + format)
	img64 := base64.StdEncoding.EncodeToString(imageBody)
	return sc.UploadCoverImage("data:" + data + ";base64," + img64)
}

func (sc *Scrobbler) UploadCoverImage(image64 string) (*sources.DiscordAsset, error) {
	return sc.AssetManager.AddAsset("cover_"+strconv.Itoa(int(time.Now().Unix())), image64, 2)
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
