package lastfm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const name = "Last FM"

//go:generate go run ../../../assets/generate.go ../../../assets/lastfm_black.png
const icon = ``

type Config struct {
	Username string `json:"username"`
	Api_key  string `json:"api_key"`
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
	lastId       string
	currentId    string
	lastDuration time.Duration
	timeThen     time.Time
}

func (s *Source) UpdateSource(instance *scrobbler.Scrobbler, newSong func(song *scrobbler.Song, source scrobbler.ScrobbleSource), conf config.ModuleConfig) {
	s.newSong = newSong
	s.instance = instance
	s.config = Config{}
	fmt.Println(mapstructure.Decode(conf, &s.config))
	fmt.Println(s.config)
	fmt.Println(conf)
	fmt.Println(s.config.Api_key)

	LASTFMKEY = s.config.Api_key
	LASTFMUSER = s.config.Username

}

func (s *Source) GetSourceIcon() string {
	return icon
}

func (s *Source) GetSourceName() string {
	return name
}

func (s *Source) GetDefaultConfig() config.ModuleConfig {
	return Config{
		Api_key:  "API_KEY_HERE",
		Username: "USERNAMEPLEASE",
	}
}

var started = false

func (s *Source) Start() error {

	if started {
		return errors.New("already started")
	}

	started = true
	fmt.Println("STARTEWD LASTFM")
	for {

		time.Sleep(time.Second * 10)
		s.QueryNewSong()

	}
}

func (s *Source) Stop() error {
	return errors.New("STOP NOTT IMPLEMENTED") // TODO
}

func (s *Source) QueryNewSong() error {

	// TODO

	song := &scrobbler.Song{}

	fmt.Println("QUERY NEW SONG")
	resp, err := http.Get("https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=" + s.config.Username + "&api_key=" + s.config.Api_key + "&format=json&limit=1")
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
		if url := track.FindImageURL(); url != "" {
			asset, err := s.instance.UploadCoverViaURL(url)
			if err == nil && asset != nil {
				song.Artwork = scrobbler.ImageID(asset.Name)
			} else {
				song.Artwork = "unknown_art"
			}
		} else {
			song.Artwork = "unknown_art"
		}
		s.timeThen = time.Now()

		song.Artist = track.Artist.Text
		song.Album = track.Album.Text
		song.Title = track.Name
		song.End = time.Now().Add(track.Duration - (13 * time.Second))
		s.newSong(song, s)
	}

	s.lastString = compareString
	return nil

}
