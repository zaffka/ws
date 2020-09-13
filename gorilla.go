package ws

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type gorilla struct {
	conn   *websocket.Conn
	ctx    context.Context
	finErr error
	url    url.URL
	header http.Header
}

func (g *gorilla) Conn(ctx context.Context) (func() error, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, g.url.String(), g.header)
	if err != nil {
		return nil, err
	}

	g.conn = conn
	g.ctx = ctx

	return func() error {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			g.finErr = fmt.Errorf("write error: %w", err)
		}
		return g.conn.Close()
	}, nil
}

func (g *gorilla) Write(subscribeMsgB []byte) (i int, err error) {
	var wcl io.WriteCloser
	wcl, err = g.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	i, err = wcl.Write(subscribeMsgB)
	if err != nil {
		return
	}

	err = wcl.Close()
	if err != nil {
		return
	}

	return
}

func (g *gorilla) DataChan() <-chan []byte {
	resChan := make(chan []byte)

	go func() {
	handleLoop:
		for {
			select {
			case <-g.ctx.Done():
				break handleLoop
			default:
				_, message, err := g.conn.ReadMessage()
				if err != nil {
					g.finErr = fmt.Errorf("read error: %w", err)
					break handleLoop
				}
				resChan <- message
			}
		}
		close(resChan)
	}()

	return resChan
}

func (g *gorilla) Err() error {
	return g.finErr
}
