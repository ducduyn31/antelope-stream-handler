package common

import "strings"

func StrArrContains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StrArrRemove(a string, list []string) []string {
	for i, s := range list {
		if a == s {
			list[i] = list[len(list) - 1]
			return list[:len(list) - 1]
		}
	}
	return list
}

func ConcatByteArr(bytes ...[]byte) []byte {
	var result []byte
	for _, el := range bytes {
		result = append(result, el...)
	}
	return result
}

func ConcatStringArr(strArr []string, char string) string {
	var sb strings.Builder
	for _, el := range strArr {
		if sb.Len() == 0 {
			sb.WriteString(el)
		} else {
			sb.WriteString(char + el)
		}
	}

	return sb.String()
}
