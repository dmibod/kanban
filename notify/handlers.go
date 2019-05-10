package notify

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/message"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
)

// API holds dependencies required by handlers
type API struct {
	sync.Mutex
	logger.Logger
	key      int
	clients  map[int]kernel.ID
	channels map[int]chan []byte
}

// CreateAPI creates new API instance
func CreateAPI(s message.Subscriber, l logger.Logger) *API {
	api := &API{
		Logger:   l,
		clients:  make(map[int]kernel.ID),
		channels: make(map[int]chan []byte),
	}

	s.Subscribe(bus.HandleFunc(func(msg []byte) {
		api.Lock()
		defer api.Unlock()
		for _, q := range api.channels {
			q <- msg
		}
	}))

	return api
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/", a.HandleConnect)
}

// HandleConnect web socket
func (a *API) HandleConnect(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.Errorln(err)
		/*
			if _, ok := err.(websocket.HandshakeError); !ok {
				a.Errorln(err)
			}
		*/
		return
	}

	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ch, key := a.subscribe()

	go a.writer(ws, ch, key)

	for {
		t, payload, err := ws.ReadMessage()
		if err != nil {
			a.Errorln(err)
			break
		} else if t == websocket.TextMessage {
			msg := &struct {
				ID kernel.ID `json:"id"`
			}{}
			if err := json.Unmarshal(payload, &msg); err != nil {
				a.Errorln(err)
			} else {
				a.Debugf("client %v switched to board %v\n", key, msg.ID)
				a.Lock()
				a.clients[key] = msg.ID
				a.Unlock()
			}
		}
	}
}

func (a *API) subscribe() (<-chan []byte, int) {
	a.Lock()
	defer a.Unlock()
	a.key++
	ch := make(chan []byte)
	a.channels[a.key] = ch
	return ch, a.key
}

func (a *API) unsubscribe(key int) {
	a.Lock()
	defer a.Unlock()
	if ch, ok := a.channels[key]; ok {
		close(ch)
		delete(a.clients, key)
		delete(a.channels, key)
	}
}

func (a *API) writer(ws *websocket.Conn, q <-chan []byte, key int) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		a.Debugf("unsubscribe client %v and close socket\n", key)
		a.unsubscribe(key)
		pingTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case m := <-q:
			if err := a.onMessage(ws, m, key); err != nil {
				a.Errorln(err)
				return
			}
		case <-pingTicker.C:
			if err := a.onPing(ws); err != nil {
				a.Errorln(err)
				return
			}
		}
	}
}

func (a *API) onMessage(ws *websocket.Conn, m []byte, key int) error {
	received := []kernel.Notification{}
	if err := json.Unmarshal(m, &received); err != nil {
		return err
	}

	if len(received) == 0 {
		a.Debugf("client %v received 0 notifications, ignore processing\n", key)
		return nil
	}

	a.Lock()
	ctx, ok := a.clients[key]
	a.Unlock()

	if !ok {
		a.Debugf("client %v has not opened any board yet, ignore notification\n", key)
		return nil
	}

	send := []kernel.Notification{}
	for _, n := range received {
		if ctx != n.BoardID {
			a.Debugf("client %v context %v != %v, ignore notification\n", key, ctx, n.BoardID)
		} else {
			send = append(send, n)
		}
	}

	if len(send) == 0 {
		a.Debugf("client %v 0 notifications to deliver, ignore processing\n", key)
		return nil
	}

	out, err := json.Marshal(send)
	if err != nil {
		a.Errorln(err)
		return err
	}

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(websocket.TextMessage, out)
}

func (a *API) onPing(ws *websocket.Conn) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(websocket.PingMessage, []byte{})
}
