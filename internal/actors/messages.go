package actors

type SimpleMessage struct {
	SenderName string
	Content    string
}

type CloseActorMessage struct {
}
