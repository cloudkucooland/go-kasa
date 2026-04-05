package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var discover = &cli.Command{
	Name:  "discover",
	Usage: "discover local devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		bctx, cancel := context.WithTimeout(ctx, time.Duration(cmd.Int("timeout"))*time.Second)
		m, err := kasa.BroadcastDiscovery(bctx, int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		defer cancel()

		keys := make([]string, 0, len(m))
		for key := range m {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabwrite, "Device\tIP/ID:\tModel\tState\tBrightness\n")
		for _, k := range keys {
			v := m[k]
			if len(v.Children) == 0 {
				fmt.Fprintf(tabwrite, "%s\t%s\t%s\t%s\t\t%3d\n", v.Alias, k, v.Model, i2o(v.RelayState), v.Brightness)
			} else {
				fmt.Fprintf(tabwrite, "%s\t%s\t%s\t\t\n", v.Alias, k, v.Model)
				for _, c := range v.Children {
					fmt.Fprintf(tabwrite, "%s/%s\t%s%s\t\t%s\t\n", v.Alias, c.Alias, v.DeviceID, c.ID, i2o(c.RelayState))
				}
			}
		}
		tabwrite.Flush()
		fmt.Fprintf(os.Stderr, "found %d devices\n", len(m))
		return nil
	},
}

func i2o(o uint) string {
	if o > 0 {
		return "On"
	}
	return "Off"
}
