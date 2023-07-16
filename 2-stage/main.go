package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

const ( // cmd
	WHOISHERE = 0x01
	IAMHERE   = 0x02
	GETSTATUS = 0x03
	STATUS    = 0x04
	SETSTATUS = 0x05
	TICK      = 0x06
)

const ( // devType
	SMARTHUB  = 0x01
	ENVSENSOR = 0x02
	SWITCH    = 0x03
	LAMP      = 0x04
	SOCKET    = 0x05
	CLOCK     = 0x06
)

const (
	BROADCAST = 0x3fff
)

const ( // sensor's bit mask
	TEMPSENS     = 0x01
	HUMIDITYSENS = 0x02
	LIGHTSENS    = 0x04
	AIRPOLLSENS  = 0x08
)

const ( // sensor's bit mask
	TEMPSENSID     = 0x00
	HUMIDITYSENSID = 0x01
	LIGHTSENSID    = 0x02
	AIRPOLLSENSID  = 0x03
)

type uleb128 struct {
	first  uint64
	second uint64
}

type switcher struct {
	id       uint16
	name     string
	devNames []string
	status   byte
}

type device struct {
	id      uint16
	name    string
	devType byte
	status  byte
}

type clock struct {
	name      string
	timestamp uleb128
}

type envSensor struct {
	name     string
	id       uint16
	sensors  byte
	triggers map[string]trigger
}

type trigger struct {
	op    byte
	value uleb128
	name  string
}

type devProps struct {
	envSensor
	switcher
}

type cmdBody struct {
	timestamp uleb128 // Clock
	devName   string
	devProps  devProps
	values    []uleb128 // EnvSensor
	status    byte      // Switch/Lamp/Socket
}

type payload struct {
	src     uleb128 //отправитель
	dst     uleb128 //получатель
	serial  uleb128
	devType byte
	cmd     byte
	cmdBody cmdBody
}

type packet struct {
	length  byte
	payload payload
	crc8    byte
}

func main() {
	serverUrl := os.Args[1]
	hubId := uleb128{0, 0}
	var err error
	hubId.first, err = hexStringToUint(os.Args[2])
	if err != nil {
		//panic(err)
	}

	var hubSerial uleb128
	var swithersDB = make(map[uint16]switcher, 0)
	var envSensorsDB = make(map[uint16]envSensor, 0)
	var devsDB = make(map[string]device, 0)
	var clock1 clock

	pkgsToSent := make([]packet, 0)
	bufout := bytes.NewBuffer([]byte{})

	pkg := getStartPkg(hubId)
	pkgsToSent = append(pkgsToSent, pkg)

	for {
		binPackets := make([]byte, 0)
		for len(pkgsToSent) > 0 {
			pkg = pkgsToSent[len(pkgsToSent)-1]
			pkgsToSent = pkgsToSent[:len(pkgsToSent)-1]
			binPkg := getBinPkg(pkg)
			binPackets = append(binPackets, binPkg...)
		}

		base64UrlEncodeData := base64.RawURLEncoding.EncodeToString(binPackets)
		bufout.WriteString(base64UrlEncodeData)

		resp, err := http.Post(serverUrl, "application/octet-stream", bufout)
		if err != nil {
			//log.Println(err)
			os.Exit(99)
		}
		respStatusCode := resp.Status[:3]

		if respStatusCode != "200" && respStatusCode != "204" {
			//log.Println(resp.Status)
			os.Exit(99)
		}
		if respStatusCode == "204" {
			os.Exit(0)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println(err)
			continue
		}
		respStr := string(respBody)

		decodeData, err := base64.RawURLEncoding.DecodeString(respStr)
		if err != nil {
			//log.Println(err)
			continue
		}

		respPkgs := getPkgsFromBytes(decodeData)
		pkgsToSent = processingPackets(respPkgs, hubId, &hubSerial, &envSensorsDB, &swithersDB, &devsDB, &clock1)
	}
}

func hexStringToUint(hexStr string) (uint64, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	var result uint64

	var pow16 uint64 = 1
	for i := len(hexStr) - 1; i >= 0; i-- {
		if (hexStr[i] >= '0') && (hexStr[i] <= '9') {
			result += pow16 * uint64(hexStr[i]-'0')
			pow16 *= 16
		} else if (hexStr[i] >= 'a') && (hexStr[i] <= 'f') {
			result += pow16 * uint64(hexStr[i]-'a'+10)
			pow16 *= 16
		} else {
			return result, errors.New("invalid hex-symbol")
		}
	}
	return result, nil
}

