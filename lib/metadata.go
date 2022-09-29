package lib

import (
	"fmt"
	"os/exec"
)

func metadataSuffix() string {
	return ".metadata"
}

func makeOutputFilePath(inputFilePath string) string {
	return inputFilePath + metadataSuffix()
}

type FfmpegMetadataWriter struct {
	ffmpegPath string
}

func NewFfmpegMetadataWriter(ffmpegPath string) *FfmpegMetadataWriter {
	return &FfmpegMetadataWriter{ffmpegPath: ffmpegPath}
}

// returns output filePath and error
func (w *FfmpegMetadataWriter) WriteMetadata(musicFilePath string, meta map[string]string) (string, error) {
	musicFilePath, err := ToAbs(musicFilePath)
	if err != nil {
		return "", err
	}
	outputFilePath := makeOutputFilePath(musicFilePath)

	args := []string{"-i", musicFilePath}
	for k, v := range meta {
		args = append(args, "-metadata:s:0", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, "-c:a", "copy", outputFilePath, "-y")
	cmd := exec.Command(w.ffmpegPath, args...)
	err = cmd.Run()
	if err != nil {
		return outputFilePath, err
	}

	return outputFilePath, nil
}
