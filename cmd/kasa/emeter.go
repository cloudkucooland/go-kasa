package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var getallemeter = &cli.Command{
	Name:  "getallemeter",
	Usage: "get emeter stats for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		m, err := kasa.BroadcastEmeter(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		var tma uint
		var tw float64
		for k, v := range m {
			kd, err := kasa.NewDevice(k)
			if err != nil {
				return err
			}
			s, err := kd.GetSettingsCtx(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("[%s]\n", s.Alias)

			if s.NumChildren > 0 {
				var ma uint
				var w float64
				for _, c := range s.Children {
					cv, err := kd.GetEmeterChildCtx(ctx, c.ID)
					if err != nil {
						return err
					}
					ma += cv.CurrentMA
					w += float64(cv.PowerMW) / 1000
					fmt.Printf("[%s]\t", c.Alias)
					fmt.Printf("Current:\t%dmA\t", cv.CurrentMA)
					fmt.Printf("Voltage:\t%2.2fV\t", float64(cv.VoltageMV)/1000)
					fmt.Printf("Power:\t%2.2fW\t", float64(cv.PowerMW)/1000)
					fmt.Printf("Total:\t%2.2fkWh\n", float64(cv.TotalWH)/1000)
				}
				fmt.Printf("Total\tCurrent:\t%dmA\tPower:\t%2.2fW\n", ma, w)
				tma += ma
				tw += w
			} else {
				fmt.Printf("Current:\t%dmA\n", v.Emeter.Realtime.CurrentMA)
				fmt.Printf("Voltage:\t%2.2fV\n", float64(v.Emeter.Realtime.VoltageMV)/1000)
				fmt.Printf("Power:\t\t%2.2fW\n", float64(v.Emeter.Realtime.PowerMW)/1000)
				fmt.Printf("Total:\t\t%2.2fkWh\n", float64(v.Emeter.Realtime.TotalWH)/1000)
				tma += v.Emeter.Realtime.CurrentMA
				tw += float64(v.Emeter.Realtime.PowerMW) / 1000
			}
		}
		fmt.Printf("Total House Current:\t%dmA\tPower:\t%2.2fW\n", tma, tw)
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
		k, err := getKasaDevice(cmd)
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
		}

		if month == 0 {
			var em *kasa.EmeterRealtime

			if cmd.String("child") != "" {
				em, err = k.GetEmeterChildCtx(ctx, cmd.String("child"))
			} else {
				em, err = k.GetEmeterCtx(ctx)
			}

			if err != nil {
				return err
			}

			fmt.Printf("Current:\t%dmA\n", em.CurrentMA)
			fmt.Printf("Voltage:\t%2.2fV\n", float64(em.VoltageMV)/1000)
			fmt.Printf("Power:\t\t%2.2fW\n", float64(em.PowerMW)/1000)
			fmt.Printf("Total:\t\t%2.2fkWh\n", float64(em.TotalWH)/1000)
			return nil
		}

		if year == 0 {
			year = time.Now().Year()
		}
		// get month/year date range
		em, err := k.GetEmeterMonthCtx(ctx, month, year)
		if err != nil {
			return err
		}
		for _, v := range em.List {
			fmt.Printf("%d-%02d-%02d:\t%2.2fkWh\n", v.Year, v.Month, v.Day, float64(v.WH)/1000)
		}
		return nil
	},
}
