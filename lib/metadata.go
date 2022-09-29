package lib

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type MetadataWriter interface {
	// should return output file path or error
	WriteMetadata(WriteMetadataJob) (string, error)
}

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
func (w *FfmpegMetadataWriter) WriteMetadata(job WriteMetadataJob) (string, error) {
	musicFilePath, err := ToAbs(job.MusicFilePath)
	if err != nil {
		return "", err
	}
	outputFilePath := makeOutputFilePath(musicFilePath)

	args := []string{"-i", musicFilePath}
	for k, v := range job.Meta {
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

type WriteMetadataJob struct {
	MusicFilePath string
	Meta          map[string]string
}

type WriteMetadataAllResults []WriteMetadataAllResult

func (r WriteMetadataAllResults) Info() string {
	errorCount := 0
	for _, v := range r {
		if v.Err != nil {
			errorCount++
		}
	}

	return fmt.Sprintf("From %d music files got %d success and %d errors", len(r), len(r)-errorCount, errorCount)
}

type WriteMetadataAllResult struct {
	MusicFilePath string
	Err           error
}

// Garantees to reserve input order.
func WriteMetadataAll(writer MetadataWriter, jobs []WriteMetadataJob) WriteMetadataAllResults {
	results := make(WriteMetadataAllResults, len(jobs))

	wg := sync.WaitGroup{}
	for i, job := range jobs {
		i := i
		job := job
		jobResult := WriteMetadataAllResult{
			MusicFilePath: job.MusicFilePath,
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				results[i] = jobResult
			}()

			outputFilePath, err := writer.WriteMetadata(job)
			defer os.Remove(outputFilePath) // just clean at the end

			if err != nil {
				jobResult.Err = fmt.Errorf("couldn't write metadata, reason: %v: %v => %v", err, job.MusicFilePath, outputFilePath)
				return
			}

			err = os.Rename(outputFilePath, job.MusicFilePath)
			if err != nil {
				jobResult.Err = fmt.Errorf("couldn't move file, reason: %v: %v => %v", err, outputFilePath, job.MusicFilePath)
				return
			}
		}()
	}
	wg.Wait()

	return results
}
