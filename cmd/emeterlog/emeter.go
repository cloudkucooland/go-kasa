package main

import (
	"context"
	"sync"
	"time"

	"github.com/cloudkucooland/go-kasa"
)

const timeout = 2 * time.Second
const repeats = 1

func queryall(ctx context.Context, results chan emeterdata) error {
	bctx, cancel := context.WithTimeout(ctx, timeout)
	m, err := kasa.BroadcastEmeter(bctx, repeats)
	if err != nil {
		return err
	}
	defer cancel()

	var wg sync.WaitGroup

	for k, v := range m {
		kd, err := kasa.NewDevice(k)
		if err != nil {
			return err
		}

		wg.Go(func() {
			s, err := kd.GetSettingsCtx(ctx)
			if err != nil {
				return
			}

			if s.NumChildren > 0 {
				for _, c := range s.Children {
					cv, err := kd.GetEmeterChildCtx(ctx, c.ID)
					if err != nil {
						return
					}
					results <- emeterdata{
						DeviceID: s.DeviceID,
						Alias:    c.Alias,
						R:        cv,
					}
				}
			} else {
				results <- emeterdata{
					DeviceID: s.DeviceID,
					Alias:    s.Alias,
					R:        &v.Emeter.Realtime,
				}
			}
		})
	}
	wg.Wait()
	return nil
}
