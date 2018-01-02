package mpris2

import (
	"fmt"
	"time"
)

import (
	"errors"
	"fmt"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
	"github.com/lann/mpris2-go"

	"github.com/guelfey/go.dbus"
	"time"
)

const name = "MPRIS2"

//go:generate go run ../../../assets/generate.go ../../../assets/lastfm_black.png
const icon = ``

func init() {
	source := &Source{}
	scrobbler.RegisterSource(name, source)
	config.AddDefaultConfig(name, source.GetDefaultConfig())
}

type Source struct {
	scrobbler.ScrobbleSource
	instance *scrobbler.Scrobbler
	newSong  func(song *scrobbler.Song, source scrobbler.ScrobbleSource)
	dbus     *dbus.Conn
}

func (s *Source) UpdateSource(instance *scrobbler.Scrobbler, newSong func(song *scrobbler.Song, source scrobbler.ScrobbleSource), conf config.ModuleConfig) {
	s.newSong = newSong
	s.instance = instance
	conn, err := mpris2.Connect()
}

func (s *Source) GetSourceIcon() string {
	return icon
}

func (s *Source) GetSourceName() string {
	return name
}

func (s *Source) GetDefaultConfig() config.ModuleConfig {
	return nil
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
