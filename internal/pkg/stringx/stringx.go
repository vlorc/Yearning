package stringx

func Coalesce(s ...string) string {
	for _, v := range s {
		if "" != v {
			return v
		}
	}
	return ""
}
