package bncode

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// Decode decodes a bencoded string into a bencode value
func Decode(bencodedString string) (interface{}, error) {
	r := strings.NewReader(bencodedString)
	return decodeNextValue(r, nil)
}

func decodeNextValue(r *strings.Reader, items []interface{}) (interface{}, error) {
	k, _, err := r.ReadRune()
	if err == io.EOF {
		return items, nil
	}
	if err != nil {
		return nil, err
	}
	err = r.UnreadRune()
	if err != nil {
		return nil, err
	}
	if unicode.IsDigit(k) {
		v, err := decodeString(r)
		if err != nil {
			return nil, err
		}
		if items != nil {
			items = append(items, v)
			return decodeNextValue(r, items)
		}
		return v, nil
	} else if k == 'i' {
		v, err := decodeInt(r)
		if err != nil {
			return nil, err
		}
		if items != nil {
			items = append(items, v)
			return decodeNextValue(r, items)
		}
		return v, nil
	} else if k == 'l' {
		v, err := decodeList(r)
		if err != nil {
			return nil, err
		}
		if items != nil {
			items = append(items, v)
			return decodeNextValue(r, items)
		}
		return v, nil
	} else if k == 'd' {
		v, err := decodeDict(r)
		if err != nil {
			return nil, err
		}
		if items != nil {
			items = append(items, v)
			return decodeNextValue(r, items)
		}
		return v, nil
	} else if k == 'e' {
		// remove the trailing e and proceed to next item
		r.ReadRune()
		return items, nil
	}
	return "", fmt.Errorf("only strings are supported at the moment")
}

func decodeString(r *strings.Reader) (string, error) {
	size := ""
	s := ""
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if c == ':' {
			l, err := strconv.Atoi(size)
			if err != nil {
				return "", err
			}
			sbuf := make([]byte, l)
			r.Read(sbuf)
			s = string(sbuf)
			break
		}
		size += string(c)
	}
	return string(s), nil
}

func decodeInt(r *strings.Reader) (int, error) {
	// make sure the first char is an i
	c, _, err := r.ReadRune()
	if err != nil {
		return 0, err
	}
	if c != 'i' {
		return 0, fmt.Errorf("%v invalid character for int", c)
	}
	num := ""
	n := 0
	for {
		c, _, err = r.ReadRune()
		if err != nil {
			return 0, err
		}
		if c == 'e' {
			n, err = strconv.Atoi(num)
			if err != nil {
				return 0, err
			}
			return n, nil
		}
		num += string(c)
	}
}

func decodeList(r *strings.Reader) (interface{}, error) {
	// remove the leading l
	r.ReadRune()
	list := make([]interface{}, 0)
	return decodeNextValue(r, list)
}

func decodeDict(r *strings.Reader) (interface{}, error) {
	// remove the leading d
	r.ReadRune()
	m := make(map[string]interface{}, 0)
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}
		if c == 'e' {
			return m, nil
		}
		r.UnreadRune()
		k, err := decodeNextValue(r, nil)
		if err != nil {
			return nil, err
		}
		v, err := decodeNextValue(r, nil)
		if err != nil {
			return nil, err
		}
		m[k.(string)] = v
	}
}
