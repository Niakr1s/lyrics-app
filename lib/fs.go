package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ToAbs(inputPath string) (string, error) {
	if !filepath.IsAbs(inputPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return inputPath, fmt.Errorf("couldn't get current working directory: %v", err)
		}
		inputPath = filepath.Join(cwd, inputPath)
	}
	inputPath = filepath.Clean(inputPath)
	return inputPath, nil
}

func GetAllFilesFromDir(dirPath string, recoursive bool) ([]string, error) {
	res := []string{}

	dirPath, err := ToAbs(dirPath)
	if err != nil {
		return res, err
	}

	dirItems, err := os.ReadDir(dirPath)
	if err != nil {
		return res, fmt.Errorf("couldn't get items of dir %s: %v", dirPath, err)
	}
	for _, dirItem := range dirItems {
		if dirItem.Type().IsRegular() {
			res = append(res, filepath.Clean(filepath.Join(dirPath, dirItem.Name())))
		} else if dirItem.IsDir() && recoursive {
			subDirItems, err := GetAllFilesFromDir(filepath.Join(dirPath, dirItem.Name()), true)
			if err != nil {
				return nil, err
			}
			res = append(res, subDirItems...)
		}
	}
	return res, nil
}

func IsMusicFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".mp3") || strings.HasSuffix(filePath, ".ogg")
}

func FilterMusicFiles(filePaths []string) []string {
	return Filter(filePaths, IsMusicFile)
}
