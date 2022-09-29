package lib

import (
	"fmt"
	"sync"

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

type GetLyricsAllResults []GetLyricsAllResult

func (r GetLyricsAllResults) Info() string {
	lyricCount := 0
	errorCount := 0
	for _, v := range r {
		if v.Lyrics != "" {
			lyricCount++
		}
		if v.Err != nil {
			errorCount++
		}
	}

	return fmt.Sprintf("From %d queries got %d lyrics and %d errors", len(r), lyricCount, errorCount)
}

type GetLyricsAllResult struct {
	Query  string
	Lyrics string
	Err    error
}

// Garantees to reserve input order.
func GetLyricsAll(queries []string) GetLyricsAllResults {
	results := make(GetLyricsAllResults, len(queries))

	wg := sync.WaitGroup{}
	for i, query := range queries {
		i := i
		query := query
		res := GetLyricsAllResult{
			Query: query,
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				results[i] = res
			}()

			lyrics, err := GetLyrics(query)
			if err != nil {
				res.Err = fmt.Errorf("lyrics not found, reason: %v: %s", err, query)
				return
			}
			res.Lyrics = lyrics
		}()
	}
	wg.Wait()

	return results
}
