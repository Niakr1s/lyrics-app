package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/niakr1s/lyricsgo/lib"
)

// command line args
type Args struct {
	InputPath  string
	Recoursive bool
	Simulate   bool
	FfmpegPath string
}

func getArgs() Args {
	inputPath := flag.String("i", "", "Required! Path for input directory, or file. "+
		"If it's directory, app will proceed all files in directory (check recoursive arg).")
	recoursive := flag.Bool("r", false, "Recoursive directory")
	simulate := flag.Bool("s", false, "If turned on, app won't change any input file.")
	ffmpegPath := flag.String("ffmpeg", "", "Path to ffmpeg executable. If not given, tries to find ffmpeg executable in PATH.")

	flag.Parse()

	return Args{
		InputPath:  *inputPath,
		Recoursive: *recoursive,
		Simulate:   *simulate,
		FfmpegPath: *ffmpegPath,
	}
}

type Settings struct {
	// needed to check
	MusicFilePaths []string
	FfmpegCmd      string

	Simulate bool
}

func (s Settings) PrintInfo() {
	fmt.Printf("Got input: %d music files:\n", len(s.MusicFilePaths))
	for _, musicFilePath := range s.MusicFilePaths {
		fmt.Printf("\t%s\n", musicFilePath)
	}
	fmt.Printf("Found ffmpeg in %s\n", s.FfmpegCmd)

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
		FfmpegCmd:      "",
		Simulate:       args.Simulate,
	}
	var err error

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

	// checking ffmpeg
	ffmpegCmds := []string{args.FfmpegPath, "ffmpeg", "ffmpeg.exe"}
	for _, ffmpegCmd := range ffmpegCmds {
		err := exec.Command(ffmpegCmd, "-h").Run()
		if err == nil {
			res.FfmpegCmd = ffmpegCmd
			break
		}
	}
	if res.FfmpegCmd == "" {
		return res, fmt.Errorf("no ffmpeg executable in path")
	}

	if err != nil {
		log.Fatalf("no ffmpeg file in PATH: %v", err)
	}

	return res, nil
}

func printSeparator() {
	fmt.Printf("-----\n")
}

type MetadataWriter interface {
	// should return output file path or error
	WriteMetadata(musicFilePath string, meta map[string]string) (string, error)
}

type SuccessJob struct {
	FileName string
	Lyrics   string
}

type FailureJob struct {
	FileName string
	Error    error
}

type JobResult struct {
	InputFiles   []string
	SuccessFiles map[string]SuccessJob
	FailureFiles map[string]FailureJob
}

func (r JobResult) String() string {
	return fmt.Sprintf("Got %d input files, success: %d, failure: %d", len(r.InputFiles), len(r.SuccessFiles), len(r.FailureFiles))
}

func doJob(settings Settings) (JobResult, error) {
	jobResult := JobResult{
		InputFiles:   settings.MusicFilePaths,
		SuccessFiles: map[string]SuccessJob{},
		FailureFiles: map[string]FailureJob{},
	}

	var metadataWriter MetadataWriter = lib.NewFfmpegMetadataWriter(settings.FfmpegCmd)

	wg := sync.WaitGroup{}
	for _, musicFilePath := range settings.MusicFilePaths {
		musicFilePath := musicFilePath
		wg.Add(1)

		go func() (err error) {
			defer wg.Done()
			defer func() {
				if _, ok := jobResult.SuccessFiles[musicFilePath]; !ok {
					jobResult.FailureFiles[musicFilePath] = FailureJob{
						FileName: musicFilePath,
						Error:    err,
					}
				}
			}()

			fmt.Printf("Start search lyrics for %s\n", musicFilePath)
			query := path.Base(musicFilePath)
			lyrics, err := lib.GetLyrics(query)
			if err != nil {
				fmt.Printf("Lyrics not found, reason: %v: %s\n", err, musicFilePath)
				return
			}

			fmt.Printf("Lyrics found, len=%d: %s\n", len(lyrics), musicFilePath)
			if settings.Simulate {
				jobResult.SuccessFiles[musicFilePath] = SuccessJob{
					FileName: musicFilePath,
					Lyrics:   lyrics,
				}
				return
			}

			outputFilePath, err := metadataWriter.WriteMetadata(musicFilePath, map[string]string{
				"Lyrics": lyrics,
			})
			defer os.Remove(outputFilePath) // just clean at the end

			if err != nil {
				fmt.Printf("Couldn't write metadata, reason: %v: %v\n", err, musicFilePath)
				return
			}

			err = os.Remove(musicFilePath)
			if err != nil {
				fmt.Printf("Couldn't remove file, reason: %v: %v\n", err, musicFilePath)
				return
			}

			err = os.Rename(outputFilePath, musicFilePath)
			if err != nil {
				fmt.Printf("Couldn't rename file, reason: %v: %v\n", err, outputFilePath)
				return
			}

			fmt.Printf("Metadata write success: %v\n", musicFilePath)
			jobResult.SuccessFiles[musicFilePath] = SuccessJob{
				FileName: musicFilePath,
				Lyrics:   lyrics,
			}
			return
		}()
	}
	wg.Wait()

	return jobResult, nil
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

	if result, err := doJob(settings); err != nil {
		log.Fatalf("%v\n", err)
	} else {
		fmt.Printf("%s\n", result)
	}

	fmt.Printf("End\n")
	printSeparator()
}
