package notify

import (
	"net/http"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/message"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 30 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
)

// API holds dependencies required by handlers
type API struct {
	sync.Mutex
	logger.Logger
	key         int
	connections map[int]*connection
}

// CreateAPI creates new API instance
func CreateAPI(s message.Subscriber, l logger.Logger) *API {
	api := &API{
		Logger:      l,
		connections: make(map[int]*connection),
	}

	s.Subscribe(bus.HandleFunc(func(msg []byte) {
		api.Lock()
		defer api.Unlock()
		l.Debugf("broadcast msg: %v\n", msg)
		for _, conn := range api.connections {
			conn.channel <- msg
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

	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	conn := a.subscribe(ws)

	go conn.write()

	conn.read()

	a.unsubscribe(conn)
}

func (a *API) subscribe(ws *websocket.Conn) *connection {
	a.Lock()
	defer a.Unlock()
	a.key++
	conn := createConnection(a.key, ws, a.Logger)
	a.connections[a.key] = conn
	return conn
}

func (a *API) unsubscribe(conn *connection) {
	a.Lock()
	defer a.Unlock()
	key := conn.client.clientID
	if conn, ok := a.connections[key]; ok {
		conn.close()
		delete(a.connections, key)
	}
}
