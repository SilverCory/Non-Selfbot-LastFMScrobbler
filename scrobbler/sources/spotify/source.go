package spotify

import (
	"errors"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
)

const name = "Spotify"

//go:generate go run ../../../assets/generate.go ../../../assets/spotify.png
const icon = ``

func init() {
	source := &Source{}
	scrobbler.RegisterSource(name, &source.ScrobbleSource)
}

type Source struct {
	scrobbler.ScrobbleSource
	instance *scrobbler.Scrobbler
	newSong  func(song *scrobbler.Song, source *scrobbler.ScrobbleSource)
}

func (s *Source) New(instance *scrobbler.Scrobbler, newSong func(song *scrobbler.Song, source *scrobbler.ScrobbleSource)) {
	s.newSong = newSong
	s.instance = instance
}

func (s *Source) GetSourceIcon() string {
	return icon
}

func (s *Source) GetSourceName() string {
	return name
}

func (s *Source) Start() error {
	return errors.New("NOT IMPLEMENTED") // TODO
}

func (s *Source) Stop() error {
	return errors.New("NOT IMPLEMENTED") // TODO
}