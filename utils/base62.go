package utils

func Base62Encode(num int64) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if num == 0 {
		return string(chars[0])
	}
	var result []byte
	for num > 0 {
		remainder := num % 62
		result = append([]byte{chars[remainder]}, result...)
		num /= 62
	}
	return string(result)
}