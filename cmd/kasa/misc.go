package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudkucooland/go-kasa"
	"github.com/urfave/cli/v3"
)

var nocloud = &cli.Command{
	Name:      "nocloud",
	Usage:     "disable the TP-Link cloud connection",
	UsageText: "kasa nocloud host",
	ArgsUsage: "host",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.DisableCloudCtx(ctx)
	},
}

var cloud = &cli.Command{
	Name:      "cloud",
	Usage:     "configure the TP-Link cloud connection",
	UsageText: "kasa cloud host username password",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.EnableCloudCtx(ctx, cmd.Args().Get(1), cmd.Args().Get(2))
	},
}

var cloudinfo = &cli.Command{
	Name:      "cloudinfo",
	Usage:     "get TP-Link cloud connection info",
	UsageText: "kasa cloudinfo host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		info, err := k.GetCloudInfoCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Username: %s\nServer: %s\nBind: %d\n", info.Username, info.Server, info.Bind)
		return nil
	},
}

var firmwarelist = &cli.Command{
	Name:      "firmwarelist",
	Usage:     "get available firmware updates",
	UsageText: "kasa firmwarelist host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		fw, err := k.GetFirmwareListCtx(ctx)
		if err != nil {
			return err
		}
		for _, f := range fw {
			fmt.Printf("Version: %s\nRelease: %s\n", f.Ver, f.Rel)
		}
		return nil
	},
}

var gettime = &cli.Command{
	Name:      "gettime",
	Usage:     "get device time",
	UsageText: "kasa gettime host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		t, err := k.GetTimeCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year, t.Month, t.Mday, t.Hour, t.Min, t.Sec)
		return nil
	},
}

var gettimezone = &cli.Command{
	Name:      "gettimezone",
	Usage:     "get device timezone info",
	UsageText: "kasa gettimezone host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		t, err := k.GetTimezoneCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year, t.Month, t.Mday, t.Hour, t.Min, t.Sec)
		return nil
	},
}

var getsensorroutines = &cli.Command{
	Name:      "getsensorroutines",
	Usage:     "get all sensor routines",
	UsageText: "kasa getsensorroutines host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		routines, err := k.GetSensorRoutinesCtx(ctx)
		if err != nil {
			return err
		}
		for _, r := range routines {
			fmt.Printf("ID: %s, Name: %s, Enabled: %d\n", r.ID, r.Name, r.En)
		}
		return nil
	},
}

var deletesensorroutine = &cli.Command{
	Name:      "deletesensorroutine",
	Usage:     "delete a sensor routine",
	UsageText: "kasa deletesensorroutine host id",
	Before:    RequireDevice,
	ArgsUsage: "host id",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "id"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.DeleteSensorRoutineCtx(ctx, cmd.StringArg("id"))
	},
}

var getmanualaction = &cli.Command{
	Name:      "getmanualaction",
	Usage:     "get default manual action",
	UsageText: "kasa getmanualaction host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		a, err := k.GetDefaultManualActionCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("OffToS: %d\n", a)
		return nil
	},
}

var sensormode = &cli.Command{
	Name:      "sensormode",
	Usage:     "get sensor mode",
	UsageText: "kasa sensormode host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		m, err := k.GetSensorModeCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Mode: %s, bAuto: %d\n", m.Mode, m.BAuto)
		return nil
	},
}

var setsensormode = &cli.Command{
	Name:      "setsensormode",
	Usage:     "set sensor mode",
	UsageText: "kasa setsensormode host mode",
	Before:    RequireDevice,
	ArgsUsage: "host mode",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "mode"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetSensorModeCtx(ctx, cmd.StringArg("mode"))
	},
}

var setmanualaction = &cli.Command{
	Name:      "setmanualaction",
	Usage:     "set default manual action",
	UsageText: "kasa setmanualaction host offToS",
	Before:    RequireDevice,
	ArgsUsage: "host offToS",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "offToS"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetDefaultManualActionCtx(ctx, int(cmd.IntArg("offToS")))
	},
}

var getschedulerules = &cli.Command{
	Name:      "getschedulerules",
	Usage:     "get all schedule rules",
	UsageText: "kasa getschedulerules host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		rules, err := k.GetScheduleRulesCtx(ctx)
		if err != nil {
			return err
		}
		for _, r := range rules {
			fmt.Printf("ID: %s, Name: %s, Enabled: %d\n", r.ID, r.Name, r.Enable)
		}
		return nil
	},
}

var deleteschedulerule = &cli.Command{
	Name:      "deleteschedulerule",
	Usage:     "delete a schedule rule",
	UsageText: "kasa deleteschedulerule host id",
	Before:    RequireDevice,
	ArgsUsage: "host id",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "id"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.DeleteScheduleRuleCtx(ctx, cmd.StringArg("id"))
	},
}

var deleteallschedulerules = &cli.Command{
	Name:      "deleteallschedulerules",
	Usage:     "delete all schedule rules",
	UsageText: "kasa deleteallschedulerules host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.DeleteAllScheduleRulesCtx(ctx)
	},
}

var setscheduleenabled = &cli.Command{
	Name:      "setscheduleenabled",
	Usage:     "enable or disable the schedule",
	UsageText: "kasa setscheduleenabled host true|false",
	Before:    RequireDevice,
	ArgsUsage: "host state",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "state"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		b, err := strconv.ParseBool(cmd.StringArg("state"))
		if err != nil {
			return err
		}
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetScheduleEnabledCtx(ctx, b)
	},
}

var bulbstate = &cli.Command{
	Name:      "bulbstate",
	Usage:     "get bulb state",
	UsageText: "kasa bulbstate host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		s, err := k.GetLightStateCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("On: %d, Brightness: %d, Hue: %d, Sat: %d, CT: %d\n", s.OnOff, s.Brightness, s.Hue, s.Saturation, s.ColorTemp)
		return nil
	},
}

