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
		results := make(chan emeterdata, 100)
		go startDBWriter(ctx, results)

		pollrate := cmd.Int("pollrate")
		if pollrate <= 0 {
			pollrate = 2
		}

		tt := cmd.Int("timeout")
		if tt > 0 {
			timeout = time.Duration(tt) * time.Second
		}

		rr := cmd.Int("repeats")
		if rr > 0 {
			repeats = rr
		}

		ticker := time.NewTicker(time.Duration(pollrate) * time.Second)

		for {
			select {
			case <-ctx.Done():
				close(results)
				return nil
			case <-ticker.C:
				// drop any lingering attempts before the next tick
				runCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
				defer cancel()
				if err := queryall(runCtx, results); err != nil {
					emlog.Error("query error", "err", err)
				}
			}
		}
		return nil
	},
}
