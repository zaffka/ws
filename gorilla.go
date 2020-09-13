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
	url    url.URL
	header http.Header

	conn   *websocket.Conn
	ctx    context.Context
	resCh  chan []byte
	finCh  chan struct{}
	finErr error
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

func (g *gorilla) Data() <-chan []byte {
	g.resCh = make(chan []byte)
	g.finCh = make(chan struct{})

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
				g.resCh <- message
			}
		}
		close(g.resCh)
		close(g.finCh)
	}()

	return g.resCh
}

func (g *gorilla) Done() <-chan struct{} {
	return g.finCh
}

func (g *gorilla) Err() error {
	return g.finErr
}
