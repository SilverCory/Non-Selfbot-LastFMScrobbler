package lastfm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var LASTFMKEY = ""
var LASTFMUSER = ""

type ScrobbleResponse struct {
	Tracks RecentTracks `json:"recenttracks"`
}

type RecentTracks struct {
	Tracks []Track `json:"track"`
}

func (rt *RecentTracks) FindNowPlaying() *Track {
	for _, track := range rt.Tracks {
		if strings.EqualFold("true", track.Attr.NowPlaying) {
			return &track
		}
	}
	return nil
}

type Artist struct {
	Text string `json:"#text"`
}

type Album struct {
	Text string `json:"#text"`
}

type Date struct {
	Date uint `json:"uts"`
}

type Image struct {
	URL  string `json:"#text"`
	Size string `json:"size"`
}

type Attr struct {
	NowPlaying string `json:"nowplaying"`
}

type TrackInfo struct {
	Track struct {
		Duration string `json:"duration"`
	} `json:"track"`
}

type Track struct {
	Artist   Artist        `json:"artist"`
	Name     string        `json:"name"`
	ID       string        `json:"mbid"`
	Images   []Image       `json:"image"`
	Attr     Attr          `json:"@attr"`
	Album    Album         `json:"album"`
	Duration time.Duration `json:"duration"`
}

func (t *Track) LoadDuration() {
	resp, err := http.Get(fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=%s&mbid=%s&format=json", LASTFMKEY, t.ID))
	if err != nil {
		fmt.Println("Error in duration load!", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 209 {
		fmt.Println("Error in duration load!", "Status code not 200, instead : "+strconv.Itoa(resp.StatusCode)+", status: "+resp.Status)
		return
	}

	jsonArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in duration load!", err)
		return
	}

	trackResponse := &TrackInfo{}
	if err := json.Unmarshal(jsonArr, trackResponse); err != nil {
		fmt.Println("Error in duration load!", err)
		return
	}

	if trackResponse != nil && trackResponse.Track.Duration != "" {
		dur, err := time.ParseDuration(trackResponse.Track.Duration + "ms")
		if err != nil {
			fmt.Println("Error loading duration..", trackResponse.Track.Duration, err)
		} else {
			t.Duration = dur
		}
	} else {
		t.Duration = time.Minute * 3
	}

}

func (t *Track) FindImageURL() string {

	for _, size := range []string{"large", "extralarge", "medium", "small"} {
		for _, v := range t.Images {
			if v.Size == size {
				return v.URL
			}
		}
	}

	return ""

}
