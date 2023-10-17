package main

import (
	"image/jpeg"
	"log"
	"os"
	"sourceafis/utils/encode/wsq"
)

func main() {
	f, err := os.Open("sample.wsq")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := wsq.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	ff, err := os.Create("img.jpeg")
	if err != nil {
		panic(err)
	}
	defer ff.Close()
	if err := jpeg.Encode(ff, img, nil); err != nil {
		log.Fatal(err)
	}
}
