package betterencode

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
)

var esafe = base64.NewEncoding("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789><").WithPadding(-1)

//EncodeURLSafe : better urlsafe encoding
func EncodeURLSafe(data []byte) string {
	tmp := esafe.EncodeToString(data)
	var last byte
	var counter int
	r := make([]byte, 0, len(tmp))
	for i := range tmp {
		if tmp[i] == last {
			counter++
		} else {
			if counter >= 3 {
				for counter > 9 {
					r = append(r, last, 42, 57)
					counter -= 9
				}
				if counter >= 3 {
					r = append(r, last, 42, []byte(strconv.Itoa(counter))[0])
				} else if counter > 0 {
					for i := 0; i < counter; i++ {
						r = append(r, last)
					}
				}
				last = tmp[i]
				counter = 1
			} else {
				for i := 0; i < counter; i++ {
					r = append(r, last)
				}
				last = tmp[i]
				counter = 1
			}
		}
	}
	if counter >= 3 {
		for counter > 9 {
			r = append(r, last, 42, 57)
			counter -= 9
		}
		if counter >= 3 {
			r = append(r, last, 42, []byte(strconv.Itoa(counter))[0])
		} else if counter > 0 {
			for i := 0; i < counter; i++ {
				r = append(r, last)
			}
		}
	} else {
		for i := 0; i < counter; i++ {
			r = append(r, last)
		}
	}
	tmp = string(r)
	tmp = strings.Replace(tmp, "qq", "!", -1)
	tmp = strings.Replace(tmp, "<<", "~", -1)
	return tmp
}

var cmatch, _ = regexp.Compile(`.\*[0-9]`)

//DecodeURLSafe : better urlsafe encoding
func DecodeURLSafe(data string) ([]byte, error) {
	data = strings.Replace(data, "!", "qq", -1)
	data = strings.Replace(data, "~", "<<", -1)
	data = string(cmatch.ReplaceAllFunc([]byte(data), func(b []byte) []byte {
		count, _ := strconv.Atoi(string(data[2]))
		return bytes.Repeat(b[:1], count)
	}))
	tmp, err := esafe.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}
