package lastfm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const name = "Last FM"

//go:generate go run ../../../assets/generate.go ../../../assets/lastfm_black.png
const icon = ``

type Config struct {
	config.ModuleConfig
	Username string `json:"username"`
	APIKey   string `json:"api_key"`
}

func init() {
	source := &Source{}
	scrobbler.RegisterSource(name, source)
	config.AddDefaultConfig(name, source.GetDefaultConfig())
}

type Source struct {
	scrobbler.ScrobbleSource
	instance *scrobbler.Scrobbler
	config   Config
	newSong  func(song *scrobbler.Song, source scrobbler.ScrobbleSource)

	lastString   string
	lastDuration time.Duration
	timeThen     time.Time
}

func (s *Source) New(instance *scrobbler.Scrobbler, newSong func(song *scrobbler.Song, source scrobbler.ScrobbleSource), conf config.ModuleConfig) {
	s.newSong = newSong
	s.instance = instance
	s.config = conf.(Config)
}

func (s *Source) GetSourceIcon() string {
	return icon
}

func (s *Source) GetSourceName() string {
	return name
}

func (s *Source) GetDefaultConfig() config.ModuleConfig {
	return Config{
		APIKey:   "API_KEY_HERE",
		Username: "USERNAMEPLEASE",
	}
}

func (s *Source) Start() error {
	return errors.New("NOT IMPLEMENTED") // TODO
}

func (s *Source) Stop() error {
	return errors.New("NOT IMPLEMENTED") // TODO
}

func (s *Source) QueryNewSong() error {

	// TODO

	resp, err := http.Get("https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=" + s.config.Username + "&api_key=" + s.config.APIKey + "&format=json&limit=1")
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 209 {
		return errors.New("Status code not 200, instead : " + strconv.Itoa(resp.StatusCode) + ", status: " + resp.Status)
	}

	defer resp.Body.Close()
	jsonArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	scrobbleResponse := &ScrobbleResponse{}
	if err := json.Unmarshal(jsonArr, scrobbleResponse); err != nil {
		return err
	}

	recentTracks := scrobbleResponse.Tracks
	track := recentTracks.FindNowPlaying()
	if track == nil && time.Now().After(s.timeThen.Add(s.lastDuration+(5*time.Second))) {
		s.newSong(nil, s)
		return nil
	}

	track.LoadDuration()
	compareString := track.Name + "^^^" + track.Artist.Text
	if s.lastString != compareString {
		if s.currentId != bot.Bot.Conf.DiscordDefaultImageID {
			s.lastId = s.currentId
		}
		s.currentId = bot.Bot.Conf.DiscordDefaultImageID
		if url := track.FindImageURL(); url != "" {
			s.UploadCover(url)
		} else {
			fmt.Println("URL was empty!")
		}
		s.Assets.RemoveAsset(s.lastId)
		s.TimeThen = time.Now()
	}

	s.lastString = compareString

}
