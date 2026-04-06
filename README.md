[![LICENSE](https://img.shields.io/badge/license-BSD-green.svg)](LICENSE)
[![GoReportCard](https://goreportcard.com/badge/cloudkucooland/go-kasa)](https://goreportcard.com/report/cloudkucooland/go-kasa)
[![GoDoc](https://godoc.org/github.com/cloudkucooland/go-kasa?status.svg)](https://godoc.org/github.com/cloudkucooland/go-kasa)


# go-kasa
Go library to control TP-Link Kasa devices.
This library uses the local API, not the cloud API.
It uses UDP rather than TCP where possible for better performance.

## Includes a CLI tool
The CLI is robust and feature-rich, including JSON output. All the commands you need to manage your Kasa environment should be covered.

## CLI install
Make sure you have Go version 1.18 or newer installed on your system. See [The Go install instructions](https://go.dev/doc/install) for details.

In your shell (terminal on macOS, PowerShell on Windows, lots of options on Linux and UNIX systems...) 
``go install github.com/cloudkucooland/go-kasa/cmd/kasa@latest``

This will place the ``kasa`` binary in ``~/go/bin/kasa`` 
Make sure ``~/go/bin`` is in your [shell's path](https://janelbrandon.medium.com/understanding-the-path-variable-6eae0936e976#:~:text=PATH%20environment%20variable,referenced%20by%20your%20operating%20environment.)

## If you need to control your Kasa devices from Apple HomeKit, I have built a bridge which works well.
[https://github.com/cloudkucooland/HomeKitBrigdges/](https://github.com/cloudkucooland/HomeKitBridges)

## Monitor your electricity usage
If you need to monitor your Kasa devices that support emetering, this includes a daemon which feeds the emeter data into InfluxDB. A sample Grafana setup is included, all running in docker.

## CLI examples
discover devices on the local subnets
```
% kasa discover
Device                               IP/ID:                                     Model     State Brightness
Hallway Side Light                   192.168.100.20                             HS200(US) Off      0
Guest Bathroom Light                 192.168.100.21                             HS200(US) Off      0
Living Room Red Lamp                 192.168.100.22                             HS103(US) Off      0
Outside Front Carriage Lights        192.168.100.23                             HS200(US) On       0
Fireplace Can Dimmer                 192.168.100.24                             HS220(US) On      75
Laundry Room Extractor Fan           192.168.100.25                             HS200(US) Off      0
Conservatory Chandelier              192.168.100.26                             HS200(US) Off      0
Living Room Overhead Lights          192.168.100.27                             HS200(US) Off      0
Tea Kettle                           192.168.100.28                             KP115(US) On       0
Jen’s Bedside Lamp                   192.168.100.29                             HS103(US) Off      0
Garage                               192.168.100.30                             HS300(US)       
Garage/Battery charger 1             8006F0636CA2DCC4AE3622D483F75865224A78C800           On    
Garage/Battery charger 2             8006F0636CA2DCC4AE3622D483F75865224A78C801           On    
Garage/Battery charger 3             8006F0636CA2DCC4AE3622D483F75865224A78C802           On    
Garage/Unused                        8006F0636CA2DCC4AE3622D483F75865224A78C805           On    
Garage/Hot plate                     8006F0636CA2DCC4AE3622D483F75865224A78C803           On    
Garage/Condenser pump                8006F0636CA2DCC4AE3622D483F75865224A78C804           On    
Kitchen Sink Light                   192.168.100.31                             HS200(US) On       0
Master Bathroom Shower Lights        192.168.100.32                             HS220(US) On      40
Living Room Fan                      192.168.100.33                             HS200(US) Off      0
Master Bathroom Closet               192.168.100.34                             HS220(US) On      75
Go Fish                              192.168.100.35                             KP303(US)       
Go Fish/Fish tank 3 filter           800661DA15771003D2531C57BE527BA61D9B40E401           Off   
Go Fish/Fish tank 3 light            800661DA15771003D2531C57BE527BA61D9B40E402           On    
Go Fish/Fish tank 3 heater           800661DA15771003D2531C57BE527BA61D9B40E400           On    
Back Porch Floodlight                192.168.100.36                             HS200(US) On       0
Master Bathroom Can Lights           192.168.100.37                             HS220(US) Off     20
Jen’s Office Light                   192.168.100.38                             HS220(US) Off     17
Potting Bench                        192.168.100.39                             HS300(US)       
Potting Bench/Fan                    800663CD078A8C16B66658BEAA8D3F981EEF3EEA00           Off   
Potting Bench/Bench light 1          800663CD078A8C16B66658BEAA8D3F981EEF3EEA01           Off   
Potting Bench/Bench light 2          800663CD078A8C16B66658BEAA8D3F981EEF3EEA02           Off   
Potting Bench/Heat mat               800663CD078A8C16B66658BEAA8D3F981EEF3EEA03           Off   
Potting Bench/Sewing machine         800663CD078A8C16B66658BEAA8D3F981EEF3EEA04           Off   
Potting Bench/Unused                 800663CD078A8C16B66658BEAA8D3F981EEF3EEA05           Off   
Front Hallway 1                      192.168.100.40                             HS210(US) Off      0
Frog Tank                            192.168.100.41                             HS300(US)       
Frog Tank/Unused                     800609CCE1FB71B08497FF57B088A1132225FA3804           Off   
Frog Tank/Unused                     800609CCE1FB71B08497FF57B088A1132225FA3805           Off   
Frog Tank/UV                         800609CCE1FB71B08497FF57B088A1132225FA3800           Off   
Frog Tank/Incandescent               800609CCE1FB71B08497FF57B088A1132225FA3801           Off   
Frog Tank/Grow light                 800609CCE1FB71B08497FF57B088A1132225FA3802           On    
Frog Tank/Heater                     800609CCE1FB71B08497FF57B088A1132225FA3803           On    
Laundry Room Light                   192.168.100.42                             HS200(US) On       0
Master Bedroom Light                 192.168.100.43                             HS220(US) Off    100
Front Room                           192.168.100.44                             HS200(US) Off      0
Outside Front Overhead Light         192.168.100.45                             HS200(US) On       0
Hallway Can Lights                   192.168.100.46                             HS210(US) On       0
Jen’s Office Fan                     192.168.100.48                             HS200(US) Off      0
Master Bathroom Fan                  192.168.100.49                             HS200(US) Off      0
Master Bathroom Mirrors              192.168.100.50                             HS200(US) Off      0
Kitchen Island Light                 192.168.100.51                             HS200(US) Off      0
Jen’s Reading Lamp                   192.168.100.52                             HS103(US) Off      0
Breakfast Nook                       192.168.100.53                             HS200(US) Off      0
Scot’s Office Overhead Light         192.168.100.54                             HS220(US) Off     46
Guppie                               192.168.100.55                             KP303(US)       
Guppie/Fish tank 1 filter            8006D442E080440F22A89B072F2E67FB1D9B3DFE02           On    
Guppie/Fish tank 1 light             8006D442E080440F22A89B072F2E67FB1D9B3DFE01           Off   
Guppie/Fish tank 1 heater            8006D442E080440F22A89B072F2E67FB1D9B3DFE00           On    
Edgar                                192.168.100.57                             KP303(US)       
Edgar/Fish tank 2 light              8006972A91D031658289D308866206E11D9B838A02           On    
Edgar/Fish tank 2 filter             8006972A91D031658289D308866206E11D9B838A01           On    
Edgar/Fish tank 2 heater             8006972A91D031658289D308866206E11D9B838A00           Off   
Scot’s Bedside Lamp                  192.168.100.58                             HS103(US) Off      0
Front Hallway 2                      192.168.100.60                             HS210(US) Off      0
Scot’s Office Fan                    192.168.100.61                             HS200(US) Off      0
Plant Shelf 1                        192.168.100.62                             HS300(US)       
Plant Shelf 1/Primary                800621D135774DF39C8C9A878A5ACF9124ACADE500           Off   
Plant Shelf 1/Secondary              800621D135774DF39C8C9A878A5ACF9124ACADE501           Off   
Plant Shelf 1/Spotlight              800621D135774DF39C8C9A878A5ACF9124ACADE502           Off   
Plant Shelf 1/Plug 4                 800621D135774DF39C8C9A878A5ACF9124ACADE503           Off   
Plant Shelf 1/Plug 5                 800621D135774DF39C8C9A878A5ACF9124ACADE504           Off   
Plant Shelf 1/Plug 6                 800621D135774DF39C8C9A878A5ACF9124ACADE505           Off   
Plant Shelf 2                        192.168.100.63                             HS300(US)       
Plant Shelf 2/Primary                8006AA385A1A397D4C1A9AF6BC6D8D8324CB87F200           Off   
Plant Shelf 2/Secondary              8006AA385A1A397D4C1A9AF6BC6D8D8324CB87F201           Off   
Plant Shelf 2/Spotlight              8006AA385A1A397D4C1A9AF6BC6D8D8324CB87F202           Off   
Plant Shelf 2/Plug 4                 8006AA385A1A397D4C1A9AF6BC6D8D8324CB87F203           Off   
Plant Shelf 2/Plug 5                 8006AA385A1A397D4C1A9AF6BC6D8D8324CB87F204           Off   
Plant Shelf 2/Plug 6                 8006AA385A1A397D4C1A9AF6BC6D8D8324CB87F205           Off   
Conservatory AV                      192.168.100.64                             HS300(US)       
Conservatory AV/Peachtree Amp        80068B68F17684F2CB930C6CDC47FC3F2375E7FB00           On    
Conservatory AV/Clearaudio Turntable 80068B68F17684F2CB930C6CDC47FC3F2375E7FB01           On    
Conservatory AV/Stream Mac           80068B68F17684F2CB930C6CDC47FC3F2375E7FB02           On    
Conservatory AV/Air Filter           80068B68F17684F2CB930C6CDC47FC3F2375E7FB03           On    
Conservatory AV/Printer              80068B68F17684F2CB930C6CDC47FC3F2375E7FB04           On    
Conservatory AV/Monitor              80068B68F17684F2CB930C6CDC47FC3F2375E7FB05           On    
Counter Fish Tank                    192.168.100.65                             HS300(US)       
Counter Fish Tank/Pi Zero            8006180DF51AF68BAB85AD990E3BD0E023760CFC00           On    
Counter Fish Tank/Heater             8006180DF51AF68BAB85AD990E3BD0E023760CFC01           Off   
Counter Fish Tank/Light              8006180DF51AF68BAB85AD990E3BD0E023760CFC02           On    
Counter Fish Tank/Plug 4             8006180DF51AF68BAB85AD990E3BD0E023760CFC03           Off   
Counter Fish Tank/Filter             8006180DF51AF68BAB85AD990E3BD0E023760CFC04           On    
Counter Fish Tank/Plug 6             8006180DF51AF68BAB85AD990E3BD0E023760CFC05           Off  
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
 % kasa dimmer
Device                        IP             Min Fade On Fade Off Gentle On Gentle Off Ramp Rate
Scot’s Office Overhead Light  192.168.100.54 11  1500    5000     3000      30000      30
Fireplace Can Dimmer          192.168.100.24 9   1500    5000     3000      10000      30
Master Bedroom Light          192.168.100.43 0   1500    5000     3000      10000      30
Master Bathroom Closet        192.168.100.34 23  1500    5000     3000      10000      30
Master Bathroom Can Lights    192.168.100.37 0   1500    5000     3000      30000      30
Master Bathroom Shower Lights 192.168.100.32 1   1500    5000     3000      60000      30
Jen’s Office Light            192.168.100.38 23  1500    5000     3000      10000      30
Entrance Chandiler            192.168.100.56 23  1500    5000     3000      30000      30
```

details about a single device
```
 % kasa info 192.168.100.56
Alias:       Entrance Chandiler
DevName:     Wi-Fi Smart Dimmer
Model:       HS220(US) (2.0)
Device ID:   800660D983102C81B1DBC3F890B96FDA1E35996A
OEM ID:      97D03CA037C71B6BCFAB8705E9B0C417
Hardware ID: CA321A94521D18706FC7C34CA84F99A4
Software:    1.0.8 Build 210423 Rel.075507
MIC:         IOT.SMARTPLUGSWITCH
MAC:         28:EE:52:AA:9C:31
LED Off:     0
Active Mode: none
Outlet       Relay State Brightness
             1           25
```

Get real-time usage
```
 % kasa emeter
Device                                 Current Voltage Power   Since Reset
Conservatory AV/Peachtree Amp          0mA     122.49V 0.00W   0.00kWh
Conservatory AV/Clearaudio Turntable   0mA     122.73V 0.00W   0.00kWh
Conservatory AV/Stream Mac             74mA    122.86V 4.50W   0.04kWh
Conservatory AV/Air Filter             163mA   122.87V 18.48W  0.15kWh
Conservatory AV/Printer                0mA     123.53V 0.00W   0.01kWh
Conservatory AV/Monitor                24mA    123.00V 1.07W   0.01kWh
Living Room Entertainment/Sony Blu-ray 0mA     123.07V 0.00W   0.00kWh
Living Room Entertainment/Sony TV      1121mA  122.86V 134.89W 0.74kWh
Living Room Entertainment/Subwoofer    0mA     122.92V 0.00W   0.00kWh
Living Room Entertainment/Onkyo Amp    462mA   122.99V 38.21W  0.16kWh
Living Room Entertainment/Turntable    0mA     122.83V 0.00W   0.00kWh
Living Room Entertainment/CD player    0mA     123.08V 0.00W   0.00kWh
Counter Fish Tank/Pi Zero              16mA    123.06V 0.82W   0.01kWh
Counter Fish Tank/Heater               0mA     123.00V 0.00W   0.09kWh
Counter Fish Tank/Light                108mA   123.08V 6.96W   0.04kWh
Counter Fish Tank/Plug 4               0mA     122.97V 0.00W   0.00kWh
Counter Fish Tank/Filter               148mA   122.98V 7.13W   0.06kWh
Counter Fish Tank/Plug 6               0mA     123.03V 0.00W   0.00kWh
Potting Bench/Fan                      0mA     122.59V 0.00W   0.00kWh
Potting Bench/Bench light 1            0mA     122.58V 0.00W   0.00kWh
Potting Bench/Bench light 2            0mA     122.69V 0.00W   0.00kWh
Potting Bench/Heat mat                 0mA     122.71V 0.00W   0.00kWh
Potting Bench/Sewing machine           0mA     122.88V 0.00W   0.46kWh
Potting Bench/Unused                   0mA     122.84V 0.00W   0.00kWh
Tea Kettle                             0mA     122.57V 0.00W   0.12kWh
Frog Tank/Unused                       0mA     122.68V 0.00W   0.09kWh
Frog Tank/Unused                       0mA     123.96V 0.00W   0.07kWh
Frog Tank/UV                           0mA     123.33V 0.00W   0.11kWh
Frog Tank/Incandescent                 0mA     123.17V 0.00W   0.24kWh
Frog Tank/Grow light                   0mA     123.35V 0.00W   0.00kWh
Frog Tank/Heater                       0mA     123.02V 0.00W   0.00kWh
Plant Shelf 2/Primary                  0mA     123.34V 0.00W   0.01kWh
Plant Shelf 2/Secondary                0mA     123.40V 0.00W   0.40kWh
Plant Shelf 2/Spotlight                0mA     123.81V 0.00W   0.06kWh
Plant Shelf 2/Plug 4                   0mA     123.39V 0.00W   0.00kWh
Plant Shelf 2/Plug 5                   0mA     123.10V 0.00W   0.00kWh
Plant Shelf 2/Plug 6                   0mA     123.44V 0.00W   0.00kWh
Garage/Battery charger 1               13mA    122.65V 0.70W   0.01kWh
Garage/Battery charger 2               0mA     123.44V 0.00W   0.00kWh
Garage/Battery charger 3               0mA     123.63V 0.00W   0.00kWh
Garage/Unused                          0mA     123.59V 0.00W   0.00kWh
Garage/Hot plate                       0mA     123.07V 0.00W   0.00kWh
Garage/Condenser pump                  0mA     123.51V 0.00W   0.00kWh
Plant Shelf 1/Primary                  0mA     123.40V 0.00W   0.40kWh
Plant Shelf 1/Secondary                0mA     123.47V 0.00W   0.00kWh
Plant Shelf 1/Spotlight                0mA     123.60V 0.00W   0.06kWh
Plant Shelf 1/Plug 4                   0mA     122.44V 0.00W   0.00kWh
Plant Shelf 1/Plug 5                   0mA     123.51V 0.00W   0.00kWh
Plant Shelf 1/Plug 6                   0mA     123.44V 0.00W   0.00kWh
Total House                            2129mA          212.76W 3.32kWh
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

Run the json output through jq
```
 % kasa -j emeter | jq '.[] | {alias: .alias, active: [.Realtime[] | select(.power_mw > 0) | {outlet: .alias, watts: (.power_mw/1000)}]}}
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
