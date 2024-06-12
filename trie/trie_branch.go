package trie

import "fmt"

func NewBranchNode() *Node {
	index := [16]int {}
	for i:=0;i<16;i++{
		index[i] = -1
	}
	return &Node{
		Flag: Branch,
		Branches: []*Node{},
		Index: index,
	}
}

func (b *Node) SetBranch(nibble byte, node *Node) error {
	if b.Flag != Branch {
		return fmt.Errorf("only branch node can call SetBranch")
	}
	if b.Index[int(nibble)] != -1 {
		b.Branches[int(nibble)] = node
	}else {
		count := 0
		for i:=int(nibble)+1; i<16; i++{
			if b.Index[i] == -1 {continue}
			b.Index[i]++
			count++
		}
		front := len(b.Branches) - count
		b.Index[int(nibble)] = front
		branches := make([]*Node,front)
		copy(branches,b.Branches[:front])
		branches = append(branches, node)
		b.Branches = append(branches, b.Branches[front:]...)
		// if front == 0 {
		// 	b.Branches = append([]*Node {node},b.Branches...)
		// }else {
		// 	b.Branches = append(append(b.Branches[:front], node), b.Branches[front:]...)
		// }
	}
	return nil
}

func (b *Node) RemoveBranch(nibble byte) error {
	if b.Flag != Branch {
		return fmt.Errorf("only branch node can call RemoveBranch")
	}
	b.Branches = append(b.Branches[:b.Index[int(nibble)]], b.Branches[b.Index[int(nibble)]+1:]...)
	b.Index[int(nibble)] = -1
	for i:=int(nibble)+1; i<16; i++{
		if b.Index[i] == -1 {continue}
		b.Index[i]--
	}
	return nil
}

func (b *Node) SetValue(value []byte) error {
	if b.Flag != Branch {
		return fmt.Errorf("only branch node can call SetValue")
	}
	b.Value = value
	return nil
}

func (b *Node) RemoveValue() error {
	if b.Flag != Branch {
		return fmt.Errorf("only branch node can call RemoveValue")
	}
	b.Value = nil
	return nil
}

func (b *Node) HasValue() (bool, error) {
	if b.Flag != Branch {
		return false, fmt.Errorf("only branch node can call HasValue")
	}
	return b.Value != nil, nil
}