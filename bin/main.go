package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/niakr1s/lyricsgo/lib"
)

// command line args
type Args struct {
	InputPath  string
	Recoursive bool
	Simulate   bool
}

func getArgs() Args {
	inputPath := flag.String("i", "", "Required! Path for input directory, or file. "+
		"If it's directory, app will proceed all files in directory (check recoursive arg).")
	recoursive := flag.Bool("r", false, "Recoursive directory")
	simulate := flag.Bool("s", false, "If turned on, app won't change any input file.")

	flag.Parse()

	return Args{
		InputPath:  *inputPath,
		Recoursive: *recoursive,
		Simulate:   *simulate,
	}
}

type Settings struct {
	// needed to check
	MusicFilePaths []string

	Simulate bool
}

func (s Settings) PrintInfo() {
	fmt.Printf("Got input: %d music files:\n", len(s.MusicFilePaths))
	for _, musicFilePath := range s.MusicFilePaths {
		fmt.Printf("\t%s\n", musicFilePath)
	}

	fmt.Printf("Simulate=%v, i will ", s.Simulate)
	if s.Simulate {
		fmt.Printf("simulate writing metadata to file.\n")
	} else {
		fmt.Printf("really write metadata to file.\n")
	}
}

// Checks if app is ready to go. If has some troubles, panics.
func makeSettings(args Args) (Settings, error) {
	res := Settings{
		MusicFilePaths: []string{},
		Simulate:       args.Simulate,
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
		res.MusicFilePaths, err = lib.GetAllFilesFromDir(args.InputPath, args.Recoursive)
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
	fmt.Printf("-----\n")
}

func doJob(settings Settings) {
	lyricQueries := make([]string, len(settings.MusicFilePaths))
	for i, musicFilePath := range settings.MusicFilePaths {
		lyricQueries[i] = filepath.Base(musicFilePath)
	}
	lyricResults := lib.GetLyricsAll(lyricQueries)
	log.Printf("GetLyricResult: %s\n", lyricResults.Info())

	if settings.Simulate {
		return
	}

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

func main() {
	args := getArgs()
	settings, err := makeSettings(args)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	settings.PrintInfo()

	printSeparator()
	fmt.Printf("Start\n")

	doJob(settings)

	fmt.Printf("End\n")
	printSeparator()
}
