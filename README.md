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
Make sure ``~/go/bin`` is in your [shell's path](https://janelbrandon.medium.com/understanding-the-path-variable-6eae0936e976#:~:text=PATH%20environment%20variable,referenced%20by%20your%20operating%20environment.)

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
[192.168.100.46] Entrance Chandiler
Min Threshold: 23	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 30000ms	Ramp Rate: 30ms
[192.168.100.44] Scot’s Office Overhead Light
Min Threshold: 11	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 30000ms	Ramp Rate: 30ms
[192.168.100.60] Master Bathroom Closet
Min Threshold: 23	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 10000ms	Ramp Rate: 30ms
[192.168.100.62] Jen’s Office Light
Min Threshold: 23	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 10000ms	Ramp Rate: 30ms
[192.168.100.73] Master Bedroom Light
Min Threshold: 0	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 10000ms	Ramp Rate: 30ms
[192.168.100.24] Fireplace Can Dimmer
Min Threshold: 9	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 10000ms	Ramp Rate: 30ms
[192.168.100.69] Master Bathroom Shower Lights
Min Threshold: 1	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 60000ms	Ramp Rate: 30ms
[192.168.100.51] Master Bathroom Can Lights
Min Threshold: 0	Fade On: 1500ms		Fade Off: 5000ms
Gentle On: 3000ms	Gentle Off: 30000ms	Ramp Rate: 30ms
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

Get daily stats for a month (Feb)
```
% kasa emeter 192.168.1.203 02
2021-02-06 Total WH:	842
2021-02-07 Total WH:	1257
```

Get daily stats for a month (March) on a multi-plug strip
```
 % kasa emeter 192.168.100.27 03     
[1]
2026-03-28:	66Wh
2026-03-29:	66Wh
2026-03-09:	65Wh
2026-03-10:	66Wh
2026-03-11:	66Wh
2026-03-12:	66Wh
2026-03-13:	66Wh
2026-03-14:	66Wh
2026-03-15:	67Wh
2026-03-16:	67Wh
2026-03-17:	68Wh
2026-03-18:	67Wh
2026-03-19:	67Wh
2026-03-20:	66Wh
2026-03-21:	67Wh
2026-03-22:	67Wh
2026-03-23:	65Wh
2026-03-24:	66Wh
2026-03-25:	66Wh
2026-03-26:	64Wh
2026-03-27:	66Wh
2026-03-01:	66Wh
2026-03-02:	66Wh
2026-03-03:	67Wh
2026-03-04:	64Wh
2026-03-05:	66Wh
2026-03-30:	47Wh
	Plug Total:	1766Wh
[2]
2026-03-25:	551Wh
2026-03-26:	1096Wh
2026-03-27:	803Wh
2026-03-28:	920Wh
2026-03-29:	1178Wh
2026-03-01:	1058Wh
2026-03-02:	1046Wh
2026-03-03:	1074Wh
2026-03-04:	572Wh
2026-03-05:	1035Wh
2026-03-06:	1070Wh
2026-03-07:	1107Wh
2026-03-08:	941Wh
2026-03-09:	1108Wh
2026-03-10:	1067Wh
2026-03-11:	1099Wh
2026-03-12:	1121Wh
2026-03-13:	918Wh
2026-03-14:	893Wh
2026-03-15:	1057Wh
2026-03-16:	907Wh
2026-03-17:	923Wh
2026-03-18:	1025Wh
2026-03-19:	1008Wh
2026-03-20:	1055Wh
2026-03-21:	884Wh
2026-03-22:	964Wh
2026-03-23:	927Wh
2026-03-24:	914Wh
2026-03-30:	391Wh
	Plug Total:	28712Wh
[3]
2026-03-21:	263Wh
2026-03-22:	261Wh
2026-03-23:	256Wh
2026-03-24:	258Wh
2026-03-25:	232Wh
2026-03-26:	261Wh
2026-03-27:	253Wh
2026-03-28:	258Wh
2026-03-29:	266Wh
2026-03-01:	269Wh
2026-03-02:	270Wh
2026-03-03:	272Wh
2026-03-04:	233Wh
2026-03-05:	265Wh
2026-03-06:	265Wh
2026-03-07:	251Wh
2026-03-08:	256Wh
2026-03-09:	269Wh
2026-03-10:	268Wh
2026-03-11:	276Wh
2026-03-12:	278Wh
2026-03-13:	259Wh
2026-03-14:	259Wh
2026-03-15:	283Wh
2026-03-16:	262Wh
2026-03-17:	272Wh
2026-03-18:	267Wh
2026-03-19:	268Wh
2026-03-20:	267Wh
2026-03-30:	166Wh
	Plug Total:	7783Wh
[4]
2026-03-07:	233Wh
2026-03-08:	185Wh
2026-03-09:	235Wh
2026-03-10:	250Wh
2026-03-11:	251Wh
2026-03-12:	270Wh
2026-03-13:	196Wh
2026-03-14:	170Wh
2026-03-15:	278Wh
2026-03-16:	201Wh
2026-03-17:	181Wh
2026-03-18:	225Wh
2026-03-19:	210Wh
2026-03-20:	225Wh
2026-03-21:	164Wh
2026-03-22:	195Wh
2026-03-23:	182Wh
2026-03-24:	194Wh
2026-03-25:	41Wh
2026-03-26:	250Wh
2026-03-27:	134Wh
2026-03-28:	174Wh
2026-03-29:	274Wh
2026-03-01:	228Wh
2026-03-02:	243Wh
2026-03-03:	243Wh
2026-03-04:	60Wh
2026-03-05:	216Wh
2026-03-06:	240Wh
2026-03-30:	29Wh
	Plug Total:	5977Wh
[5]
2026-03-15:	6Wh
2026-03-30:	0Wh
	Plug Total:	6Wh
[6]
2026-03-30:	0Wh
	Plug Total:	0Wh
	Strip Total:	44244Wh
```

Emeter status on multi-plug power-strip
 ```
 % kasa emeter 192.168.100.27   
[1]	Current:	39mA	Voltage:	125.12V	Power:	2.83W	Total:	0.15kWh
[2]	Current:	329mA	Voltage:	124.92V	Power:	23.10W	Total:	2.29kWh
[3]	Current:	187mA	Voltage:	125.06V	Power:	9.88W	Total:	0.60kWh
[4]	Current:	32mA	Voltage:	125.08V	Power:	1.72W	Total:	0.46kWh
[5]	Current:	0mA	Voltage:	125.00V	Power:	0.00W	Total:	0.00kWh
[6]	Current:	0mA	Voltage:	125.24V	Power:	0.00W	Total:	0.00kWh
Total	Current:	587mA	Power:	37.52W
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

# Provisioning a new device without the cloud

. Connect to the device's WiFi network
. Set the device name name
```
kasa alias 192.168.0.1 "New Dev Name" 
```
. Turn off cloud
```
kasa nocloud 192.168.0.1
```
. Set the WiFi net
```
kasa setwifi 192.168.0.1 "MySecureSSID" "securenetpw!"
```

===


# If you are researching TP-Link Kasa devices, here are some resources

https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents

https://github.com/whitslack/kasa/blob/master/API.md

http://rat.admin.lv/wp-content/uploads/2018/08/TR17_fgont_-iot_tp_link_hacking.pdf

https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/#TP-Link%20Smart%20Home%20Protocol

https://medium.com/@hu3vjeen/reverse-engineering-tp-link-kc100-bac4641bf1cd

https://machinekoder.com/controlling-tp-link-hs100110-smart-plugs-with-machinekit

https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents

https://github.com/p-doyle/Python-KasaSmartPowerStrip

https://community.hubitat.com/t/release-tp-link-kasa-plug-switch-and-bulb-integration/1675/482

https://www.wiredtron.com/2023/11/28/setting-up-tp-link-kasa-devices-without-smart-home-app.html

[![GoDoc](https://godoc.org/github.com/cloudkucooland/go-kasa?status.svg)](https://godoc.org/github.com/cloudkucooland/go-kasa)
