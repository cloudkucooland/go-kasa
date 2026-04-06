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

type DimmerResult struct {
	Host string `json:"host"`
	kasa.Sysinfo
	Realtime []Realtime
}

type Realtime struct {
	kasa.Child
	kasa.EmeterRealtime
}

var allemeter = &cli.Command{
	Name:  "allemeter",
	Usage: "get emeter stats for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		bctx, cancel := context.WithTimeout(ctx, time.Duration(cmd.Int("timeout"))*time.Second)
		defer cancel()

		m, err := kasa.BroadcastEmeter(bctx, int(cmd.Int("repeats")))
		if err != nil {
			return err
		}

		var d []DimmerResult
		var mu sync.Mutex
		g, gctx := errgroup.WithContext(ctx)

		for k, v := range m {
			kk, vv := k, v
			g.Go(func() error {
				dr := DimmerResult{
					Host: kk,
				}
				kd, err := kasa.NewDevice(kk)
				if err != nil {
					return err
				}

				drs, err := kd.GetSettingsCtx(gctx)
				if err != nil {
					return err
				}
				dr.Sysinfo = *drs

				if drs.NumChildren > 0 {
					for _, c := range drs.Children {
						drc, err := kd.GetEmeterChildCtx(gctx, c.ID)
						if err != nil {
							continue
						}
						dr.Realtime = append(dr.Realtime, Realtime{c, *drc})
					}
				} else {
					dr.Realtime = append(dr.Realtime, Realtime{kasa.Child{}, vv.Emeter.Realtime})
				}
				mu.Lock()
				d = append(d, dr)
				mu.Unlock()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		return formatOutput(cmd, d, func() {
			sort.Slice(d, func(i, j int) bool {
				return d[i].Alias < d[j].Alias
			})

			var tma, twh, tw uint // total MA, Wh, W and
			tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			if !cmd.Bool("no-header") {
				fmt.Fprintf(tabwrite, "Device\tCurrent\t%s\tPower\tSince Reset\n", color.GreenString("Voltage"))
			}

			for _, s := range d {
				var ma, w, tsr uint
				for _, cv := range s.Realtime {
					displayName := s.Alias
					if cv.Alias != "" {
						displayName = fmt.Sprintf("%s/%s", s.Alias, cv.Alias)
					}
					fmt.Fprintf(tabwrite, "%s\t%dmA\t%s\t%2.2fW\t%2.2fkWh\n", displayName, cv.CurrentMA, colorVolts(cv.VoltageMV), float64(cv.PowerMW)/1000, float64(cv.TotalWH)/1000)
					ma += cv.CurrentMA
					w += cv.PowerMW
					tsr += cv.TotalWH
				}
				tma += ma
				twh += tsr
				tw += w
			}
			fmt.Fprintf(tabwrite, "Total House\t%dmA\t%s\t%2.2fW\t%2.2fkWh\n", tma, color.GreenString(" "), float64(tw)/1000, float64(twh)/1000)
			_ = tabwrite.Flush()
		})
	},
}

var emeter = &cli.Command{
	Name:      "emeter",
	Usage:     "check energy usage",
	ArgsUsage: "[host] [month] [year]",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "month"},
		&cli.IntArg{Name: "year"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		s, err := k.GetSettingsCtx(ctx)
		if err != nil {
			return err
		}

		month := 0
		mm := cmd.IntArg("month")
		if mm != 0 {
			if mm < 1 || mm > 12 {
				return fmt.Errorf("invalid month")
			}
			month = mm
		}

		year := time.Now().Year()
		yy := cmd.IntArg("year")
		if yy != 0 {
			year = yy
		}

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		if month == 0 {
			fmt.Fprintf(tabwrite, "Device\tCurrent\t%s\tPower\tSince Reset\n", color.GreenString("Voltage"))
			child := cmd.String("child")
			if child == "" {
				var ma uint
				var w, twh float64
				for _, c := range s.Children {
					cv, err := k.GetEmeterChildCtx(ctx, c.ID)
					if err != nil {
						continue
					}
					ma += cv.CurrentMA
					w += float64(cv.PowerMW) / 1000
					twh += float64(cv.TotalWH) / 1000
					fmt.Fprintf(tabwrite, "%s\t", c.Alias)
					fmt.Fprintf(tabwrite, "%dmA\t", cv.CurrentMA)
					fmt.Fprintf(tabwrite, "%s\t", colorVolts(cv.VoltageMV))
					fmt.Fprintf(tabwrite, "%2.2fW\t", float64(cv.PowerMW)/1000)
					fmt.Fprintf(tabwrite, "%2.2fkWh\n", float64(cv.TotalWH)/1000)
				}
				fmt.Fprintf(tabwrite, "Total\t%dmA\t%s\t%2.2fW\t%2.2fkWh\n", ma, color.GreenString(" "), w, twh)
			} else {
				var em *kasa.EmeterRealtime

				if child != "" {
					em, err = k.GetEmeterChildCtx(ctx, child)
				} else {
					em, err = k.GetEmeterCtx(ctx)
				}
				if err != nil {
					return err
				}
				fmt.Fprintf(tabwrite, "%s\t%dmA\t", s.Alias, em.CurrentMA)
				fmt.Fprintf(tabwrite, "%s\t", colorVolts(em.VoltageMV))
				fmt.Fprintf(tabwrite, "%2.2fW\t", float64(em.PowerMW)/1000)
				fmt.Fprintf(tabwrite, "%2.2fkWh\n", float64(em.TotalWH)/1000)
			}
			_ = tabwrite.Flush()
			return nil
		}

		// get month/year date range
		if s.NumChildren > 0 {
			var stripTotal uint
			for _, c := range s.Children {
				fmt.Fprintf(tabwrite, "%s\t\n", c.Alias)
				em, err := k.GetEmeterChildMonthCtx(ctx, month, year, c.ID)
				if err != nil {
					continue
				}
				var plugTotal uint
				for _, v := range em.List {
					fmt.Fprintf(tabwrite, "%d-%02d-%02d:\t%dWh\n", v.Year, v.Month, v.Day, v.WH)
					plugTotal += v.WH
				}
				fmt.Fprintf(tabwrite, "Plug Total:\t%dWh\n", plugTotal)
				stripTotal += plugTotal
			}
			fmt.Fprintf(tabwrite, "Strip Total:\t%dWh\n", stripTotal)
		} else {
			em, err := k.GetEmeterMonthCtx(ctx, month, year)
			if err != nil {
				return err
			}
			for _, v := range em.List {
				fmt.Fprintf(tabwrite, "%d-%02d-%02d:\t%dWh\n", v.Year, v.Month, v.Day, v.WH)
			}
		}
		_ = tabwrite.Flush()
		return nil
	},
}

func colorVolts(mv uint) string {
	vVolts := float64(mv) / 1000.0
	vStr := fmt.Sprintf("%2.2fV", vVolts)
	coloredVolt := ""
	switch {
	case vVolts > 127.0:
		coloredVolt = color.RedString(vStr) // High Voltage Alarm
	case (vVolts > 124.0 && vVolts <= 127.0):
		coloredVolt = color.YellowString(vStr)
	case vVolts < 114.0:
		coloredVolt = color.RedString(vStr) // Sag/Under-voltage
	case (vVolts < 116.0 && vVolts >= 114.0):
		coloredVolt = color.YellowString(vStr)
	default:
		coloredVolt = color.GreenString(vStr) // Nominal
	}
	return coloredVolt
}
