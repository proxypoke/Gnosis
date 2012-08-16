// Author: slowpoke <mail at slowpoke dot io>
//
// This program is free software under the non-terms
// of the Anti-License. Do whatever the fuck you want.

package bencode_test

import (
	. "gnosis/bencode"
	"testing"
)

func TestIntDecode(t *testing.T) {
	var (
		// valid streams
		positive []byte = []byte("i1234567890e")
		negative []byte = []byte("i-1234567890e")
		zero     []byte = []byte("i0e")
		// check values
		pos_check  Int = 1234567890
		neg_check  Int = -1234567890
		zero_check Int = 0
		// invalid streams
		no_starting_i []byte = []byte("x12345667890e")
		leading_zero  []byte = []byte("i000000e")
		negative_zero []byte = []byte("i-0e")
		invalid_chars []byte = []byte("ifoobare")
		too_short     []byte = []byte("i1234")
	)

	valid := [][]byte{
		positive,
		negative,
		zero,
	}
	check := []Int{
		pos_check,
		neg_check,
		zero_check,
	}
	invalid := [][]byte{
		no_starting_i,
		leading_zero,
		negative_zero,
		invalid_chars,
		too_short,
	}

	for i, stream := range valid {
		result, err := DecodeInt(stream)
		if err != nil {
			t.Errorf("Couldn't decode valid stream: %s", stream)
		}
		if check[i] != result {
			t.Errorf("Result (%d) doesn't match expected value (%d)",
				result, check[i])
		}
	}

	for _, stream := range invalid {
		_, err := DecodeInt(stream)
		if err == nil {
			t.Errorf("Didn't error on invalid stream: %s", stream)
		}
	}
}

func TestStringDecode(t *testing.T) {
	valid := [][]byte{
		[]byte("6:foobar"),
		[]byte("0:"),
	}
	check := []String{
		"foobar",
		"",
	}

	invalid := [][]byte{
		// No length
		[]byte(":dungoofd"),
		// Too short
		[]byte("9001:overninethousand...not"),
	}

	for i, stream := range valid {
		result, err := DecodeString(stream)
		if err != nil {
			t.Errorf("Couldn't decode valid stream: %s", stream)
		}
		if check[i] != result {
			t.Errorf("Result (%#v) doesn't match expected value (%#v)",
				result, check[i])
		}
	}

	for _, stream := range invalid {
		_, err := DecodeString(stream)
		if err == nil {
			t.Errorf("Didn't error on invalid stream: %s", stream)
		}
	}

}

func TestListDecode(t *testing.T) {
	// TODO: Add dicts once they are implemented.
	valid := [][]byte{
		[]byte("le"),            // empty list
		[]byte("l3:fooe"),       // list containing a string
		[]byte("li42ee"),        // list containing an int
		[]byte("llee"),          // list containing another list
		[]byte("l3:fooi42elee"), // comination of all of the above
	}
	/*
		// It's not possible to (trivially) compare slices, so we'll skip
		// checks.
		check := []List{
			List{},
			List{String("foo")},
			List{Int(42)},
			List{List{}},
			List{String("Foo"), Int(42), List{}},
		}
	*/

	invalid := [][]byte{
		[]byte("l"), // unterminated list
	}

	//for i, stream := range valid {
	for _, stream := range valid {
		//result, err := DecodeList(stream)
		_, err := DecodeList(stream)
		if err != nil {
			t.Errorf("Couldn't decode valid stream: %s", stream)
			t.Errorf("%#v", err)
		}
		/*
			// Catch panics when one of the results doesn't have enough or too
			// many elements.
			defer func() {
				if err := recover(); err != nil {
					t.Error("Encountered panic while comparing result to",
						"expected value. The result either had not enough",
						"or too many elements.")
					t.Logf("%#v", err)
				}
			}()

			for j, value := range check[i] {
				if value != result[j] {
					t.Errorf("Result (%#v) doesn't match expected value (%#v)",
						result, check[i])
				}
			}
		*/
	}

	for _, stream := range invalid {
		_, err := DecodeList(stream)
		if err == nil {
			t.Errorf("Didn't error on invalid stream: %s", stream)
		}
	}

}
