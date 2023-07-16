package main

//func getBinPkg(pkg packet) (buf []byte) {
//	payloadSize := pkg.length
//	buf = make([]byte, 0, 2+payloadSize)
//
//	buf = append(buf, pkg.length)
//
//	buf = append(buf, getBinPayload(&pkg.payload)...)
//
//	buf = append(buf, pkg.crc8)
//
//	return buf
//}
//
//func getBinPayload(payload1 *payload) []byte {
//	buf := make([]byte, 0, 5)
//	buf = append(buf, getBytesFromUleb(payload1.src)...)
//	buf = append(buf, getBytesFromUleb(payload1.dst)...)
//	buf = append(buf, getBytesFromUleb(payload1.serial)...)
//	buf = append(buf, payload1.devType)
//	buf = append(buf, payload1.cmd)
//	buf = append(buf, getBinCmdBody(&payload1.cmdBody, payload1.cmd, payload1.devType)...)
//
//	return buf
//}
//
//func getBinCmdBody(cbody *cmdBody, cmd uint8, devType uint8) []byte {
//	buf := make([]byte, 0)
//
//	if cmd == TICK && devType == CLOCK {
//		buf = append(buf, getBytesFromUleb(cbody.timestamp)...)
//	} else if cmd == GETSTATUS { //cmdBody is empty
//		return buf
//	} else if cmd == WHOISHERE || cmd == IAMHERE {
//		nameLen := len(cbody.devName)
//		buf = append(buf, byte(nameLen))
//		buf = append(buf, cbody.devName...)
//		if devType == ENVSENSOR {
//			buf = append(buf, cbody.devProps.sensors)
//			buf = append(buf, byte(len(cbody.devProps.triggers)))
//			for i, _ := range cbody.devProps.triggers {
//				buf = append(buf, cbody.devProps.triggers[i].op)
//				buf = append(buf, getBytesFromUleb(cbody.devProps.triggers[i].value)...)
//				nameLen = len(cbody.devProps.triggers[i].name)
//				buf = append(buf, byte(nameLen))
//				buf = append(buf, cbody.devProps.triggers[i].name...)
//			}
//		} else if devType == SWITCH {
//			buf = append(buf, byte(len(cbody.devProps.devNames)))
//			for i, _ := range cbody.devProps.devNames {
//				nameLen = len(cbody.devProps.devNames[i])
//				buf = append(buf, byte(nameLen))
//				buf = append(buf, cbody.devProps.devNames[i]...)
//			}
//		}
//	} else if cmd == STATUS {
//		if devType == ENVSENSOR {
//			buf = append(buf, byte(len(cbody.values)))
//			for i, _ := range cbody.values {
//				buf = append(buf, getBytesFromUleb(cbody.values[i])...)
//			}
//		} else if devType == SWITCH || devType == LAMP || devType == SOCKET {
//			buf = append(buf, cbody.status)
//		}
//	} else if cmd == SETSTATUS {
//		buf = append(buf, cbody.status)
//	}
//
//	return buf
//}
//
//func getPkgsFromBytes(buf []byte) []packet {
//	pkgs := make([]packet, 0, 1)
//
//	for len(buf) > 0 {
//		pkg := packet{}
//		pkg.length = buf[0]
//		buf = buf[1:]
//
//		pkg.payload = getPayloadFromBytes(buf[0:pkg.length])
//		buf = buf[pkg.length:]
//		pkg.crc8 = buf[0]
//		buf = buf[1:]
//		pkgs = append(pkgs, pkg)
//	}
//
//	return pkgs
//}
//
//func getPayloadFromBytes(buf []byte) payload {
//	payld := payload{}
//	if buf[0] >= 128 {
//		payld.src = getUlebFromBytes(buf[0:2])
//		buf = buf[2:]
//	} else {
//		payld.src = getUlebFromBytes(buf[0:1])
//		buf = buf[1:]
//	}
//	if buf[0] >= 128 {
//		payld.dst = getUlebFromBytes(buf[0:2])
//		buf = buf[2:]
//	} else {
//		payld.dst = getUlebFromBytes(buf[0:1])
//		buf = buf[1:]
//	}
//	ulebBytes := make([]byte, 0)
//	for buf[0] >= 128 {
//		ulebBytes = append(ulebBytes, buf[0])
//		buf = buf[1:]
//	}
//	ulebBytes = append(ulebBytes, buf[0])
//	buf = buf[1:]
//	payld.serial = getUlebFromBytes(ulebBytes)
//
//	payld.devType = buf[0]
//	buf = buf[1:]
//	payld.cmd = buf[0]
//	buf = buf[1:]
//	payld.cmdBody = getCmdBodyFromBytes(buf, payld.cmd, payld.devType)
//
//	return payld
//}
//
//func getCmdBodyFromBytes(buf []byte, cmd byte, devType byte) cmdBody {
//	cbody := cmdBody{}
//	if cmd == TICK && devType == CLOCK {
//		ulebBytes := make([]byte, 0)
//		for buf[0] >= 128 {
//			ulebBytes = append(ulebBytes, buf[0])
//			buf = buf[1:]
//		}
//		ulebBytes = append(ulebBytes, buf[0])
//		buf = buf[1:]
//		cbody.timestamp = getUlebFromBytes(ulebBytes)
//	} else if cmd == GETSTATUS {
//		return cbody
//	} else if cmd == WHOISHERE || cmd == IAMHERE {
//		nameLen := buf[0]
//		buf = buf[1:]
//		devName := make([]byte, nameLen)
//		for i := 0; i < int(nameLen); i++ {
//			devName[i] = buf[i]
//		}
//		buf = buf[nameLen:]
//		cbody.devName = string(devName)
//		if devType == ENVSENSOR {
//			cbody.devProps.sensors = buf[0]
//			buf = buf[1:]
//			triggersLen := buf[0]
//			buf = buf[1:]
//			cbody.devProps.triggers = make(map[string]trigger)
//			for i := 0; i < int(triggersLen); i++ {
//				trig := trigger{}
//				trig.op = buf[0]
//				buf = buf[1:]
//				ulebBytes := make([]byte, 0)
//				for buf[0] >= 128 {
//					ulebBytes = append(ulebBytes, buf[0])
//					buf = buf[1:]
//				}
//				ulebBytes = append(ulebBytes, buf[0])
//				buf = buf[1:]
//				trig.value = getUlebFromBytes(ulebBytes)
//				nameLen := buf[0]
//				buf = buf[1:]
//				devName := make([]byte, nameLen)
//				for i := 0; i < int(nameLen); i++ {
//					devName[i] = buf[i]
//				}
//				buf = buf[nameLen:]
//				trig.name = string(devName)
//				cbody.devProps.triggers[trig.name] = trig
//			}
//		} else if devType == SWITCH {
//			devNamesLen := buf[0]
//			buf = buf[1:]
//			for i := 0; i < int(devNamesLen); i++ {
//				nameLen := buf[0]
//				buf = buf[1:]
//				devName := make([]byte, nameLen)
//				for i := 0; i < int(nameLen); i++ {
//					devName[i] = buf[i]
//				}
//				buf = buf[nameLen:]
//				cbody.devProps.devNames = append(cbody.devProps.devNames, string(devName))
//			}
//		}
//	} else if cmd == STATUS {
//		if devType == ENVSENSOR {
//			valsLen := buf[0]
//			buf = buf[1:]
//			for i := 0; i < int(valsLen); i++ {
//				ulebBytes := make([]byte, 0)
//				for buf[0] >= 128 {
//					ulebBytes = append(ulebBytes, buf[0])
//					buf = buf[1:]
//				}
//				ulebBytes = append(ulebBytes, buf[0])
//				buf = buf[1:]
//				cbody.values = append(cbody.values, getUlebFromBytes(ulebBytes))
//			}
//		} else if devType == SWITCH || devType == LAMP || devType == SOCKET {
//			cbody.status = buf[0]
//			buf = buf[1:]
//		}
//	}
//	return cbody
//}
