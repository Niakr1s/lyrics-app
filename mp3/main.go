package main

import (
	"log"

	"github.com/mikkyang/id3-go"
	v2 "github.com/mikkyang/id3-go/v2"
)

func main() {
	f, err := id3.Open("/mnt/d/Ayano Mashiro - Gentou.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ft := v2.V23FrameTypeMap["USLT"]
	text := `-----------
SUPER DUPER LYRICS`
	textFrame := v2.NewTextFrame(ft, text)
	f.AddFrames(textFrame)

}
