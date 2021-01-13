# go-kasa
Go library to TP-Link Kasa devices

## Includes a small cli tool
This is still a work-in-progress, but works for what I need.

## examples
```
scot@covert:/home/scot/go-kasa % kasa discover
read udp [::]:40569: i/o timeout
found 11 devices
192.168.1.166: HS200(US) [state: 1] [brightness:   0] Front Door Pendant
192.168.1.164: HS220(US) [state: 1] [brightness:  28] Master Bath Can Lights
192.168.1.172: HS220(US) [state: 0] [brightness:  66] Master Bedroom Light
192.168.1.165: HS200(US) [state: 0] [brightness:   0] Master Bath Mirrors
192.168.1.162: HS200(US) [state: 1] [brightness:   0] Breakfast Nook
192.168.1.144: HS220(US) [state: 1] [brightness: 100] Fireplace Can Dimmer
192.168.1.163: HS220(US) [state: 1] [brightness:  39] Master Bath Shower Lights
192.168.1.161: HS200(US) [state: 1] [brightness:   0] Back Porch Floodlight
192.168.1.167: HS200(US) [state: 1] [brightness:   0] Front Room
192.168.1.171: HS210(US) [state: 1] [brightness:   0] Front Hallway 2
192.168.1.170: HS210(US) [state: 1] [brightness:   0] Front Hallway 1
scot@covert:/home/scot/go-kasa % kasa nocloud 255.255.255.255
scot@covert:/home/scot/go-kasa % kasa switch 192.168.1.171 false
scot@covert:/home/scot/go-kasa % kasa brightness 192.168.1.164 100
scot@covert:/home/scot/go-kasa % kasa discover
read udp [::]:50027: i/o timeout
found 11 devices
192.168.1.172: HS220(US) [state: 0] [brightness:  66] Master Bedroom Light
192.168.1.164: HS220(US) [state: 1] [brightness: 100] Master Bath Can Lights
192.168.1.170: HS210(US) [state: 0] [brightness:   0] Front Hallway 1
192.168.1.165: HS200(US) [state: 0] [brightness:   0] Master Bath Mirrors
192.168.1.144: HS220(US) [state: 1] [brightness: 100] Fireplace Can Dimmer
192.168.1.161: HS200(US) [state: 1] [brightness:   0] Back Porch Floodlight
192.168.1.163: HS220(US) [state: 1] [brightness:  39] Master Bath Shower Lights
192.168.1.167: HS200(US) [state: 1] [brightness:   0] Front Room
192.168.1.171: HS210(US) [state: 0] [brightness:   0] Front Hallway 2
192.168.1.166: HS200(US) [state: 1] [brightness:   0] Front Door Pendant
192.168.1.162: HS200(US) [state: 1] [brightness:   0] Breakfast Nook
```
