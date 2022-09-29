package lib

import (
	"fmt"

	"github.com/raitonoberu/ytmusic"
)

func GetLyrics(query string) (string, error) {

	ytSearchResult, err := ytmusic.Search(query).Next()
	if err != nil {
		return "", err
	}

	tracks := ytSearchResult.Tracks
	if len(tracks) == 0 {
		return "", fmt.Errorf("tracks are empty for query %s", query)
	}

	track := tracks[0]
	lyrics, err := ytmusic.GetLyrics(track.VideoID)
	if err != nil {
		return "", err
	}
	if len(lyrics) == 0 {
		return "", fmt.Errorf("lyrics are empty")
	}
	return lyrics, nil
}
