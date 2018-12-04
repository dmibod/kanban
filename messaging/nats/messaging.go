package nats

type Client struct{
}

func (c *Client) Send(subj string, msg []byte){
}

func (c *Client) Receive(group string, subj string, msg []byte){
}

func New() *Client {
  return &Client{}
}
