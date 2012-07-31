// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode

import (
	"fmt"
	"sort"
)

// An object that can be bencoded (string, integer, slice or map)
type Interface interface {
	Encode() []byte
}

// Bencode type for a string.
type String string

func (self String) Encode() []byte {
	length := len(self)
	result := fmt.Sprintf("%d:%s", length, string(self))
	return []byte(result)
}

// Bencode type for an int
type Int int64

func (self Int) Encode() []byte {
	result := fmt.Sprintf("i%de", self)
	return []byte(result)
}

// Bencode type for a list
type List []Interface

func (self List) Encode() []byte {
	result := []byte("l")
	for _, thing := range self {
		result = append(result, Encode(thing)...)
	}
	result = append(result, 'e')
	return result
}

// Bencode type for a dict
type Dict map[String]Interface

func (self Dict) Encode() []byte {
	// keys must be in lexiographic order
	var keys sort.StringSlice
	for key := range self {
		keys = append(keys, string(key))
	}
	keys.Sort()

	result := []byte("d")
	for _, key := range keys {
		result = append(result, String(key).Encode()...)
		val := Interface(self[String(key)])
		result = append(result, val.Encode()...)
	}
	result = append(result, 'e')
	return []byte(result)
}

// Wrapper function for a more consistent interface.
func Encode(data Interface) []byte {
	return data.Encode()
}