func EncodeUleb128(val uint64) []byte {
	var result []byte
	for {
		b := byte(val & 0x7F)
		val >>= 7
		if val != 0 {
			b |= 0x80
		}
		result = append(result, b)
		if val == 0 {
			break
		}
	}
	return result
}

func DecodeUleb128(data []byte) uint64 {
	var result uint64
	var shift uint
	for _, b := range data {
		result |= uint64(b&0x7F) << shift
		shift += 7
		if b&0x80 == 0 {
			break
		}
	}
	return result
}

func getBytesFromUleb(u uleb128) []byte {
	buf := make([]byte, 0, 16)
	var encodedFirst []byte
	var encodedSecond []byte
	if u.second == 0 {
		encodedFirst = EncodeUleb128(u.first)
		buf = append(buf, encodedFirst...)
		return buf
	}
	encodedFirst = EncodeUleb128(u.first)
	encodedSecond = EncodeUleb128(u.second)
	buf = append(buf, encodedFirst...)
	buf = append(buf, encodedSecond...)
	return buf
}

func getUlebFromBytes(bytes []byte) uleb128 { //bytes' len <= 16
	var u uleb128
	if len(bytes) <= 8 {
		u.first = DecodeUleb128(bytes)
		return u
	}
	u.first = DecodeUleb128(bytes[:8])
	u.second = DecodeUleb128(bytes[8:])
	return u
}

func ulebPlus(u *uleb128, val uint64) *uleb128 {
	if u.first+val > math.MaxUint64 {
		u.second += (u.first + val) - math.MaxUint64
		u.first = math.MaxUint64
	} else {
		u.first += val
	}
	return u
}

func getBinPkg(pkg packet) (buf []byte) {
	payloadSize := pkg.length
	buf = make([]byte, 0, 2+payloadSize)

	buf = append(buf, pkg.length)

	buf = append(buf, getBinPayload(&pkg.payload)...)

	buf = append(buf, pkg.crc8)

	return buf
}

func getBinPayload(payload1 *payload) []byte {
	buf := make([]byte, 0, 5)
	buf = append(buf, getBytesFromUleb(payload1.src)...)
	buf = append(buf, getBytesFromUleb(payload1.dst)...)
	buf = append(buf, getBytesFromUleb(payload1.serial)...)
	buf = append(buf, payload1.devType)
	buf = append(buf, payload1.cmd)
	buf = append(buf, getBinCmdBody(&payload1.cmdBody, payload1.cmd, payload1.devType)...)

	return buf
}

func getBinCmdBody(cbody *cmdBody, cmd uint8, devType uint8) []byte {
	buf := make([]byte, 0)

	if cmd == TICK && devType == CLOCK {
		buf = append(buf, getBytesFromUleb(cbody.timestamp)...)
	} else if cmd == GETSTATUS { //cmdBody is empty
		return buf
	} else if cmd == WHOISHERE || cmd == IAMHERE {
		nameLen := len(cbody.devName)
		buf = append(buf, byte(nameLen))
		buf = append(buf, cbody.devName...)
		if devType == ENVSENSOR {
			buf = append(buf, cbody.devProps.sensors)
			buf = append(buf, byte(len(cbody.devProps.triggers)))
			for i, _ := range cbody.devProps.triggers {
				buf = append(buf, cbody.devProps.triggers[i].op)
				buf = append(buf, getBytesFromUleb(cbody.devProps.triggers[i].value)...)
				nameLen = len(cbody.devProps.triggers[i].name)
				buf = append(buf, byte(nameLen))
				buf = append(buf, cbody.devProps.triggers[i].name...)
			}
		} else if devType == SWITCH {
			buf = append(buf, byte(len(cbody.devProps.devNames)))
			for i, _ := range cbody.devProps.devNames {
				nameLen = len(cbody.devProps.devNames[i])
				buf = append(buf, byte(nameLen))
				buf = append(buf, cbody.devProps.devNames[i]...)
			}
		}
	} else if cmd == STATUS {
		if devType == ENVSENSOR {
			buf = append(buf, byte(len(cbody.values)))
			for i, _ := range cbody.values {
				buf = append(buf, getBytesFromUleb(cbody.values[i])...)
			}
		} else if devType == SWITCH || devType == LAMP || devType == SOCKET {
			buf = append(buf, cbody.status)
		}
	} else if cmd == SETSTATUS {
		buf = append(buf, cbody.status)
	}

	return buf
}

