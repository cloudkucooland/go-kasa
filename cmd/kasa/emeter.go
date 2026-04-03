package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var getallemeter = &cli.Command{
	Name:  "getallemeter",
	Usage: "get emeter stats for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		m, err := kasa.BroadcastEmeter(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		var tma, twh, tw uint // total MA, Wh, W and
		for k, v := range m {
			// ctx is already cancelled by now
			nctx := context.Background()

			kd, err := kasa.NewDevice(k)
			if err != nil {
				return err
			}
			s, err := kd.GetSettingsCtx(nctx)
			if err != nil {
				continue
				// return err
			}
			fmt.Fprintf(tabwrite, "[%s]\t\t\t\t\t\t\t\t\n", s.Alias)

			if s.NumChildren > 0 {
				var ma, w, todaytwh uint
				for _, c := range s.Children {
					cv, err := kd.GetEmeterChildCtx(nctx, c.ID)
					if err != nil {
						continue
						// return err
					}
					ma += cv.CurrentMA
					w += cv.PowerMW
					fmt.Fprintf(tabwrite, "[%s]\t", c.Alias)
					fmt.Fprintf(tabwrite, "Current:\t%dmA\t", cv.CurrentMA)
					fmt.Fprintf(tabwrite, "Voltage:\t%2.2fV\t", float64(cv.VoltageMV)/1000)
					fmt.Fprintf(tabwrite, "Power:\t%2.2fW\t", float64(cv.PowerMW)/1000)
					fmt.Fprintf(tabwrite, "Today:\t%2.2fkWh\n", float64(cv.TotalWH)/1000)
					todaytwh += cv.TotalWH
				}
				fmt.Fprintf(tabwrite, "Total\tCurrent:\t%dmA\t\t\tPower:\t%2.2fW\tTotal:\t%2.2fkWh\n", ma, float64(w)/1000, float64(todaytwh)/1000)
				tma += ma
				tw += todaytwh
			} else {
				fmt.Fprintf(tabwrite, "\tCurrent:\t%dmA\t", v.Emeter.Realtime.CurrentMA)
				fmt.Fprintf(tabwrite, "Voltage:\t%2.2fV\t", float64(v.Emeter.Realtime.VoltageMV)/1000)
				fmt.Fprintf(tabwrite, "Power:\t%2.2fW\t", float64(v.Emeter.Realtime.PowerMW)/1000)
				fmt.Fprintf(tabwrite, "Today:\t%2.2fkWh\n", float64(v.Emeter.Realtime.TotalWH)/1000)
				tma += v.Emeter.Realtime.CurrentMA
				tw += v.Emeter.Realtime.PowerMW
				twh += v.Emeter.Realtime.TotalWH
			}
		}
		fmt.Fprintf(tabwrite, "Total House\tCurrent:\t%dmA\t\t\tPower:\t%2.2fW\tTotal:\t%2.2fkWh\n", tma, float64(tw/1000), float64(twh/1000))
		tabwrite.Flush()
		return nil
	},
}

var emeter = &cli.Command{
	Name:      "emeter",
	Usage:     "check energy usage",
	ArgsUsage: "host month year",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host", Destination: &host},
		&cli.StringArg{Name: "month", Destination: &value},
		&cli.StringArg{Name: "year", Destination: &secondary},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		k, err := getKasaDevice(cmd)
		if err != nil {
			return err
		}
		s, err := k.GetSettingsCtx(ctx)
		if err != nil {
			return err
		}

		month := 0
		year := 0

		if value != "" {
			month, err = strconv.Atoi(value)
			if err != nil {
				return err
			}
			if month < 1 || month > 12 {
				return fmt.Errorf("invalid month")
			}
		}

		yy := cmd.String("year")
		if yy != "" {
			year, err = strconv.Atoi(yy)
			if err != nil {
				return err
			}
		} else {
			year = time.Now().Year()
		}

		if month == 0 {
			if cmd.String("child") == "" && s.NumChildren > 0 {
				var ma uint
				var w float64
				for _, c := range s.Children {
					cv, err := k.GetEmeterChildCtx(ctx, c.ID)
					if err != nil {
						continue
						// return err
					}
					ma += cv.CurrentMA
					w += float64(cv.PowerMW) / 1000
					fmt.Fprintf(tabwrite, "[%s]\t", c.Alias)
					fmt.Fprintf(tabwrite, "Current:\t%dmA\t", cv.CurrentMA)
					fmt.Fprintf(tabwrite, "Voltage:\t%2.2fV\t", float64(cv.VoltageMV)/1000)
					fmt.Fprintf(tabwrite, "Power:\t%2.2fW\t", float64(cv.PowerMW)/1000)
					fmt.Fprintf(tabwrite, "Total:\t%2.2fkWh\n", float64(cv.TotalWH)/1000)
				}
				fmt.Fprintf(tabwrite, "Total\tCurrent:\t%dmA\tPower:\t%2.2fW\n", ma, w)
			} else {
				var em *kasa.EmeterRealtime

				if cmd.String("child") != "" {
					em, err = k.GetEmeterChildCtx(ctx, cmd.String("child"))
				} else {
					em, err = k.GetEmeterCtx(ctx)
				}
				if err != nil {
					return err
				}
				fmt.Fprintf(tabwrite, "Current:\t%dmA\t", em.CurrentMA)
				fmt.Fprintf(tabwrite, "Voltage:\t%2.2fV\t", float64(em.VoltageMV)/1000)
				fmt.Fprintf(tabwrite, "Power:\t%2.2fW\t", float64(em.PowerMW)/1000)
				fmt.Fprintf(tabwrite, "Total:\t%2.2fkWh\n", float64(em.TotalWH)/1000)
			}
			return nil
		}

		// get month/year date range
		if s.NumChildren > 0 {
			var stripTotal uint
			for _, c := range s.Children {
				fmt.Fprintf(tabwrite, "[%s]\t\n", c.Alias)
				em, err := k.GetEmeterChildMonthCtx(ctx, month, year, c.ID)
				if err != nil {
					// continue
					return err
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
		tabwrite.Flush()
		return nil
	},
}
