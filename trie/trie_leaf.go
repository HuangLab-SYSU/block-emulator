package trie

import "fmt"

func NewLeafNodeFromNibbleBytes(nibbles []byte, value []byte) (*Node, error) {
	ns, err := FromNibbleBytes(nibbles)
	if err != nil {
		return nil, fmt.Errorf("could not leaf node from nibbles: %w", err)
	}

	return NewLeafNodeFromNibbles(ns, value), nil
}

func NewLeafNodeFromNibbles(nibbles []byte, value []byte) *Node {
	return &Node{
		Flag:  Leaf,
		KeySlice:  nibbles,
		Value: value,
	}
}

func NewLeafNodeFromKeyValue(key, value string) *Node {
	return NewLeafNodeFromBytes([]byte(key), []byte(value))
}

func NewLeafNodeFromBytes(key, value []byte) *Node {
	return NewLeafNodeFromNibbles(FromBytes(key), value)
}