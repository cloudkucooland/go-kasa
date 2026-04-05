package main

import (
	"context"
	"os"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3/batching"
	"github.com/cloudkucooland/go-kasa"
	"github.com/urfave/cli/v3"
)

type emeterdata struct {
	DeviceID string
	Alias    string
	R        *kasa.EmeterRealtime
}

var client *influxdb3.Client

func setupdb(ctx context.Context, cmd *cli.Command) error {
	var err error

	if h := os.Getenv("INFLUX_HOST"); h != "" {
		log.InfoContext(ctx, "INFLUX_HOST", h)
	} else {
		log.InfoContext(ctx, "INFLUX_HOST not set")
	}
	if d := os.Getenv("INFLUX_DATABASE"); d != "" {
		log.InfoContext(ctx, "INFLUX_DATABASE", d)
	} else {
		log.InfoContext(ctx, "INFLUX_DATABASE not set")
	}

	client, err = influxdb3.NewFromEnv()
	if err != nil {
		return err
	}

	return nil
}

func dbwriter(ctx context.Context, r <-chan emeterdata) error {
	batch := batching.NewBatcher(batching.WithSize(80))

	for v := range r {
		p := influxdb3.NewPoint("emeter",
			map[string]string{
				"device": v.DeviceID,
				"alias":  v.Alias,
			},
			map[string]any{
				"slot":      v.R.Slot,
				"VoltageMV": v.R.VoltageMV,
				"CurrentMA": v.R.CurrentMA,
				"PowerMW":   v.R.PowerMW,
			},
			time.Now())

		batch.Add(p)
		if batch.Ready() {
			err := client.WritePoints(ctx, batch.Emit())
			if err != nil {
				continue
			}
		}
	}

	if err := client.WritePoints(ctx, batch.Emit()); err != nil {
		return err
	}

	return nil
}
