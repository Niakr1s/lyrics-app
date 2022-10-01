package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/niakr1s/lyricsgo/lib"
)

// command line args
type Args struct {
	InputPath string
}

func getArgs() Args {
	if len(os.Args) != 2 {
		log.Fatalf("USAGE: app [input]\n\tWHERE [input] is path to music file or directory with music files")
	}

	return Args{
		InputPath: os.Args[1],
	}
}

type Settings struct {
	// needed to check
	MusicFilePaths []string
}

func (s Settings) PrintInfo() {
	log.Printf("Got input: %d music files:\n", len(s.MusicFilePaths))
	for _, musicFilePath := range s.MusicFilePaths {
		log.Printf("\t%s\n", musicFilePath)
	}
}

// Checks if app is ready to go. If has some troubles, panics.
func makeSettings(args Args) (Settings, error) {
	res := Settings{
		MusicFilePaths: []string{},
	}
	var err error

	args.InputPath = filepath.Clean(args.InputPath)
	args.InputPath, err = lib.ToAbs(args.InputPath)
	if err != nil {
		return res, err
	}
	stat, err := os.Stat(args.InputPath)
	if err != nil {
		return res, fmt.Errorf("couldn't get stat of input path %s: %v", args.InputPath, err)
	}
	if stat.Mode().IsRegular() {
		res.MusicFilePaths = []string{args.InputPath}
	} else if stat.Mode().IsDir() {
		res.MusicFilePaths, err = lib.GetAllFilesFromDir(args.InputPath, true)
		if err != nil {
			return res, err
		}
	} else {
		return res, fmt.Errorf("input path %s is nor regular file, nor directory", args.InputPath)
	}
	res.MusicFilePaths = lib.FilterMusicFiles(res.MusicFilePaths)

	return res, nil
}

func printSeparator() {
	log.Printf("-----\n")
}

func doJob(settings Settings) {
	lyricQueries := make([]string, len(settings.MusicFilePaths))
	for i, musicFilePath := range settings.MusicFilePaths {
		lyricQueries[i] = filepath.Base(musicFilePath)
	}
	lyricResults := lib.GetLyricsAll(lyricQueries)
	log.Printf("GetLyricResult: %s\n", lyricResults.Info())

	lyricsJobs := []lib.WriteLyricsJob{}
	for i, lyricJobResult := range lyricResults {
		if lyricJobResult.Err != nil {
			continue
		}
		lyricsJob := lib.WriteLyricsJob{
			MusicFilePath: settings.MusicFilePaths[i],
			Lyrics:        lyricJobResult.Lyrics,
		}
		lyricsJobs = append(lyricsJobs, lyricsJob)
	}

	lyricsResults := lib.WriteLyricsAll(lyricsJobs, false)
	log.Printf("Write lyrics result: %s\n", lyricsResults.Info())
}

func setLogger(filename string) io.Closer {
	fi, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("error getting stat for file %v: %v", filename, err)
	}

	var logFilePath string
	if fi.Mode().IsRegular() {
		ext := filepath.Ext(filename)
		logFilePath = filepath.Join(strings.TrimSuffix(filename, ext) + ".log")
	} else if fi.Mode().IsDir() {
		logFilePath = filepath.Join(filename, "lyrics.log")
	}
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

func main() {
	args := getArgs()
	defer setLogger(args.InputPath).Close()

	settings, err := makeSettings(args)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	settings.PrintInfo()

	printSeparator()
	log.Printf("Start\n")

	doJob(settings)

	log.Printf("End\n")
	printSeparator()
}
