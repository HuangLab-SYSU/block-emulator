package new_trie

func NewExtensionNode(nibbles []byte, next *Node, epochID int) *Node {
	return &Node{
		Flag:     Extension,
		KeySlice: nibbles,
		Next:     next,
		Epoch:    epochID,
	}
}
