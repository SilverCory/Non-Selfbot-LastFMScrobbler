package main

import (
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
	_ "github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler/sources/lastfm"
	_ "github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler/sources/spotify"
	"time"
)

func main() {
	conf := config.New()

	sc, err := scrobbler.New(conf)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	println(sc)

}
