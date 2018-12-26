package notify

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/dmibod/kanban/shared/tools/msg"
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
	logger     logger.Logger
	subscriber msg.Subscriber
	queue      <-chan []byte
}

// CreateAPI creates new API instance
func CreateAPI(l logger.Logger, s msg.Subscriber) *API {
	q := make(chan []byte)
	_, err := s.Subscribe("", func(msg []byte) {
		q <- msg
	})
	if err != nil {
		l.Errorln("error subscribe queue", err)
	}
	return &API{
		logger:     l,
		subscriber: s,
		queue:      q,
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
			a.logger.Errorln("error reading message", err)
			break
		}
	}
}

func (a *API) writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case m := <-a.queue:
			n := Notification{}
			err := json.Unmarshal(m, &n)
			if err != nil {
				a.logger.Errorln("error parsing json", err)
				return
			} else {
				a.logger.Debugln(n)
			}
			if len(n) > 0 {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
					a.logger.Errorln("error writing message", err)
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				a.logger.Errorln("error ping message", err)
				return
			}
		}
	}
}

func (a *API) Ws(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			a.logger.Errorln(err)
		}

		return
	}

	go a.writer(ws)

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
