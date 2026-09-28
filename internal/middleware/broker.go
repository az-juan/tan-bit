package middleware

import (
	"charm.land/log/v2"
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

func ExecBroker() {
	nc, err := nats.Connect("nats://broker:4222")
	if err != nil {
		log.Fatalf("Error en broker: %v", err)
	}
	defer nc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	js, _ := jetstream.New(nc)

	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ARTICULOS",
		Subjects: []string{"articulos.>"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created stream: %s\n", stream.CachedInfo().Config.Name)

	stream, err = js.Stream(ctx, "ARTICULOS")
	if err != nil {
		panic(err)
	}

	info, err := stream.Info(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name:     %s\n", info.Config.Name)
	fmt.Printf("Subjects: %v\n", info.Config.Subjects)
	fmt.Printf("Messages: %d\n", info.State.Msgs)
}
