[![LICENSE](https://img.shields.io/badge/license-BSD-green.svg)](LICENSE)
[![GoReportCard](https://goreportcard.com/badge/cloudkucooland/go-kasa)](https://goreportcard.com/report/cloudkucooland/go-kasa)
[![GoDoc](https://godoc.org/github.com/cloudkucooland/go-kasa?status.svg)](https://godoc.org/github.com/cloudkucooland/go-kasa)


# go-kasa
Go library to control TP-Link Kasa devices.
This library uses the local API, not the cloud API.
It uses UDP rather than TCP where possible for better performance.

## Includes a small cli tool
This is still a work-in-progress, but works for most operations.

## CLI install
Make sure you have Go version 1.18 or newer installed on your system. See [The Go install instructions](https://go.dev/doc/install) for details.

In your shell (terminal on macOS, PowerShell on Windows, lots of options on Linux and UNIX systems...) 
``go install github.com/cloudkucooland/go-kasa/cmd/kasa@latest``

This will place the ``kasa`` binary in ``~/go/bin/kasa`` 
Make sure ``~/go/bin`` is in your shell's path and 

## If you need to control your Kasa devices from Apple HomeKit, I have built a bridge which works well.
[https://github.com/cloudkucooland/HomeKitBrigdges/](https://github.com/cloudkucooland/HomeKitBridges)

## CLI examples
discover devices on the local subnets
```
% kasa discover
read udp [::]:48781: i/o timeout
found 25 devices
  192.168.1.144: HS220(US)             Fireplace Can Dimmer [state: 0] [brightness:  25]
  192.168.1.145: HS200(US)           Dining Room Chandelier [state: 0] [brightness:   0]
  192.168.1.161: HS200(US)            Back Porch Floodlight [state: 1] [brightness:   0]
  192.168.1.162: HS200(US)                   Breakfast Nook [state: 0] [brightness:   0]
  192.168.1.163: HS220(US)        Master Bath Shower Lights [state: 0] [brightness:  50]
  192.168.1.164: HS220(US)           Master Bath Can Lights [state: 0] [brightness:  50]
  192.168.1.165: HS200(US)              Master Bath Mirrors [state: 0] [brightness:   0]
  192.168.1.166: HS200(US)               Front Door Pendant [state: 0] [brightness:   0]
  192.168.1.167: HS200(US)                       Front Room [state: 0] [brightness:   0]
  192.168.1.170: HS210(US)                  Front Hallway 1 [state: 0] [brightness:   0]
  192.168.1.171: HS210(US)                  Front Hallway 2 [state: 0] [brightness:   0]
  192.168.1.172: HS220(US)             Master Bedroom Light [state: 0] [brightness:  40]
  192.168.1.175: HS200(US)               Laundry Room Light [state: 1] [brightness:   0]
  192.168.1.176: HS200(US)       Laundry Room Extractor Fan [state: 0] [brightness:   0]
  192.168.1.177: HS200(US)               Hallway Side Light [state: 1] [brightness:   0]
  192.168.1.178: HS200(US)      Living Room Overhead Lights [state: 0] [brightness:   0]
  192.168.1.179: HS200(US)                  Living Room Fan [state: 0] [brightness:   0]
  192.168.1.180: HS200(US)               Kitchen Sink Light [state: 0] [brightness:   0]
  192.168.1.183: KP303(US) TP-LINK_Power Strip_2BAB
    ID: 8006D442E080440F22A89B072F2E67FB1D9B3DFE02               Guppie Light [state: 0]
    ID: 8006D442E080440F22A89B072F2E67FB1D9B3DFE01              Guppie Heater [state: 1]
    ID: 8006D442E080440F22A89B072F2E67FB1D9B3DFE00              Guppie Filter [state: 1]
  192.168.1.184: KP303(US) TP-LINK_Power Strip_2C77
    ID: 8006972A91D031658289D308866206E11D9B838A02               Edgar Heater [state: 1]
    ID: 8006972A91D031658289D308866206E11D9B838A01              Edgar Bubbler [state: 1]
    ID: 8006972A91D031658289D308866206E11D9B838A00               Edgar Filter [state: 1]
  192.168.1.185: KP303(US) TP-LINK_Power Strip_34EC
    ID: 800661DA15771003D2531C57BE527BA61D9B40E400               Gofish Light [state: 0]
    ID: 800661DA15771003D2531C57BE527BA61D9B40E401              Gofish Heater [state: 1]
    ID: 800661DA15771003D2531C57BE527BA61D9B40E402              Gofish Filter [state: 1]
  192.168.1.187: HS103(US)                Living Room Spare [state: 0] [brightness:   0]
  192.168.1.188: HS103(US)                Scot Bedside Lamp [state: 0] [brightness:   0]
  192.168.1.189: HS103(US)                 Jen Bedside Lamp [state: 0] [brightness:   0]
  192.168.1.193: HS103(US)              Scot’s Office Spare [state: 0] [brightness:   0]
```

disable the cloud service for all devices on the local subnets
```
% kasa nocloud 255.255.255.255
```

toggle one switch
```
% kasa switch 192.168.1.171 false
```

adjust the brightness on a dimmer switch
```
% kasa brightness 192.168.1.164 100
```

show dimmer status and timeings for all dimmer-enabled devices (this needs to be prettier...)
```
% kasa getalldimmer
read udp [::]:58322: i/o timeout
192.168.1.172: &{MinThreshold:0 FadeOnTime:2000 FadeOffTime:2000 GentleOnTime:3000 GentleOffTime:10000 RampRate:30 BulbType:1 ErrCode:0 ErrMsg:}
192.168.1.163: &{MinThreshold:1 FadeOnTime:1000 FadeOffTime:2000 GentleOnTime:3000 GentleOffTime:60000 RampRate:30 BulbType:1 ErrCode:0 ErrMsg:}
192.168.1.144: &{MinThreshold:9 FadeOnTime:1000 FadeOffTime:1000 GentleOnTime:3000 GentleOffTime:10000 RampRate:30 BulbType:0 ErrCode:0 ErrMsg:}
192.168.1.164: &{MinThreshold:0 FadeOnTime:1000 FadeOffTime:2000 GentleOnTime:3000 GentleOffTime:30000 RampRate:30 BulbType:1 ErrCode:0 ErrMsg:}
```

details about a single device
```
% kasa info 192.168.1.144
Alias:		Fireplace Can Dimmer
DevName:	Wi-Fi Smart Dimmer
Model:		HS220(US) [2.0]
Device ID:	xxx
OEM ID:		xxx
Hardware ID:	xxx
Software:	1.0.5 Build 201211 Rel.085320
MIC:		IOT.SMARTPLUGSWITCH
MAC:		60:32:B1:00:00:00
LED Off:	0
Active Mode:	none
Relay:	0	Brightness:	25%
```

Get real-time usage
```
% kasa emeter 192.168.1.203
CurrentMA:	1807
VoltageMV:	122209
PowerMW:	175494
TotalWH:	2097
```

Get daily stats for a month (Feb 2021)
```
% kasa emeter 192.168.1.203 02 2021
2021-02-06 Total WH:	842
2021-02-07 Total WH:	1257
```

Get Countdown Rules (needs to be prettier)
```
% kasa countdown 192.168.1.206
{ID:8725326BB2D0C0DD8D521379163C7D67 Name:TooFar Enable:0 Delay:0 Active:1 Remaining:0}
```

Clear Countdown rules
```
% kasa countdown 192.168.1.206 delete
```
