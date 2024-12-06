package trie

func NewExtensionNode(nibbles []byte, next *Node) *Node {
	return &Node{
		Flag:     Extension,
		KeySlice: nibbles,
		Next:     next,
	}
}
