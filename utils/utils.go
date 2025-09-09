package utils

func ReverseMap(m map[string]string) map[string]string {
	reversed := make(map[string]string)
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}
