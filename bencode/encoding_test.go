// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode_test

import (
	. "gnosis/bencode"
	"sort"
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	str := String("foobar")
	bstr := str.Encode()
	expected := "6:foobar"
	if bstr != "6:foobar" {
		t.Errorf("Encoded string isn't a valid bencoding: %s", bstr)
		t.Logf("Expected result: %s", expected)
	}
}

func TestInt(t *testing.T) {
	var i Int
	// test a negative (-1) and positive (1) number as well as zero
	for i = -1; i <= 1; i++ {
		bint := i.Encode()
		expected := "i" + strconv.Itoa(int(i)) + "e"
		if bint != expected {
			t.Errorf("Encoded int isn't a valid bencoding: %s", bint)
			t.Logf("Expected result: %s", expected)
		}
	}
}

func TestList(t *testing.T) {
	var list List
	expected := "l"

	// Add 10 integers to the list
	for i := 0; i < 10; i++ {
		list = append(list, Int(i))
		expected += Int(i).Encode()
	}
	// Add some strings to the list
	strings := []string{"foo", "bar", "baz"}
	for _, str := range strings {
		list = append(list, String(str))
		expected += String(str).Encode()
	}
	expected += "e"
	blist := list.Encode()
	if blist != expected {
		t.Errorf("Encoded list isn't a valid bencoding: %s", blist)
		t.Logf("Expected result: %s", expected)
	}
}

func TestDict(t *testing.T) {
	dict := make(Dict)
	expected := "d"

	// Add 10 integers to the dict
	for i := 0; i < 10; i++ {
		key := String(strconv.Itoa(i))
		dict[key] = Int(i)
		expected += key.Encode() + Int(i).Encode()
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
		key := String(str).Encode()
		val := key
		expected += key + val
	}
	expected += "e"

	bdict := dict.Encode()
	if bdict != expected {
		t.Errorf("Encoded dict isn't a valid bencoding: %s:", bdict)
		t.Logf("Expected result: %s", expected)
	}
}
