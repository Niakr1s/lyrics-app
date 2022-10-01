package lib

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestToAbs(t *testing.T) {
	type args struct {
		inputPath string
	}
	tests := []struct {
		name       string
		args       args
		wantSuffix string
		wantErr    bool
	}{
		{"empty", args{""}, filepath.Join("lyrics-app", "lib"), false},
		{"current", args{"."}, filepath.Join("lyrics-app", "lib"), false},
		{"test", args{"test"}, filepath.Join("lyrics-app", "lib", "test"), false},
		{"up", args{"../test"}, filepath.Join("lyrics-app", "test"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToAbs(tt.args.inputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToAbs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !filepath.IsAbs(got) {
				t.Errorf("%v is not Abs", got)
			}
			if !strings.HasSuffix(got, tt.wantSuffix) {
				t.Errorf("%v doesn't have suffiex %v", got, tt.wantSuffix)
			}
		})
	}
}

func TestGetAllFilesFromDir(t *testing.T) {
	type args struct {
		dirPath    string
		recoursive bool
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
		wantErr  bool
	}{
		{"non-recoursive test", args{"../test", false}, 3, false},
		{"non-recoursive subdir", args{"../test/subdir", false}, 6, false},
		{"non-recoursive error", args{"../test/error", false}, 0, true},

		{"recoursive test", args{"../test", true}, 9, false},
		{"recoursive subdir", args{"../test/subdir", true}, 6, false},
		{"recoursive error", args{"../test/error", true}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllFilesFromDir(tt.args.dirPath, tt.args.recoursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFilesFromDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantSize {
				t.Errorf("len(got)=%v != tt.wantSize=%v", got, tt.wantSize)
			}
		})
	}
}

func TestFilterMusicFiles(t *testing.T) {
	type args struct {
		dirPath string
	}
	tests := []struct {
		name       string
		args       args
		wantLength int
	}{
		{"test", args{"../test"}, 6},
		{"subdir", args{"../test/subdir"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePaths, err := GetAllFilesFromDir(tt.args.dirPath, true)
			if err != nil {
				t.Fatal(err)
			}
			if got := FilterMusicFiles(filePaths); len(got) != tt.wantLength {
				t.Errorf("len(got)=%v != tt.wantLength=%v", len(got), tt.wantLength)
			}
		})
	}
}
