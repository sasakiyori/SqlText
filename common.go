package sqltext

func SkipSpacesFromHead(s string) string {
	for i := range s {
		if !IsSpaces(s[i]) {
			return s[i:]
		}
	}
	return ""
}

func IsSpaces(ch byte) bool {
	switch ch {
	case ' ', '\t', '\n':
		return true
	default:
		return false
	}
}
