package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/peterh/liner"
)

const Version = "v0.1.0"

var (
	origin       = flag.String("origin", "", "websocket origin")
	printVersion = flag.Bool("version", false, "print version")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: ws URL")
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetFlags(0)

	if *printVersion {
		fmt.Printf("ws %s\n", Version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if *origin == "" {
		originURL := *u
		switch originURL.Scheme {
		case "wss":
			originURL.Scheme = "https"
		default:
			originURL.Scheme = "http"
		}
		*origin = originURL.String()
	}

	err = echo(u.String(), *origin)
	switch err {
	case nil, io.EOF, liner.ErrPromptAborted:
	default:
		log.Fatal(err)
	}
}
