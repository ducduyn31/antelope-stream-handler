package common

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
