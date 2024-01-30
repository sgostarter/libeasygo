package strutils

import "strings"

// StringTrim .
func StringTrim(s string) string {
	return strings.Trim(s, "\r\n\t ")
}

// StringTrimBlank .
func StringTrimBlank(s string) string {
	return strings.Trim(s, " ")
}

/*
abc "aa bb" dd 'aaa "ab" xx ' "aaa bbb 'x x ' zz"
aa"xx
*/

func StringSplit(s string) (rets []string) {
	var inSubString bool

	var stringBorder rune

	var curSlice []rune

	for _, r := range s {
		if !inSubString {
			if r == ' ' {
				if len(curSlice) > 0 {
					rets = append(rets, string(curSlice))

					curSlice = nil
				}

				continue
			}

			if r == '"' || r == '\'' {
				stringBorder = r
				inSubString = true

				if len(curSlice) > 0 {
					rets = append(rets, string(curSlice))

					curSlice = nil
				}

				continue
			}

			curSlice = append(curSlice, r)
		} else {
			if r == stringBorder {
				if len(curSlice) > 0 {
					rets = append(rets, string(curSlice))
					curSlice = nil
				} else {
					rets = append(rets, "")
				}

				inSubString = false
				stringBorder = 0

				continue
			}

			curSlice = append(curSlice, r)
		}
	}

	if len(curSlice) > 0 {
		rets = append(rets, string(curSlice))
	}

	if len(rets) == 0 {
		rets = append(rets, "")
	}

	return
}
