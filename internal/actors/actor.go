package actors

import (
	"context"
	zerologger "go-actors/pkg/log"
	"go-actors/pkg/utils"
	"log"
)

type Actor struct {
	Name         string
	state        string
	messageQueue chan any
}

func NewActor(ctx context.Context, name string) *Actor {
	ctx = utils.GetCtxWithScope(ctx, "Actor Constructor")
	logger := zerologger.GetCtxLogger(ctx)

	actor := &Actor{
		Name:         name,
		state:        "created",
		messageQueue: make(chan any),
	}

	go actor.start(ctx) // запуск обработчика сообщений

	logger.Debug().Str("name", actor.Name).Str("state", actor.state).Msgf("actor created")
	return actor
}

func (a *Actor) start(ctx context.Context) {
	for msg := range a.messageQueue {
		a.handleMessage(ctx, msg)
	}
}

func (a *Actor) handleMessage(ctx context.Context, msg any) {
	ctx = utils.GetCtxWithScope(ctx, "Actor Handler")
	logger := zerologger.GetCtxLogger(ctx)

	switch m := msg.(type) {
	case SimpleMessage:
		logger.Info().Msgf("%s received message: %s from %s", a.Name, m.Content, m.SenderName)
	case CloseActorMessage:
		logger.Info().Msgf("actor %s state: %s", a.Name, a.state)
		a.state = "closed"
	}
}

// sending messages to actor using actor system
func (senderActor Actor) SendMessage(ctx context.Context, msg any, receiverActor *Actor) {
	ctx = utils.GetCtxWithScope(ctx, "ActorSystem Sending Message")
	logger := zerologger.GetCtxLogger(ctx)

	if receiverActor != nil {
		if receiverActor.state == "closed" {
			logger.Info().Msgf("actor:%s we are closed! come tomorrow!!", receiverActor.Name)
			return
		}

		switch m := msg.(type) {
		case SimpleMessage:
			m.SenderName = senderActor.Name
			logger.Info().Msgf("actor:%s is sending message:%s to the actor:%s", senderActor.Name, m.Content, receiverActor.Name)
			receiverActor.messageQueue <- m
		case CloseActorMessage:
			receiverActor.messageQueue <- m
		}
	} else {
		log.Printf("Actor not found\n")
	}
}
