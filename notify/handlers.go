package notify

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sync"
	"time"

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
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
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
func CreateAPI(l logger.Logger) *API {
	api := &API{
		Logger:   l,
		channels: make(map[int]chan []byte),
	}

	bus.Subscribe("notification", bus.HandleFunc(func(msg []byte) {
		api.Lock()
		defer api.Unlock()
		for _, q := range api.channels {
			q <- msg
		}
	}))

	return api
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

// Routes export API router
func (a *API) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", a.Home)
	router.Get("/ws", a.Ws)
	return router
}

func (a *API) reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			a.Errorln("error reading message", err)
			break
		}
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
			n := Notification{}
			err := json.Unmarshal(m, &n)
			if err != nil {
				a.Errorln("error parsing json", err)
				return
			}
			a.Debugln(n)
			if len(n) > 0 {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
					a.Errorln("error writing message", err)
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				a.Errorln("error ping message", err)
				return
			}
		}
	}
}

func (a *API) Ws(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			a.Errorln(err)
		}

		return
	}

	ch, key := a.subscribe()

	go a.writer(ws, ch, key)

	a.reader(ws)
}

func (a *API) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		"",
	}
	homeTempl.Execute(w, &v)
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Notifications</title>
    </head>
    <body>
        <pre id="data">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("data");
                var conn = new WebSocket("ws://{{.Host}}/v1/api/notify/ws");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('notification received');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
