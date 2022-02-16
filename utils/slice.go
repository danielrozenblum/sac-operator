package utils

func Subtruct(a, b []string) []string {
	var c []string
	for i := range a {
		exist := false
		for i2 := range b {
			if a[i] == b[i2] {
				exist = true
			}
		}
		if !exist {
			c = append(c, a[i])
		}
	}
	return c
}
