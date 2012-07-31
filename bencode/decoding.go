// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode

import (
	"strconv"
)

type DecodeError string

func (self DecodeError) Error() string {
	return string(self)
}

type Decoder struct {
	stream []byte
	pos    int
}

func NewDecoder(stream []byte) *Decoder {
	return &Decoder{stream, 0}
}

// Decode an integer from a stream.
func DecodeInt(stream []byte) (int64, error) {
	return NewDecoder(stream).decodeInt()
}

func (self *Decoder) decodeInt() (result int64, err error) {
	if self.stream[self.pos] != 'i' {
		err = DecodeError("Cannot decode int: doesn't start with 'i'.")
		return
	}

	i := self.pos + 1
	// check for negative number
	negative := false
	if self.stream[i] == '-' {
		negative = true
		i++
	}

	// check for leading zeros
	zero := false
	if self.stream[i] == '0' {
		if !negative {
			zero = true
			i++
		} else {
			err = DecodeError("Cannot decode int: negative zero not allowed.")
			return
		}
	}
	if zero && self.stream[i] != 'e' {
		err = DecodeError("Cannot decode int: leading zeros not allowed.")
		return
	}

	for {
		c := self.stream[i]

		if c == 'e' {
			break
		}
		if c < '0' || c > '9' {
			err = DecodeError("Cannot decode int: invalid characters.")
			return
		}

		i++
		if i >= len(self.stream) {
			err = DecodeError("Cannot decode string: reached end of stream " +
				"before 'e' was found.")
			return
		}
	}

	str := string(self.stream[self.pos+1 : i])
	result, err = strconv.ParseInt(str, 10, 64)
	if err == nil {
		self.pos = i + 1
	}
	return
}