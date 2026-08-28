# Kasa Device Capability Mapping

This document maps Kasa device models and software versions to supported `kasa` CLI commands.

| Model | Software | Supported Commands |
| :--- | :--- | :--- |
| **All** | Any | `info`, `status`, `emeter`, `wifi`, `gettime`, `gettimezone`, `getschedulerules`, `bulbstate`, `getcountdown`, `ambient`, `diagnose`, `onboarding`, `sefinfo` |
| **HS220(US)** | 1.0.8 | + `dimmer` |
| **HS210(US)** | 1.1.5 | + `btncheck`, `testmode`, `mcudiagnose` |
| **KS200M(US)** | 1.0.12 | + `btncheck`, `testmode`, `mcudiagnose`, `getsensorroutines`, `getmanualaction` |
| **HS200(US)** | 1.1.5 | + `mcudiagnose` |
| **HS300(US)** | 1.1.2 | + `mcudiagnose` |

> **Note:** Capability routing should be implemented using a map keyed by `(Model, SoftwareVersion)`. Commands failing with "Not Supported" should be added to a dynamic denylist for that specific `(Model, SoftwareVersion)`.
