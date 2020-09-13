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
	doneCh chan error
}

func (g *gorilla) Conn(ctx context.Context) (func() error, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, g.url.String(), g.header)
	if err != nil {
		return nil, err
	}
	g.conn = conn

	return func() error {
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

func (g *gorilla) unsubscribe() {
	g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (g *gorilla) handle(ctx context.Context, res chan []byte) {
	defer close(g.doneCh)
	defer g.unsubscribe()

handleLoop:
	for {
		_, message, err := g.conn.ReadMessage()
		if err != nil {
			g.doneCh <- fmt.Errorf("read error: %w", err)
			break handleLoop
		}

		select {
		case <-ctx.Done():
			g.doneCh <- ctx.Err()
			break handleLoop
		default:
			res <- message
		}
	}
}

func (g *gorilla) Handle(ctx context.Context) <-chan []byte {
	res := make(chan []byte)
	go g.handle(ctx, res)
	return res
}

func (g *gorilla) Done() <-chan error {
	return g.doneCh
}
