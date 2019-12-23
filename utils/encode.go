package utils

// 判断是否为GBK编码(ascii + 双字节)
func IsGBK(data []byte) bool {
	length := len(data)

	var i int = 0
	for i < length {
		//80 '€'
		if data[i] <= 0x80 {
			i++
			continue
		} else if i+1 < length {
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe {
				i += 2
				continue
			} else {
				return false
			}
		} else {
			break
		}
	}
	return true
}