func getPkgsFromBytes(buf []byte) []packet {
	pkgs := make([]packet, 0, 1)

	for len(buf) > 0 {
		pkg := packet{}
		pkg.length = buf[0]
		buf = buf[1:]

		pkg.payload = getPayloadFromBytes(buf[0:pkg.length])
		buf = buf[pkg.length:]
		pkg.crc8 = buf[0]
		buf = buf[1:]
		pkgs = append(pkgs, pkg)
	}

	return pkgs
}

func getPayloadFromBytes(buf []byte) payload {
	payld := payload{}
	if buf[0] >= 128 {
		payld.src = getUlebFromBytes(buf[0:2])
		buf = buf[2:]
	} else {
		payld.src = getUlebFromBytes(buf[0:1])
		buf = buf[1:]
	}
	if buf[0] >= 128 {
		payld.dst = getUlebFromBytes(buf[0:2])
		buf = buf[2:]
	} else {
		payld.dst = getUlebFromBytes(buf[0:1])
		buf = buf[1:]
	}
	ulebBytes := make([]byte, 0)
	for buf[0] >= 128 {
		ulebBytes = append(ulebBytes, buf[0])
		buf = buf[1:]
	}
	ulebBytes = append(ulebBytes, buf[0])
	buf = buf[1:]
	payld.serial = getUlebFromBytes(ulebBytes)

	payld.devType = buf[0]
	buf = buf[1:]
	payld.cmd = buf[0]
	buf = buf[1:]
	payld.cmdBody = getCmdBodyFromBytes(buf, payld.cmd, payld.devType)

	return payld
}

func getCmdBodyFromBytes(buf []byte, cmd byte, devType byte) cmdBody {
	cbody := cmdBody{}
	if cmd == TICK && devType == CLOCK {
		ulebBytes := make([]byte, 0)
		for buf[0] >= 128 {
			ulebBytes = append(ulebBytes, buf[0])
			buf = buf[1:]
		}
		ulebBytes = append(ulebBytes, buf[0])
		buf = buf[1:]
		cbody.timestamp = getUlebFromBytes(ulebBytes)
	} else if cmd == GETSTATUS {
		return cbody
	} else if cmd == WHOISHERE || cmd == IAMHERE {
		nameLen := buf[0]
		buf = buf[1:]
		devName := make([]byte, nameLen)
		for i := 0; i < int(nameLen); i++ {
			devName[i] = buf[i]
		}
		buf = buf[nameLen:]
		cbody.devName = string(devName)
		if devType == ENVSENSOR {
			cbody.devProps.sensors = buf[0]
			buf = buf[1:]
			triggersLen := buf[0]
			buf = buf[1:]
			cbody.devProps.triggers = make(map[string]trigger)
			for i := 0; i < int(triggersLen); i++ {
				trig := trigger{}
				trig.op = buf[0]
				buf = buf[1:]
				ulebBytes := make([]byte, 0)
				for buf[0] >= 128 {
					ulebBytes = append(ulebBytes, buf[0])
					buf = buf[1:]
				}
				ulebBytes = append(ulebBytes, buf[0])
				buf = buf[1:]
				trig.value = getUlebFromBytes(ulebBytes)
				nameLen := buf[0]
				buf = buf[1:]
				devName := make([]byte, nameLen)
				for i := 0; i < int(nameLen); i++ {
					devName[i] = buf[i]
				}
				buf = buf[nameLen:]
				trig.name = string(devName)
				cbody.devProps.triggers[trig.name] = trig
			}
		} else if devType == SWITCH {
			devNamesLen := buf[0]
			buf = buf[1:]
			for i := 0; i < int(devNamesLen); i++ {
				nameLen := buf[0]
				buf = buf[1:]
				devName := make([]byte, nameLen)
				for i := 0; i < int(nameLen); i++ {
					devName[i] = buf[i]
				}
				buf = buf[nameLen:]
				cbody.devProps.devNames = append(cbody.devProps.devNames, string(devName))
			}
		}
	} else if cmd == STATUS {
		if devType == ENVSENSOR {
			valsLen := buf[0]
			buf = buf[1:]
			for i := 0; i < int(valsLen); i++ {
				ulebBytes := make([]byte, 0)
				for buf[0] >= 128 {
					ulebBytes = append(ulebBytes, buf[0])
					buf = buf[1:]
				}
				ulebBytes = append(ulebBytes, buf[0])
				buf = buf[1:]
				cbody.values = append(cbody.values, getUlebFromBytes(ulebBytes))
			}
		} else if devType == SWITCH || devType == LAMP || devType == SOCKET {
			cbody.status = buf[0]
			buf = buf[1:]
		}
	}
	return cbody
}

