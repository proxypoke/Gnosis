// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode

import (
	"strconv"
)

// Constant identifiers for the various Bencode structures.
const (
	INVALID = iota
	INT
	STRING
	LIST
	DICT
	TERMINATOR // termination character for int, list and dict ('e')
)

// Get the type of the stream. Returns one of the above identifiers.
func GetType(stream []byte) (type_ int) {
	return NewDecoder(stream).getType()
}

// ================================ [ DECODER ] ================================

type DecodeError string

func (self DecodeError) Error() string {
	return string(self)
}

type Decoder struct {
	stream []byte
	pos    int
}

// Construct a Decoder for a given stream.
func NewDecoder(stream []byte) *Decoder {
	return &Decoder{stream, 0}
}

// Decode a stream.
func Decode(stream []byte) (Interface, error) {
	return NewDecoder(stream).decodeNext()
}

// Decode an integer from a stream.
func DecodeInt(stream []byte) (Int, error) {
	return NewDecoder(stream).decodeInt()
}

// Decode a string from a stream.
func DecodeString(stream []byte) (String, error) {
	return NewDecoder(stream).decodeString()
}

// Decode a list from a stream.
func DecodeList(stream []byte) (List, error) {
	return NewDecoder(stream).decodeList()
}

// Decode a dict from a stream.
func DecodeDict(stream []byte) (Dict, error) {
	return NewDecoder(stream).decodeDict()
}

// ============================ [ DECODER METHODS ] ============================

// Internal getType method.
func (self *Decoder) getType() (type_ int) {
	if self.pos > len(self.stream) {
		type_ = INVALID
	}
	switch c := self.stream[self.pos]; {
	default:
		type_ = INVALID
	case c >= '0' && c <= '9':
		type_ = STRING
	case c == 'i':
		type_ = INT
	case c == 'l':
		type_ = LIST
	case c == 'd':
		type_ = DICT
	case c == 'e':
		type_ = TERMINATOR
	}
	return
}

// Internal method for getting the next decodable object.
// TODO: Add dicts once they are implemented.
func (self *Decoder) decodeNext() (result Interface, err error) {
	switch t := self.getType(); {
	case t == INT:
		result, err = self.decodeInt()
	case t == STRING:
		result, err = self.decodeString()
	case t == LIST:
		result, err = self.decodeList()
	case t == DICT:
		result, err = self.decodeDict()
	default:
		err = DecodeError("Cannot decode stream: invalid characters.")
	}
	return
}

// Internal decoder method for ints.
func (self *Decoder) decodeInt() (result Int, err error) {
	if self.getType() != INT {
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
			err = DecodeError("Cannot decode int: reached end of stream " +
				"before 'e' was found.")
			return
		}
	}

	str := string(self.stream[self.pos+1 : i])
	x, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		self.pos = i + 1
		result = Int(x)
	}
	return
}

// Internal decoding method for Strings.
func (self *Decoder) decodeString() (result String, err error) {
	i := self.pos
	for {
		c := self.stream[i]

		if c == ':' {
			break
		}
		if !isNum(c) {
			err = DecodeError(
				"Cannot decode string: invalid characters in length.")
			return
		}
		i++
		if i >= len(self.stream) {
			err = DecodeError("Cannot decode string: reached end of stream " +
				"before ':' was found.")
			return
		}
	}

	// Check if there have been any numbers at all.
	if i == self.pos {
		err = DecodeError("Cannot decode string: no length found.")
		return
	}

	lenstr := string(self.stream[self.pos:i])
	length, err := strconv.ParseInt(lenstr, 10, 64)
	if err != nil {
		return
	}

	i++
	var chars []byte
	for {
		if length == 0 {
			break
		}
		if i >= len(self.stream) {
			err = DecodeError("Cannot decode string: too short.")
			return
		}
		chars = append(chars, self.stream[i])
		length--
		i++
	}

	result = String(chars)
	self.pos = i
	return
}

// Internal decoding method for lists.
func (self *Decoder) decodeList() (result List, err error) {
	if self.getType() != LIST {
		err = DecodeError("Cannot decode list: doesn't start with 'l'.")
		return
	}

	// Just the nil element is a list containing nil, which isn't what we want.
	result = List{}

	self.pos++
	for {
		if self.pos >= len(self.stream) {
			err = DecodeError("Cannot decode list: reached end of stream " +
				"before terminating 'e' was found.")
			return
		}
		if self.stream[self.pos] == 'e' {
			break
		}

		var obj Interface
		obj, err = self.decodeNext()
		if err != nil {
			return
		}
		result = append(result, obj)
	}

	if err == nil {
		self.pos++
	}
	return
}

func (self *Decoder) decodeDict() (result Dict, err error) {
	if self.getType() != DICT {
		err = DecodeError("Cannot decode dict: doesn't start with 'd'.")
		return
	}

	// Create empty dict
	result = Dict{}

	self.pos++
	for {
		if self.pos >= len(self.stream) {
			err = DecodeError("Cannot decode dict: reached end of stream " +
				"before terminating 'e' was found.")
			return
		}
		if self.stream[self.pos] == 'e' {
			break
		}

		var str String
		str, err = self.decodeString()
		if err != nil {
			return
		}

		var obj Interface
		obj, err = self.decodeNext()
		if err != nil {
			return
		}

		result[str] = obj
	}

	if err == nil {
		self.pos++
	}
	return

	return
}

// =========================== [ HELPER FUNCTIONS ] ============================

// Check if a byte represents an ASCII number.
func isNum(b byte) bool {
	if b < '0' || b > '9' {
		return false
	}
	return true
}
