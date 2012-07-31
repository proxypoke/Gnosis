// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode_test

import (
	"bytes"
	. "gnosis/bencode"
	"sort"
	"strconv"
	"testing"
)

func TestStringEncode(t *testing.T) {
	str := String("foobar")
	bstr := Encode(str)
	expected := []byte("6:foobar")
	if !bytes.Equal(bstr, expected) {
		t.Errorf("encoded string isn't a valid bencoding: %s", bstr)
		t.Logf("Expected result: %s", expected)
	}
}

func TestIntEncode(t *testing.T) {
	var i Int
	// test a negative (-1) and positive (1) number as well as zero
	for i = -1; i <= 1; i++ {
		bint := Encode(i)
		expected := []byte("i" + strconv.Itoa(int(i)) + "e")
		if !bytes.Equal(bint, expected) {
			t.Errorf("encoded int isn't a valid bencoding: %s", bint)
			t.Logf("Expected result: %s", expected)
		}
	}
}

func TestListEncode(t *testing.T) {
	var list List
	expected := []byte("l")

	// Add 10 integers to the list
	for i := 0; i < 10; i++ {
		list = append(list, Int(i))
		expected = append(expected, Encode(Int(i))...)
	}
	// Add some strings to the list
	strings := []string{"foo", "bar", "baz"}
	for _, str := range strings {
		list = append(list, String(str))
		expected = append(expected, Encode(String(str))...)
	}
	expected = append(expected, 'e')
	blist := Encode(list)
	if !bytes.Equal(blist, expected) {
		t.Errorf("encoded list isn't a valid bencoding: %s", blist)
		t.Logf("Expected result: %s", expected)
	}
}

func TestDictEncode(t *testing.T) {
	dict := make(Dict)
	expected := []byte("d")

	// Add 10 integers to the dict
	for i := 0; i < 10; i++ {
		key := String(strconv.Itoa(i))
		dict[key] = Int(i)
		expected = append(expected, Encode(key)...)
		expected = append(expected, Encode(Int(i))...)
	}
	// Add some strings to the dict
	strings := []string{"foo", "bar", "baz"}
	for _, str := range strings {
		key := String(str)
		dict[key] = String(str)
	}
	// Add the strings to expected in lexigraphic order
	sort.Strings(strings)
	for _, str := range strings {
		key := Encode(String(str))
		val := key
		expected = append(expected, key...)
		expected = append(expected, val...)
	}
	expected = append(expected, 'e')
	bdict := Encode(dict)
	if !bytes.Equal(bdict, expected) {
		t.Errorf("encoded dict isn't a valid bencoding: %s:", bdict)
		t.Logf("Expected result: %s", expected)
	}
}
