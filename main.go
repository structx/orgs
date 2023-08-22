// Package main contains the main function
package main

import (
	"context"
	"fmt"

	"github.com/structx/common/database"
	"github.com/structx/common/logging"
	"github.com/structx/orgs/internal/pubsub"
	"go.uber.org/fx"
)

func main() {

	fx.New(
		fx.Provide(logging.NewZap),
		fx.Provide(database.NewPGXPool),
		fx.Provide(pubsub.New),
		fx.Invoke(registerHooks),
	).Run()
}

func registerHooks(lc fx.Lifecycle, client *pubsub.Client) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			// monitor publication channel
			err := client.Monitor(ctx)
			if err != nil {
				return fmt.Errorf("failed to monitor publication channel: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {

			// close client
			client.Close()

			return nil
		},
	})

	return nil
}
