package main

import (
	"context"
	"fmt"
	"sourceafis"
	"sourceafis/config"
	"time"

	"log"
)

type TransparencyContents struct {
}

func (c *TransparencyContents) Accepts(key string) bool {
	return true
}

func (c *TransparencyContents) Accept(key, mime string, data []byte) error {
	fmt.Printf("%d B  %s %s \n", len(data), mime, key)
	return nil
}

func main() {
	now := time.Now()
	config.LoadDefaultConfig()

	probeImg, err := sourceafis.LoadImage("probe.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	l := sourceafis.NewTransparencyLogger(new(TransparencyContents))
	tc := sourceafis.NewTemplateCreator(l)
	probe, err := tc.Template(probeImg)
	if err != nil {
		log.Fatal(err.Error())
	}
	candidateImg, err := sourceafis.LoadImage("matching.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	candidate, err := tc.Template(candidateImg)
	if err != nil {
		log.Fatal(err.Error())
	}

	candidateImg2, err := sourceafis.LoadImage("nonmatching.png")
	if err != nil {
		log.Fatal(err.Error())
	}
	candidate2, err := tc.Template(candidateImg2)
	if err != nil {
		log.Fatal(err.Error())
	}

	matcher, err := sourceafis.NewMatcher(l, probe)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()
	fmt.Println("matching score ===> ", matcher.Match(ctx, candidate))
	fmt.Println("non-matching score ===> ", matcher.Match(ctx, candidate2))
	fmt.Println("elapsed: ", time.Since(now))
}
