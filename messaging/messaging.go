package messaging

type Client interface{
  Send(subj string, msg []byte)
  Receive(group string, subj string, msg []byte)
}