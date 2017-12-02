package main

import (
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/config"
	"github.com/SilverCory/Non-Selfbot-LastFMScrobbler/scrobbler"
)

func main() {
	conf := config.New()

	sc, err := scrobbler.New(conf)
	if err != nil {
		panic(err)
	}

	println(sc)

}
