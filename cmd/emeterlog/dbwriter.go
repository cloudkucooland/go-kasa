package main

import (
	"context"
	"os"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/cloudkucooland/go-kasa"
	"github.com/urfave/cli/v3"
)

type emeterdata struct {
	DeviceID string
	Alias    string
	R        *kasa.EmeterRealtime
}

var (
	client   influxdb2.Client
	writeAPI api.WriteAPI
)

func setupdb(ctx context.Context, cmd *cli.Command) error {
	host := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")

	emlog.InfoContext(ctx, "Setting up Kasa InfluxDB V2 Client", "host", host, "bucket", bucket)

	client = influxdb2.NewClient(host, token)
	
	// Use the async WriteAPI to handle batching automatically
	writeAPI = client.WriteAPI(org, bucket)

	// Log async errors
	go func() {
		for err := range writeAPI.Errors() {
			emlog.Error("kasa influx write error", "err", err)
		}
	}()

	return nil
}

func startDBWriter(ctx context.Context, r <-chan emeterdata) {
	emlog.Info("Kasa DB Writer started (V2 Downgrade)")

	for {
		select {
		case <-ctx.Done():
			writeAPI.Flush()
			client.Close()
			return
		case v, ok := <-r:
			if !ok {
				return
			}
			
			// Explicit uint64 casts ensure the 'u' suffix is added in Line Protocol
			p := influxdb2.NewPoint("emeter",
				map[string]string{
					"device": v.DeviceID, 
					"alias":  v.Alias,
				},
				map[string]interface{}{
					"slot":      uint64(v.R.Slot),
					"VoltageMV": uint64(v.R.VoltageMV),
					"CurrentMA": uint64(v.R.CurrentMA),
					"PowerMW":   uint64(v.R.PowerMW),
				},
				time.Now())

			writeAPI.WritePoint(p)
		}
	}
}
