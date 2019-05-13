package notify

import (
	"encoding/json"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type connection struct {
	sync.Mutex
	logger.Logger
	socket  *websocket.Conn
	channel chan []byte
	client  *client
}

func createConnection(clientID int, socket *websocket.Conn, l logger.Logger) *connection {
	return &connection{
		Logger:  l,
		socket:  socket,
		channel: make(chan []byte),
		client:  createClient(clientID, l),
	}
}

func (c *connection) read() {
	ws := c.socket

	for {
		messageType, payload, err := ws.ReadMessage()
		if err != nil {
			c.Errorln(err)
			break
		} else if messageType == websocket.TextMessage {
			msg := &struct {
				ID kernel.ID `json:"id"`
			}{}
			if err := json.Unmarshal(payload, &msg); err != nil {
				c.Errorln(err)
			} else {
				c.client.receive(msg.ID)
			}
		}
	}
}

func (c *connection) write() {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		c.close()
	}()

	for {
		select {
		case m := <-c.channel:
			if err := c.message(m); err != nil {
				c.Errorln(err)
				return
			}
		case <-pingTicker.C:
			if err := c.ping(); err != nil {
				c.Errorln(err)
				return
			}
		}
	}
}

func (c *connection) close() {
	c.Lock()
	defer c.Unlock()

	ws := c.socket
	ch := c.channel

	c.Debugf("unsubscribe client %v and close socket\n", c.client.clientID)

	if ch != nil {
		close(ch)
		c.channel = nil
	}
	if ws != nil {
		ws.Close()
		c.socket = nil
	}
}

func (c *connection) message(m []byte) error {
	received := []kernel.Notification{}
	if err := json.Unmarshal(m, &received); err != nil {
		return err
	}

	if len(received) == 0 {
		c.Debugf("client %v received 0 notifications, ignore processing\n", c.client.clientID)
		return nil
	}

	send := c.client.send(received)

	if len(send) == 0 {
		c.Debugf("client %v 0 notifications to deliver, ignore processing\n", c.client.clientID)
		return nil
	}

	out, err := json.Marshal(send)
	if err != nil {
		c.Errorln(err)
		return err
	}

	ws := c.socket

	ws.SetWriteDeadline(time.Now().Add(writeWait))

	return ws.WriteMessage(websocket.TextMessage, out)
}

func (c *connection) ping() error {
	ws := c.socket

	ws.SetWriteDeadline(time.Now().Add(writeWait))

	return ws.WriteMessage(websocket.PingMessage, []byte{})
}
