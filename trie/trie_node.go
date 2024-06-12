package trie

import (
	"blockEmulator/account"
	"encoding/hex"
	"fmt"
)

const (
	Leaf      = "leaf"
	Branch    = "branch"
	Extension = "extension"
)

type Node struct {
	//To indicate which type the node is
	//  leaf branch extension
	Flag     string
	Next     *Node   //For extension
	KeySlice []byte  //For leaf/extension
	Value    []byte  //For leaf/branch
	Index    [16]int //For branch
	Branches []*Node //For branch

}

func (node *Node) PrintState(key []byte) {
	if node.Flag == Leaf {
		key = append(key, node.KeySlice...)
		fmt.Printf("Key: %v, Value: %v\n", hex.EncodeToString(ToBytes(key)), account.DecodeAccountState(node.Value))
	} else if node.Flag == Extension {
		key = append(key, node.KeySlice...)
		node.Next.PrintState(key)
	} else {
		if has, _ := node.HasValue(); has {
			fmt.Printf("Key: %v, Value: %v\n", hex.EncodeToString(ToBytes(key)), account.DecodeAccountState(node.Value))
		}
		for i := 0; i < 16; i++ {
			if node.Index[i] != -1 {
				tmp := append(key, byte(i))
				node.Branches[node.Index[i]].PrintState(tmp)
			}
		}
	}
}

func (node *Node) Delete(key []byte) (*Node, error) {
	if node.Flag == Leaf {
		node = nil
		return node, nil
	}
	if node.Flag == Extension {
		length := len(node.KeySlice)
		next, err := node.Next.Delete(key[length:])
		if err != nil {
			return nil, err
		}
		if next.Flag == Leaf {
			node.Flag = Leaf
			node.Next = nil
			node.KeySlice = append(node.KeySlice, next.KeySlice...)
			node.Value = next.Value
			return node, nil
		}
		if next.Flag == Extension {
			node.Next = next.Next
			node.KeySlice = append(node.KeySlice, next.KeySlice...)
			return node, nil
		}
		if next.Flag == Branch {
			node.Next = next
			return node, nil
		}
	}
	if node.Flag == Branch {
		if len(key) == 0 {
			node.Value = nil
			//如果不止一个分支
			if len(node.Branches) != 1 {
				return node, nil
			}
			//如果只有一个分支
			nibble := 0
			for node.Index[nibble] == -1 {
				nibble++
			}
			branch := node.Branches[0]
			if branch.Flag == Leaf {
				node.Flag = Leaf
				node.Index = [16]int{}
				node.KeySlice = append([]byte{byte(nibble)}, branch.KeySlice...)
				node.Value = branch.Value
				node.Branches = nil
				return node, nil
			}
			if branch.Flag == Extension {
				node.Flag = Extension
				node.Next = branch.Next
				node.Index = [16]int{}
				node.KeySlice = append([]byte{byte(nibble)}, branch.KeySlice...)
				node.Branches = nil
				return node, nil
			}
			if branch.Flag == Branch {
				node.Flag = Extension
				node.Next = branch
				node.Index = [16]int{}
				node.KeySlice = []byte{byte(nibble)}
				node.Branches = nil
				return node, nil
			}
		}
		if len(key) != 0 {
			if len(node.Branches) == 1 {
				node.Flag = Leaf
				node.Branches = nil
				node.Index = [16]int{}
				return node, nil
			}
			if len(node.Branches) != 1 {
				nibble := key[0]
				branch := node.Branches[node.Index[int(nibble)]]
				if branch.Flag == Leaf {
					node.Branches = append(node.Branches[:node.Index[int(nibble)]], node.Branches[node.Index[int(nibble)]+1:]...)
					for k := range node.Index {
						if node.Index[k] > node.Index[int(nibble)] {
							node.Index[k] -= 1
						}
					}
					node.Index[int(nibble)] = -1
					length := len(node.Branches)
					if length == 1 {
						if node.Value != nil {
							return node, nil
						}
						if node.Value == nil {
							nibble := 0
							for node.Index[nibble] == -1 {
								nibble++
							}
							branch := node.Branches[0]
							if branch.Flag == Leaf {
								node.Flag = Leaf
								node.Index = [16]int{}
								node.KeySlice = append([]byte{byte(nibble)}, branch.KeySlice...)
								node.Value = branch.Value
								node.Branches = nil
								return node, nil
							}
							if branch.Flag == Extension {
								node.Flag = Extension
								node.Next = branch.Next
								node.Index = [16]int{}
								node.KeySlice = append([]byte{byte(nibble)}, branch.KeySlice...)
								node.Branches = nil
								return node, nil
							}
							if branch.Flag == Branch {
								node.Flag = Extension
								node.Next = branch
								node.Index = [16]int{}
								node.KeySlice = []byte{byte(nibble)}
								node.Branches = nil
								return node, nil
							}
						}
					}
					if length != 1 {
						return node, nil
					}
				}
				if branch.Flag == Extension {
					var err error
					node.Branches[node.Index[int(nibble)]], err = branch.Delete(key[1:])
					if err != nil {
						return nil, err
					}
					return node, nil
				}
				if branch.Flag == Branch {
					var err error
					node.Branches[node.Index[int(nibble)]], err = branch.Delete(key[1:])
					if err != nil {
						return nil, err
					}
					return node, nil
				}
			}

		}
	}
	return nil, fmt.Errorf("delete error")
}

// func (node *Node) Serialize() []byte {
// 	var encoded bytes.Buffer
// 	encode := gob.NewEncoder(&encoded)
// 	err := encode.Encode(node)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return encoded.Bytes()
// }

// func (node *Node) Hash() []byte {
// 	hash := sha256.Sum256(node.Serialize())
// 	return hash[:]
// }
