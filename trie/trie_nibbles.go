package trie

import (
	"fmt"
)

func IsNibble(nibble byte) bool {
	n := int(nibble)
	// 0-9 && a-f
	return n >= 0 && n < 16
}

func FromNibbleByte(n byte) (byte, error) {
	if !IsNibble(n) {
		return 0, fmt.Errorf("non-nibble byte: %v", n)
	}
	return byte(n), nil
}

// nibbles contain one nibble per byte
func FromNibbleBytes(nibbles []byte) ([]byte, error) {
	ns := make([]byte, 0, len(nibbles))
	for _, n := range nibbles {
		nibble, err := FromNibbleByte(n)
		if err != nil {
			return nil, fmt.Errorf("contains non-nibble byte: %w", err)
		}
		ns = append(ns, nibble)
	}
	return ns, nil
}

func FromByte(b byte) []byte {
	return []byte{
		byte(byte(b >> 4)),
		byte(byte(b % 16)),
	}
}

func FromBytes(bs []byte) []byte {
	ns := make([]byte, 0, len(bs)*2)
	for _, b := range bs {
		ns = append(ns, FromByte(b)...)
	}
	return ns
}

func FromString(s string) []byte {
	return FromBytes([]byte(s))
}

// ToPrefixed add nibble prefix to a slice of nibbles to make its length even
// the prefix indicts whether a node is a leaf node.
func ToPrefixed(ns []byte, isLeafNode bool) []byte {
	// create prefix
	var prefixBytes []byte
	// odd number of nibbles
	if len(ns)%2 > 0 {
		prefixBytes = []byte{1}
	} else {
		// even number of nibbles
		prefixBytes = []byte{0, 0}
	}

	// append prefix to all nibble bytes
	prefixed := make([]byte, 0, len(prefixBytes)+len(ns))
	prefixed = append(prefixed, prefixBytes...)
	for _, n := range ns {
		prefixed = append(prefixed, byte(n))
	}

	// update prefix if is leaf node
	if isLeafNode {
		prefixed[0] += 2
	}

	return prefixed
}

// ToBytes converts a slice of nibbles to a byte slice
// assuming the nibble slice has even number of nibbles.
func ToBytes(ns []byte) []byte {
	buf := make([]byte, 0, len(ns)/2)

	for i := 0; i < len(ns); i += 2 {
		b := byte(ns[i]<<4) + byte(ns[i+1])
		buf = append(buf, b)
	}

	return buf
}

// [0,1,2,3], [0,1,2] => 3
// [0,1,2,3], [0,1,2,3] => 4
// [0,1,2,3], [0,1,2,3,4] => 4
func PrefixMatchedLen(node1 []byte, node2 []byte) int {
	matched := 0
	for i := 0; i < len(node1) && i < len(node2); i++ {
		n1, n2 := node1[i], node2[i]
		if n1 == n2 {
			matched++
		} else {
			break
		}
	}

	return matched
}
