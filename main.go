package main

import (
	"go-actors/internal/actors"
	zerologger "go-actors/pkg/log"
	"time"
)

const scope = "main"

func main() {
	ctx, cancelFn := zerologger.InitLogger(scope, true)
	defer cancelFn()

	logger := zerologger.GetCtxLogger(ctx)
	logger.Info().Msg("ActorSystem ready")

	// wg := sync.WaitGroup{}
	// wg.Add(100_000)
	// msg := actors.SimpleMessage{Content: "Ping"}
	// for i := 0; i < 100_000; i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		actor := actors.NewActor(ctx, "Actor"+strconv.Itoa(i))
	// 		recipient := actors.NewActor(ctx, "Actor"+strconv.Itoa(i+1))
	// 		actor.SendMessage(ctx, msg, recipient)
	// 	}()
	// }

	// wg.Wait()

	// creating actors
	actor1 := actors.NewActor(ctx, "Actor1")
	actor2 := actors.NewActor(ctx, "Actor2")

	go actor1.SendMessage(ctx, actors.SimpleMessage{Content: "Ping"}, actor2)
	go actor2.SendMessage(ctx, actors.SimpleMessage{Content: "Hi there!"}, actor1)

	actor1.SendMessage(ctx, actors.CloseActorMessage{}, actor2)
	actor2.SendMessage(ctx, actors.CloseActorMessage{}, actor1)

	time.Sleep(3 * time.Second)

	actor1.SendMessage(ctx, actors.SimpleMessage{Content: "P1"}, actor2)
	actor2.SendMessage(ctx, actors.SimpleMessage{Content: "P2"}, actor1)

	time.Sleep(1 * time.Second)
}
