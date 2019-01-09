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

type Notification map[kernel.Id]int

// API holds dependencies required by handlers
type API struct {
	sync.Mutex
	logger.Logger
	key      int
	channels map[int]chan []byte
}

// CreateAPI creates new API instance
func CreateAPI(s message.Subscriber, l logger.Logger) *API {
	api := &API{
		Logger:   l,
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
		if _, ok := err.(websocket.HandshakeError); !ok {
			a.Errorln(err)
		}

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
		_, _, err := ws.ReadMessage()
		if err != nil {
			a.Errorln("error reading message", err)
			break
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
		delete(a.channels, key)
	}
}

func (a *API) writer(ws *websocket.Conn, q <-chan []byte, key int) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		a.unsubscribe(key)
		pingTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case m := <-q:
			if err := onMessage(ws, m); err != nil {
				a.Errorln(err)
				return
			}
		case <-pingTicker.C:
			if err := onPing(ws); err != nil {
				a.Errorln(err)
				return
			}
		}
	}
}

func onMessage(ws *websocket.Conn, m []byte) error {
	n := Notification{}
	if err := json.Unmarshal(m, &n); err != nil {
		return err
	}
	if len(n) > 0 {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		return ws.WriteMessage(websocket.TextMessage, m)
	}
	return nil
}

func onPing(ws *websocket.Conn) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(websocket.PingMessage, []byte{})
}
