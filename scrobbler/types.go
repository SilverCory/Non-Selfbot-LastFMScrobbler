package scrobbler

import (
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"io"
)

type ScrobbleSource interface {
	io.Closer
	New(instance *Scrobbler, newSong func(song *Song, source *ScrobbleSource), moduleConfig config.ModuleConfig)
	GetSourceIcon() string
	GetSourceName() string
	Start() error
	Stop() error
	GetDefaultConfig() *config.ModuleConfig
}

type ImageID string

type Song struct {
	Title   string
	Artist  string
	Album   string
	Artwork ImageID
}