func processingPackets(respPkgs []packet, hubId uleb128, hubSerialP *uleb128, envSensorsP *map[uint16]envSensor, swithcersP *map[uint16]switcher, devsP *map[string]device, clock *clock) []packet {
	envSensors := *envSensorsP
	swithcers := *swithcersP
	devs := *devsP

	pkgsToSent := make([]packet, 0)
	for i := 0; i < len(respPkgs); i++ {
		respPkg := &respPkgs[i]
		if respPkg.crc8 != calcCRC8(getBinPayload(&respPkg.payload)) {
			continue
		}
		if respPkg.payload.cmd == TICK {
			clock.timestamp = respPkg.payload.cmdBody.timestamp
		} else if respPkg.payload.cmd == IAMHERE || respPkg.payload.cmd == WHOISHERE {
			if respPkg.payload.devType == SWITCH {
				swtch := switcher{}
				swtch.id = uint16(respPkg.payload.src.first)
				swtch.name = respPkg.payload.cmdBody.devName
				swtch.devNames = append(swtch.devNames, respPkg.payload.cmdBody.devProps.devNames...)
				swithcers[swtch.id] = swtch
			} else if respPkg.payload.devType == LAMP || respPkg.payload.devType == SOCKET {
				dev := device{}
				dev.id = uint16(respPkg.payload.src.first)
				dev.name = respPkg.payload.cmdBody.devName
				dev.devType = respPkg.payload.devType
				devs[dev.name] = dev
			} else if respPkg.payload.devType == CLOCK {
				clock.name = respPkg.payload.cmdBody.devName
			} else if respPkg.payload.devType == ENVSENSOR {
				envSens := envSensor{}
				envSens.id = uint16(respPkg.payload.src.first)
				envSens.name = respPkg.payload.cmdBody.devName
				envSens.sensors = respPkg.payload.cmdBody.devProps.sensors
				envSens.triggers = make(map[string]trigger)
				for i, v := range respPkg.payload.cmdBody.devProps.triggers {
					envSens.triggers[i] = v
				}
				envSensors[envSens.id] = envSens
			}

			if respPkg.payload.cmd == WHOISHERE {
				hubSerialP = ulebPlus(hubSerialP, 1)
				npayld := payload{
					src:     hubId,
					dst:     uleb128{16383, 0},
					serial:  *hubSerialP,
					devType: SMARTHUB,
					cmd:     IAMHERE,
					cmdBody: cmdBody{
						devName: "HUB01",
					},
				}
				binPayld := getBinPayload(&npayld)
				npkg := packet{
					length:  byte(len(binPayld)),
					payload: npayld,
					crc8:    calcCRC8(binPayld),
				}

				pkgsToSent = append(pkgsToSent, npkg)
			}

			// create packet with GETSTATUS cmd
			hubSerialP = ulebPlus(hubSerialP, 1)
			var npayload payload
			initPayload(&npayload, hubId, respPkg.payload.src, *hubSerialP, respPkg.payload.devType, GETSTATUS, nil)
			var npkg packet
			initPacket(&npkg, &npayload)

			pkgsToSent = append(pkgsToSent, npkg)
		} else if respPkg.payload.cmd == STATUS {
			if respPkg.payload.devType == SWITCH {
				swtch := swithcers[uint16(respPkg.payload.src.first)]
				swtch.status = respPkg.payload.cmdBody.status
				swithcers[swtch.id] = swtch

				for _, v := range swtch.devNames {
					if devs[v].status != swtch.status { // create packet with cmd SETSTATUS
						hubSerialP = ulebPlus(hubSerialP, 1)
						var cmdbody = cmdBody{
							status: swtch.status,
						}
						dev := devs[v]
						var npayload payload
						initPayload(&npayload, hubId, uleb128{uint64(dev.id), 0}, *hubSerialP, dev.devType, SETSTATUS, &cmdbody)
						var npkg packet
						initPacket(&npkg, &npayload)

						pkgsToSent = append(pkgsToSent, npkg)
					}
				}
			} else if respPkg.payload.devType == LAMP || respPkg.payload.devType == SOCKET {
				var dev device
				ok := false
				for _, v := range devs {
					if uint16(respPkg.payload.src.first) == v.id {
						dev = v
						ok = true
					}
				}
				if !ok {
					//log.Println("device with id ", respPkg.payload.src.first, "not exist in database")
					continue
				}
				dev.status = respPkg.payload.cmdBody.status
				devs[dev.name] = dev
			} else if respPkg.payload.devType == ENVSENSOR {
				envSens := envSensors[uint16(respPkg.payload.src.first)]
				envSensProp := make([]uint64, 4)
				i := 0
				if envSens.sensors&TEMPSENS != 0 {
					envSensProp[TEMPSENSID] = respPkg.payload.cmdBody.values[i].first
					i++
				}
				if envSens.sensors&HUMIDITYSENS != 0 {
					envSensProp[HUMIDITYSENSID] = respPkg.payload.cmdBody.values[i].first
					i++
				}
				if envSens.sensors&LIGHTSENS != 0 {
					envSensProp[LIGHTSENSID] = respPkg.payload.cmdBody.values[i].first
					i++
				}
				if envSens.sensors&AIRPOLLSENS != 0 {
					envSensProp[AIRPOLLSENSID] = respPkg.payload.cmdBody.values[i].first
					i++
				}

				for _, trig := range envSens.triggers {
					if trig.op&((1<<2)|(1<<3)) == TEMPSENSID {
						if trig.op&(1<<1) == 0 {
							if envSensProp[TEMPSENSID] < trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						} else {
							if envSensProp[TEMPSENSID] > trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						}
					} else if trig.op&((1<<2)|(1<<3)) == HUMIDITYSENSID {
						if trig.op&(1<<1) == 0 {
							if envSensProp[HUMIDITYSENSID] < trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						} else {
							if envSensProp[HUMIDITYSENSID] > trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						}
					} else if trig.op&((1<<2)|(1<<3)) == LIGHTSENSID {
						if trig.op&(1<<1) == 0 {
							if envSensProp[LIGHTSENSID] < trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						} else {
							if envSensProp[LIGHTSENSID] > trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						}
					} else if trig.op&((1<<2)|(1<<3)) == AIRPOLLSENSID {
						if trig.op&(1<<1) == 0 {
							if envSensProp[AIRPOLLSENSID] < trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						} else {
							if envSensProp[AIRPOLLSENSID] > trig.value.first {
								npkg, ok := processingTrig(&trig, &devs, hubId, *hubSerialP)
								if ok {
									pkgsToSent = append(pkgsToSent, npkg)
								}
							}
						}
					}
				}
			}
		}
	}

	envSensorsP = &envSensors
	swithcersP = &swithcers
	devsP = &devs
	return pkgsToSent
}

