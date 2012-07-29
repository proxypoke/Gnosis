// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode

import (
	"sort"
	"strconv"
)

// An object that can be bencoded (string, integer, slice or map)
type Interface interface {
	Encode() string
}

// Bencode type for a string.
type String string

func (self String) Encode() string {
	length := strconv.Itoa(len(self))
	return length + ":" + string(self)
}

// Bencode type for an int
type Int int

func (self Int) Encode() string {
	return "i" + strconv.Itoa(int(self)) + "e"
}

// Bencode type for a list
type List []Interface

func (self List) Encode() string {
	result := "l"
	for _, thing := range self {
		result += thing.Encode()
	}
	result += "e"
	return result
}

// Bencode type for a dict
type Dict map[String]Interface

func (self Dict) Encode() string {
	var keys sort.StringSlice
	for key := range self {
		keys = append(keys, string(key))
	}
	// keys must be in lexiographic order
	keys.Sort()

	result := "d"
	for _, key := range keys {
		result += String(key).Encode()
		val := Interface(self[String(key)])
		result += val.Encode()
	}
	result += "e"
	return result
}

// Wrapper function for a more consistent interface.
func Encode(data Interface) string {
	return data.Encode()
}
