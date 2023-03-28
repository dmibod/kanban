package notify

import (
	"sync"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type client struct {
	sync.Mutex
	logger.Logger
	clientID int
	boardID  kernel.ID
}

func createClient(clientID int, l logger.Logger) *client {
	return &client{
		clientID: clientID,
		Logger:   l,
	}
}

func (c *client) receive(boardID kernel.ID) {
	c.Lock()
	c.boardID = boardID
	c.Unlock()

	c.Debugf("client %v switched to board %v\n", c.clientID, boardID)
}

func (c *client) send(notifications []kernel.Notification) []kernel.Notification {
	c.Lock()
	boardID := c.boardID
	c.Unlock()

	var send []kernel.Notification

	for _, n := range notifications {
		if isDeliverable(boardID, n) {
			send = append(send, n)
		} else {
			c.Debugf("client %v context %v != %v, ignore notification\n", c.clientID, boardID, n.BoardID)
		}
	}

	return send
}

func isDeliverable(id kernel.ID, n kernel.Notification) bool {
	if id == n.BoardID {
		return true
	}

	if id == "" {
		return n.Type == kernel.RefreshBoardNotification || n.Type == kernel.RemoveBoardNotification || n.Type == kernel.CreateBoardNotification
	}

	return false
}
