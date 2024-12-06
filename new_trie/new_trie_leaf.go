package new_trie

import "fmt"

func NewLeafNodeFromNibbleBytes(nibbles []byte, value []byte, epochID int) (*Node, error) {
	ns, err := FromNibbleBytes(nibbles)
	if err != nil {
		return nil, fmt.Errorf("could not leaf node from nibbles: %w", err)
	}

	return NewLeafNodeFromNibbles(ns, value, epochID), nil
}

func NewLeafNodeFromNibbles(nibbles []byte, value []byte, epochID int) *Node {
	return &Node{
		Flag:     Leaf,
		KeySlice: nibbles,
		Value:    value,
		Epoch:    epochID,
	}
}

func NewLeafNodeFromKeyValue(key, value string, epochID int) *Node {
	return NewLeafNodeFromBytes([]byte(key), []byte(value), epochID)
}

func NewLeafNodeFromBytes(key, value []byte, epochID int) *Node {
	return NewLeafNodeFromNibbles(FromBytes(key), value, epochID)
}
