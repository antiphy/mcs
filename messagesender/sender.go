package messagesender

// fake 3rd message api sender

type MessageSender interface {
	SendMsg(msg []byte, PhoneNumber string, name string) error
}

func NewMsgSender() MessageSender {
	return &sender{}
}

type sender struct {
}

// TODO:
func (s *sender) SendMsg(msg []byte, PhoneNumber string, name string) error {
	return nil
}
