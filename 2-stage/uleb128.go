package main

//type uleb128 struct {
//	first  uint64
//	second uint64
//}

//func hexStringToUint(hexStr string) (uint64, error) {
//	hexStr = strings.TrimPrefix(hexStr, "0x")
//
//	var result uint64
//
//	var pow16 uint64 = 1
//	for i := len(hexStr) - 1; i >= 0; i-- {
//		if (hexStr[i] >= '0') && (hexStr[i] <= '9') {
//			result += pow16 * uint64(hexStr[i]-'0')
//			pow16 *= 16
//		} else if (hexStr[i] >= 'a') && (hexStr[i] <= 'f') {
//			result += pow16 * uint64(hexStr[i]-'a'+10)
//			pow16 *= 16
//		} else {
//			return result, errors.New("invalid hex-symbol")
//		}
//	}
//	return result, nil
//}
//
//func EncodeUleb128(val uint64) []byte {
//	var result []byte
//	for {
//		b := byte(val & 0x7F)
//		val >>= 7
//		if val != 0 {
//			b |= 0x80
//		}
//		result = append(result, b)
//		if val == 0 {
//			break
//		}
//	}
//	return result
//}
//
//func DecodeUleb128(data []byte) uint64 {
//	var result uint64
//	var shift uint
//	for _, b := range data {
//		result |= uint64(b&0x7F) << shift
//		shift += 7
//		if b&0x80 == 0 {
//			break
//		}
//	}
//	return result
//}
//
//func getBytesFromUleb(u uleb128) []byte {
//	buf := make([]byte, 0, 16)
//	var encodedFirst []byte
//	var encodedSecond []byte
//	if u.second == 0 {
//		encodedFirst = EncodeUleb128(u.first)
//		buf = append(buf, encodedFirst...)
//		return buf
//	}
//	encodedFirst = EncodeUleb128(u.first)
//	encodedSecond = EncodeUleb128(u.second)
//	buf = append(buf, encodedFirst...)
//	buf = append(buf, encodedSecond...)
//	return buf
//}
//
//func getUlebFromBytes(bytes []byte) uleb128 { //bytes' len <= 16
//	var u uleb128
//	if len(bytes) <= 8 {
//		u.first = DecodeUleb128(bytes)
//		return u
//	}
//	u.first = DecodeUleb128(bytes[:8])
//	u.second = DecodeUleb128(bytes[8:])
//	return u
//}
//
//func printBytesHexf(mes string, mbytes []byte) {
//	fmt.Printf("%s: ", mes)
//	for _, b := range mbytes {
//		fmt.Printf("%2x ", b)
//	}
//	fmt.Printf("\n")
//}
//
//func ulebPlus(u *uleb128, val uint64) *uleb128 {
//	if u.first+val > math.MaxUint64 {
//		u.second += (u.first + val) - math.MaxUint64
//		u.first = math.MaxUint64
//	} else {
//		u.first += val
//	}
//	return u
//}
