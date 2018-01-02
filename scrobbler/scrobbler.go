package scrobbler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler/sources"
	"github.com/SilverCory/go_discordrpc"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:generate go run ../assets/generate.go ../assets/unknown.png
const unknown_icon = ``

type Scrobbler struct {
	nowPlaying        *Song
	sourcesByPriority map[ScrobbleSource]int
	AssetManager      *sources.AssetManager
	DiscordRPC        *go_discordrpc.API
	Config            *config.Config
}

func New(config *config.Config) (*Scrobbler, error) {

	ret := &Scrobbler{
		sourcesByPriority: make(map[ScrobbleSource]int),
		Config:            config,
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
	assets, err := ret.AssetManager.GetAllAssets()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, v := range assets {
			if strings.HasPrefix(v.Name, "cover_") {
				go ret.AssetManager.RemoveAsset(v.ID)
			}
		}
	}

	for configName, moduleConfig := range config.ModuleConfigs {
		for module := range ret.sourcesByPriority {
			fmt.Println(module)
			if module.GetSourceName() == configName {
				fmt.Println(module.GetSourceIcon())
				module.UpdateSource(
					ret,
					ret.newSong,
					moduleConfig)
			}
		}
	}

	ret.uploadDefaults()

	return ret, nil
}

func (sc *Scrobbler) Start() error {

	fmt.Println("STARTING SCROBBLER")
	for k := range sc.sourcesByPriority {
		go func() {
			if err := k.Start(); err != nil {
				fmt.Println(err)
			}
		}()
	}

	return nil

}

func (sc *Scrobbler) newSong(song *Song, source ScrobbleSource) {

	if song == nil {
		fmt.Println("SONG IS NIL!")
		return
	}

	assets := &go_discordrpc.Assets{}

	assets.LargeImageID = fmt.Sprint(song.Artwork)
	assets.LargeText = song.Album

	if asset, err := sc.AssetManager.GetAssetViaName(source.GetSourceName()); err == nil {
		assets.SmallImageID = asset.Name
		assets.SmallText = source.GetSourceName()
	}

	timestamps := &go_discordrpc.TimeStamps{}
	if song.End.After(time.Now()) {
		timestamps.EndTimestamp = song.End.Unix()
	}

	sc.DiscordRPC.SetRichPresence(&go_discordrpc.Activity{
		TimeStamps: timestamps,
		State:      song.Title,
		Details:    song.Artist,
		Assets:     assets,
	})

	if sc.nowPlaying != nil {
		go sc.AssetManager.RemoveAssetViaName(string(sc.nowPlaying.Artwork))
	}
	sc.nowPlaying = song
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
