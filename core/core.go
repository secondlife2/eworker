package core

import "strconv"

/* interface to int */
func StringToInt(ifc string) int {
	if ifc == "" {
		return 0
	} else {
		if i, err := strconv.Atoi(ifc); err == nil {
			return i
		}
		return 0
	}
}

/* interface to string */
func IntToString(ifc int) string {
	if ifc == 0 {
		return ""
	} else {
		str1 := strconv.Itoa(ifc)
		return str1
	}
}
