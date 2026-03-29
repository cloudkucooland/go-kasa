package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var discover = &cli.Command{
	Name:  "discover",
	Usage: "discover local devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		m, err := kasa.BroadcastDiscovery(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
		if err != nil {
			return err
		}

		keys := make([]string, 0, len(m))
		for key := range m {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		fmt.Printf("found %d devices\n", len(m))
		for _, k := range keys {
			v := m[k]
			if len(v.Children) == 0 {
				fmt.Printf("%15s: %s %32s [state: %d] [brightness: %3d]\n", k, v.Model, v.Alias, v.RelayState, v.Brightness)
			} else {
				fmt.Printf("%15s: %s %s\n", k, v.Model, v.Alias)
				for _, c := range v.Children {
					fmt.Printf("    ID: %40s%s %26s [state: %d]\n", v.DeviceID, c.ID, c.Alias, c.RelayState)
				}
			}
		}
		return nil
	},
}
