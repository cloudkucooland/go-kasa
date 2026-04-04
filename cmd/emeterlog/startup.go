package main

import (
	"context"
	"time"

	"github.com/urfave/cli/v3"
)

var startup = &cli.Command{
	Name:  "startup",
	Usage: "start the daemon",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if err := setupdb(ctx, cmd); err != nil {
			return err
		}

		ticker := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				queryall(ctx)
			}
		}
		return nil
	},
}
