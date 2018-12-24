package msg

type Receive func([]byte)

type Receiver interface {
	Receive(string, Receive) error
}

type Sender interface {
	Send([]byte) error
}

type Transport interface {
	CreateReceiver(string) Receiver

	CreateSender(string) Sender
}
