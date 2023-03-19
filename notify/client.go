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

	send := []kernel.Notification{}

	for _, n := range notifications {
		if boardID == n.BoardID || n.Type == kernel.RefreshBoardNotification || n.Type == kernel.RemoveBoardNotification || n.Type == kernel.CreateBoardNotification {
			send = append(send, n)
		} else {
			c.Debugf("client %v context %v != %v, ignore notification\n", c.clientID, boardID, n.BoardID)
		}
	}

	return send
}
