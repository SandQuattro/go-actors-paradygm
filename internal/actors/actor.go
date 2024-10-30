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

func NewActor(ctx context.Context, name string) Actor {
	ctx = utils.GetCtxWithScope(ctx, "Actor Constructor")
	logger := zerologger.GetCtxLogger(ctx)

	actor := Actor{
		Name:         name,
		state:        "created",
		messageQueue: make(chan any),
	}

	go actor.start(ctx) // запуск обработчика сообщений

	logger.Debug().Str("name", actor.Name).Str("state", actor.state).Msgf("actor created")
	return actor
}

// Обработчик сообщений актора
func (a Actor) start(ctx context.Context) {
	for msg := range a.messageQueue {
		a.handleMessage(ctx, msg)
	}
}

// Обработка сообщения и изменение состояния
func (a *Actor) handleMessage(ctx context.Context, msg any) {
	ctx = utils.GetCtxWithScope(ctx, "Actor Handler")
	logger := zerologger.GetCtxLogger(ctx)

	switch m := msg.(type) {
	case SimpleMessage:
		logger.Info().Msgf("%s received message: %s from %s", a.Name, m.Content, m.SenderName)
		a.state = "closed"
	case StateMessage:
		logger.Info().Msgf("actor %s state: %s", a.Name, a.state)
	}
}

// sending messages to actor using actor system
func (senderActor Actor) SendMessage(ctx context.Context, msg any, receiverActor *Actor) {
	ctx = utils.GetCtxWithScope(ctx, "ActorSystem Sending Message")
	logger := zerologger.GetCtxLogger(ctx)

	if receiverActor != nil && receiverActor.state != "closed" {
		switch m := msg.(type) {
		case SimpleMessage:
			m.SenderName = senderActor.Name
			logger.Info().Msgf("actor:%s sending message:%s to actor:%s", senderActor.Name, m.Content, receiverActor.Name)
			receiverActor.messageQueue <- m
		case StateMessage:
			receiverActor.messageQueue <- m
		}
	} else {
		log.Printf("Actor with name '%s' not found\n", receiverActor.Name)
	}
}
