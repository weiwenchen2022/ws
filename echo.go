package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/peterh/liner"
	"nhooyr.io/websocket"
)

func echo(u, origin string) error {
	header := make(http.Header)
	header.Add("Origin", origin)

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	c, _, err := websocket.Dial(ctx, u, &websocket.DialOptions{
		HTTPHeader:   header,
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		return err
	}

	cl := &client{
		c:    c,
		errc: make(chan error, 1),
	}

	go cl.readStdin()
	go cl.readWs()

	return <-cl.errc
}

type client struct {
	c    *websocket.Conn
	errc chan error
}

func (cl *client) readStdin() {
	defer cl.c.CloseNow()

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	var historyFile string
	if homeDir, err := os.UserHomeDir(); err == nil {
		historyFile = filepath.Join(homeDir, ".ws_history")
	}

	if historyFile != "" {
		if f, err := os.Open(historyFile); err == nil {
			_, _ = line.ReadHistory(f)
			f.Close()
		}
	}

	for {
		l, err := line.Prompt("> ")
		if err != nil {
			cl.errc <- err
			break
		}

		line.AppendHistory(l)

		if err := cl.writeTimeout([]byte(l)); err != nil {
			cl.errc <- err
			break
		}
	}

	if historyFile != "" {
		if f, err := os.Create(historyFile); err != nil {
			log.Println("Error writing history file: ", err)
		} else {
			_, _ = line.WriteHistory(f)
			f.Close()
		}
	}
}

func (cl *client) readWs() {
	defer cl.c.CloseNow()

	for {
		typ, b, err := cl.c.Read(context.Background())
		if err != nil {
			cl.errc <- err
			return
		}

		var text string
		switch typ {
		case websocket.MessageText:
			text = string(b)
		case websocket.MessageBinary:
			text = fmt.Sprintf("% x", b)
		}
		log.Println(color.CyanString("< %s", text))
	}
}

func (cl *client) writeTimeout(b []byte) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	return cl.c.Write(ctx, websocket.MessageText, b)
}
