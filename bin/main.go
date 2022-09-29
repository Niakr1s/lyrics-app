package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/niakr1s/lyricsgo/src/fs"
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
	InputFiles []string
	FfmpegCmd  string

	Simulate bool
}

// Checks if app is ready to go. If has some troubles, panics.
func makeSettings(args Args) (Settings, error) {
	res := Settings{
		InputFiles: []string{},
		FfmpegCmd:  "",
		Simulate:   args.Simulate,
	}
	var err error

	args.InputPath, err = fs.ToAbs(args.InputPath)
	if err != nil {
		return res, err
	}
	stat, err := os.Stat(args.InputPath)
	if err != nil {
		return res, fmt.Errorf("couldn't get stat of input path %s: %v", args.InputPath, err)
	}
	if stat.Mode().IsRegular() {
		res.InputFiles = []string{args.InputPath}
	} else if stat.Mode().IsDir() {
		res.InputFiles, err = fs.GetAllFilesFromDir(args.InputPath, args.Recoursive)
		if err != nil {
			return res, err
		}
	} else {
		return res, fmt.Errorf("input path %s is nor regular file, nor directory", args.InputPath)
	}
	res.InputFiles = fs.FilterMusicFiles(res.InputFiles)

	// checking ffmpeg
	ffmpegCmds := []string{"ffmpeg", "ffmpeg.exe"}
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

func main() {
	args := getArgs()
	settings, err := makeSettings(args)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("%+v | %+v\n", args, settings)
}
