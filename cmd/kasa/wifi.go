package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/cloudkucooland/go-kasa"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

var setwifi = &cli.Command{
	Name:      "setwifi",
	Usage:     "configure wifi",
	ArgsUsage: "host ssid key",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "ssid"},
		&cli.StringArg{Name: "key"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		_, err := k.SetWIFICtx(ctx, cmd.StringArg("ssid"), cmd.StringArg("key"))
		return err
	},
}

var wifi = &cli.Command{
	Name:      "wifi",
	Usage:     "check device wifi status",
	ArgsUsage: "[host]",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		res, err := k.GetWIFIStatusCtx(ctx)
		if err != nil {
			return err
		}

		return formatOutput(cmd, res, func() {
			tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if !cmd.Bool("no-header") {
				fmt.Fprintf(tabwrite, "Device\tIP\tSSID\t%s\t%s\n", color.GreenString("Key Type"), color.GreenString("RSSI"))
			}
			fmt.Fprintf(tabwrite, "%s\t%s\t%s\t%s\t%s\n", "", k.IP, res.SSID, keyType(res.KeyType), colorRSSI(res.RSSI))
			_ = tabwrite.Flush()
		})
	},
}

type ws struct {
	Host    string
	StaInfo *kasa.StaInfo
	Sysinfo *kasa.Sysinfo
}

var allwifi = &cli.Command{
	Name:  "allwifi",
	Usage: "get wifi stats for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		bctx, cancel := context.WithTimeout(context.Background(), time.Duration(cmd.Int("timeout"))*time.Second)
		defer cancel()

		m, err := kasa.BroadcastWifiParameters(bctx, int(cmd.Int("repeats")))
		if err != nil {
			return err
		}

		w := make([]ws, 0)
		var mu sync.Mutex
		g, gctx := errgroup.WithContext(ctx)

		for host, info := range m {
			h, i := host, info // shadow for closure
			g.Go(func() error {
				kd, _ := kasa.NewDevice(h)
				s, err := kd.GetSettingsCtx(gctx)
				if err != nil {
					return nil // skip if offline
				}
				mu.Lock()
				w = append(w, ws{h, i, s})
				mu.Unlock()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		return formatOutput(cmd, w, func() {
			sort.Slice(w, func(i, j int) bool {
				return w[i].Sysinfo.Alias < w[j].Sysinfo.Alias
			})

			tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if !cmd.Bool("no-header") {
				fmt.Fprintf(tabwrite, "Device\tIP\tSSID\t%s\t%s\n", color.GreenString("Key Type"), color.GreenString("RSSI"))
			}
			for _, v := range w {
				fmt.Fprintf(tabwrite, "%s\t%s\t%s\t%s\t%s\n", v.Sysinfo.Alias, v.Host, v.StaInfo.SSID, keyType(v.StaInfo.KeyType), colorRSSI(v.StaInfo.RSSI))
			}
			_ = tabwrite.Flush()
		})
	},
}

func colorRSSI(rssi int) string {
	s := fmt.Sprintf("%ddB", rssi)
	switch {
	case rssi <= -80:
		return color.RedString(s)
	case rssi <= -70:
		return color.YellowString(s)
	default:
		return color.GreenString(s)
	}
}

func keyType(t int) string {
	switch t {
	case 3:
		return color.GreenString("WPA3")
	case 2:
		return color.YellowString("WPA2")
	default:
		return color.RedString("Unknown")
	}
}