var bulbon = &cli.Command{
	Name:      "bulbon",
	Usage:     "turn bulb on",
	UsageText: "kasa bulbon host [brightness] [hue] [sat] [temp]",
	Before:    RequireDevice,
	ArgsUsage: "host [brightness] [hue] [sat] [temp]",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		// default to 50% brightness, 0 hue, 0 sat, 2700 temp
		b := 50
		h := 0
		s := 0
		t := 2700
		if cmd.Args().Len() > 1 {
			b, _ = strconv.Atoi(cmd.Args().Get(1))
		}
		if cmd.Args().Len() > 2 {
			h, _ = strconv.Atoi(cmd.Args().Get(2))
		}
		if cmd.Args().Len() > 3 {
			s, _ = strconv.Atoi(cmd.Args().Get(3))
		}
		if cmd.Args().Len() > 4 {
			t, _ = strconv.Atoi(cmd.Args().Get(4))
		}
		return k.TransitionLightStateCtx(ctx, 1, b, h, s, t, 1000, 1)
	},
}

var bulboff = &cli.Command{
	Name:      "bulboff",
	Usage:     "turn bulb off",
	UsageText: "kasa bulboff host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.TransitionLightStateCtx(ctx, 0, 0, 0, 0, 0, 1000, 1)
	},
}

var ledoff = &cli.Command{
	Name:      "ledoff",
	Usage:     "disable status LED",
	ArgsUsage: "host true|false",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "state"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		b, err := strconv.ParseBool(cmd.StringArg("state"))
		if err != nil {
			return err
		}
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetLEDOffCtx(ctx, b)
	},
}

var reboot = &cli.Command{
	Name:      "reboot",
	Usage:     "reboot device",
	ArgsUsage: "host",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.RebootCtx(ctx)
	},
}

var alias = &cli.Command{
	Name:      "alias",
	Usage:     "update device name (alias)",
	ArgsUsage: "host new-name",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "newname"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)

		nn := cmd.StringArg("newname")
		if nn == "" {
			return fmt.Errorf("need a valid name")
		}

		child := cmd.String("child")
		if child != "" {
			fmt.Fprintf(os.Stderr, "using child %s", child)
			return k.SetChildAliasCtx(ctx, child, nn)
		}

		return k.SetAlias(nn)
	},
}

var raw = &cli.Command{
	Name:      "raw",
	Usage:     "send raw command",
	ArgsUsage: "host command",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "command"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		b, err := k.SendRawCommandCtx(ctx, cmd.StringArg("command"))
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

var diagnose = &cli.Command{
	Name:      "diagnose",
	Usage:     "get device diagnose status",
	UsageText: "kasa diagnose host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		d, err := k.GetDiagnoseStatusCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Uptime: %d\nFree Mem: %d\nRSSI: %d\n%+v\n%+v\n%+v\n", d.Sysinfo.Uptime, d.Sysinfo.FreeMem, d.Wireless.RSSI, d.Cloud, d.Wireless, d.Sysinfo)
		return nil
	},
}

var mcudiagnose = &cli.Command{
	Name:      "mcudiagnose",
	Usage:     "get device mcu diagnose status",
	UsageText: "kasa mcudiagnose host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		d, err := k.GetMCUDiagnoseCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("I2C Total Num: %d\nMCU I2C Reset Num: %d\nMaster I2C Abnormal Num: %d\n", d.I2CTotalNum, d.MCUI2CResetNum, d.MasterI2CAbnormal)
		return nil
	},
}

var pirconfig = &cli.Command{
	Name:      "pirconfig",
	Usage:     "get PIR configuration",
	UsageText: "kasa pirconfig host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		p, err := k.GetPIRConfigCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Enabled: %d, Cold Time: %d, Trigger Index: %d, Min ADC: %d, Max ADC: %d\n", p.Enable, p.ColdTime, p.TriggerIndex, p.MinADC, p.MaxADC)
		return nil
	},
}

var setpirenable = &cli.Command{
	Name:      "setpirenable",
	Usage:     "enable or disable PIR",
	UsageText: "kasa setpirenable host true|false",
	Before:    RequireDevice,
	ArgsUsage: "host state",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "state"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		b, err := strconv.ParseBool(cmd.StringArg("state"))
		if err != nil {
			return err
		}
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetPIREnableCtx(ctx, b)
	},
}

var onboarding = &cli.Command{
	Name:      "onboarding",
	Usage:     "get onboarding status",
	UsageText: "kasa onboarding host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		o, err := k.GetOnboardingStatusCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Onboarding: %s\n", o.Value)
		return nil
	},
}

var sefinfo = &cli.Command{
	Name:      "sefinfo",
	Usage:     "get TP-Link cloud SEF info",
	UsageText: "kasa sefinfo host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		info, err := k.GetSefInfoCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("SEF Server: %s\nDefault Server: %s\nCached Server: %s\n", info.SefServer, info.DefaultServer, info.CachedServer)
		return nil
	},
}

var btncheck = &cli.Command{
	Name:      "btncheck",
	Usage:     "get button check result",
	UsageText: "kasa btncheck host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		b, err := k.GetBtnCheckResCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Reset Button: %v\nSwitch Button: %v\n", b.ResetBtn, b.SwitchBtn)
		return nil
	},
}

var testmode = &cli.Command{
	Name:      "testmode",
	Usage:     "get test mode result",
	UsageText: "kasa testmode host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		t, err := k.GetTestModeCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Factory Mode: %d\n", t.FactoryMode)
		return nil
	},
}
