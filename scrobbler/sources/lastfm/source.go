package lastfm

import (
	"errors"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
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
	newSong  func(song *scrobbler.Song, source scrobbler.ScrobbleSource)
}

func (s *Source) New(instance *scrobbler.Scrobbler, newSong func(song *scrobbler.Song, source scrobbler.ScrobbleSource), conf config.ModuleConfig) {
	s.newSong = newSong
	s.instance = instance
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
