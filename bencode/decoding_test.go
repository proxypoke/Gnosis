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
		pos_check  int64 = 1234567890
		neg_check  int64 = -1234567890
		zero_check int64 = 0
		// invalid streams
		no_starting_i []byte = []byte("x12345667890e")
		leading_zero  []byte = []byte("i000000e")
		negative_zero []byte = []byte("i-0e")
		invalid_chars []byte = []byte("ifoobare")
	)

	valid := [][]byte{
		positive,
		negative,
		zero}
	check := []int64{
		pos_check,
		neg_check,
		zero_check}
	invalid := [][]byte{
		no_starting_i,
		leading_zero,
		negative_zero,
		invalid_chars}

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
