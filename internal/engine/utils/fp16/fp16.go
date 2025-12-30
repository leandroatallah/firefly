package fp16

const scale = 16

func To16(value int) int {
	return value * scale
}
func From16(value int) int {
	return value / scale
}
