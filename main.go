// Package main contains the main function
package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/structx/common/database"
	"github.com/structx/common/logging"
)

func main() {

	fx.New(
		fx.Provide(logging.NewZap),
		fx.Provide(database.NewPGXPool),
		fx.Invoke(registerHooks),
	).Run()
}

func registerHooks(lc fx.Lifecycle) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			return nil
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})

	return nil
}
