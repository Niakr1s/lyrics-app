package lib

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

//go:embed write_lyrics.py
var write_lyrics_py string
var write_lyrics_py_filepath string

func init() {
	var err error

	// preinstall mutagen
	err = exec.Command("pip", "install", "mutagen").Run()
	if err != nil {
		log.Fatalf("pip not found, please, install python version 3.10+: %v", err)
	}

	write_lyrics_py_filepath = filepath.Join(os.TempDir(), "set_lyrics.py")
	err = os.WriteFile(write_lyrics_py_filepath, []byte(write_lyrics_py), 0666)
	if err != nil {
		log.Fatalf("couldn't write %s", write_lyrics_py_filepath)
	}
}

type WriteLyricsJob struct {
	MusicFilePath string
	Lyrics        string
}

type WriteLyricsResults []WriteLyricsAllResult

func (r WriteLyricsResults) Info() string {
	successCount := 0
	for _, v := range r {
		if v.Err == nil {
			successCount++
		}
	}

	return fmt.Sprintf("From %d music files got %d success", len(r), successCount)
}

type WriteLyricsAllResult struct {
	MusicFilePath string
	Err           error
}

// Garantees to reserve input order.
func WriteLyricsAll(jobs []WriteLyricsJob, withOutput bool) WriteLyricsResults {
	results := make(WriteLyricsResults, len(jobs))

	wg := sync.WaitGroup{}
	for i, job := range jobs {
		i := i
		job := job
		jobResult := WriteLyricsAllResult{
			MusicFilePath: job.MusicFilePath,
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				results[i] = jobResult
			}()

			err := WriteLyrics(job.MusicFilePath, job.Lyrics, withOutput)

			if err != nil {
				jobResult.Err = fmt.Errorf("couldn't write lyrics for file %v, reason: %v", job.MusicFilePath, err)
				return
			}
		}()
	}
	wg.Wait()

	return results
}

func WriteLyrics(filepath string, lyrics string, withOutput bool) error {
	cmd := exec.Command("python", write_lyrics_py_filepath, filepath, lyrics)

	if withOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error while writing lyrics to %v, lyrics len %d, reason: %v", filepath, len(lyrics), err)
	}
	return nil
}