func processingTrig(trig *trigger, devsP *map[string]device, hubId uleb128, hubSerial uleb128) (packet, bool) {
	ncmdBody := cmdBody{
		status: trig.op & (1),
	}
	if dev, ok := (*devsP)[trig.name]; ok {
		var npayload payload
		initPayload(&npayload, hubId, uleb128{uint64(dev.id), 0}, hubSerial, dev.devType, SETSTATUS, &ncmdBody)
		var npkg packet
		initPacket(&npkg, &npayload)

		return npkg, true

	}
	return packet{}, false
}

func initPacket(pkg *packet, payld *payload) {
	binNpayload := getBinPayload(payld)
	pkg.length = byte(len(binNpayload))
	pkg.payload = *payld
	pkg.crc8 = calcCRC8(binNpayload)
}

func initPayload(payld *payload, src uleb128, dst uleb128, serial uleb128, devType byte, cmd byte, cmdBody *cmdBody) {
	payld.src = src
	payld.dst = dst
	payld.serial = serial
	payld.devType = devType
	payld.cmd = cmd
	if cmdBody != nil {
		payld.cmdBody = *cmdBody
	}
}

func initCmdBody(cmdbody *cmdBody, timestamp uleb128, devName *string, devProps *devProps, values *[]uleb128, status byte) {
	cmdbody.timestamp = timestamp
	if devName != nil {
		cmdbody.devName = *devName
	}
	if devProps != nil {
		cmdbody.devProps = *devProps
	}
	if values != nil {
		cmdbody.values = *values
	}
	cmdbody.status = status
}

func getStartPkg(hubId uleb128) packet {
	payl := payload{
		src:     hubId,
		dst:     uleb128{16383, 0},
		serial:  uleb128{1, 0},
		devType: SMARTHUB,
		cmd:     WHOISHERE,
		cmdBody: cmdBody{
			devName: "HUB01",
		},
	}
	binPayld := getBinPayload(&payl)
	pkg := packet{
		length:  byte(len(binPayld)),
		payload: payl,
		crc8:    calcCRC8(binPayld),
	}
	return pkg
}

func calcCRC8(bytes []byte) byte {
	const generator byte = 0x1D
	var crc byte = 0

	for _, currByte := range bytes {
		crc ^= currByte

		for i := 0; i < 8; i++ {
			if (crc & 0x80) != 0 {
				crc = (crc << 1) ^ generator
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}
