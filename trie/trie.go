package trie

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
)

type Trie struct {
	Root *Node
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Hash() []byte {
	hash := sha256.Sum256(t.Encode())
	return hash[:]
}

func (t *Trie) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeStateTree(to_decode []byte) *Trie {
	var tree Trie

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tree)
	if err != nil {
		log.Panic(err)
	}

	return &tree
}

func (t *Trie) Get(key []byte) ([]byte, bool) {
	node := t.Root
	nibbles := FromBytes(key)
	for {
		if IsEmptyNode(node) {
			return nil, false
		}

		if node.Flag == Leaf {
			matched := PrefixMatchedLen(node.KeySlice, nibbles)
			if matched != len(node.KeySlice) || matched != len(nibbles) {
				return nil, false
			}
			return node.Value, true
		}

		if node.Flag == Branch {
			if len(nibbles) == 0 {
				if yes, err := node.HasValue(); err != nil {
					panic(err)
				} else {
					return node.Value, yes
				}
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			if node.Index[int(b)] == -1 {
				return nil, false
			}
			node = node.Branches[node.Index[int(b)]]
			continue
		}

		if node.Flag == Extension {
			matched := PrefixMatchedLen(node.KeySlice, nibbles)
			// E 01020304
			//   010203
			if matched < len(node.KeySlice) {
				return nil, false
			}

			nibbles = nibbles[matched:]
			node = node.Next
			continue
		}

		panic("not found")
	}
}

// Put adds a key value pair to the Trie
// In general, the rule is:
// - When stopped at an EmptyNode, replace it with a new LeafNode with the remaining KeySlice.
// - When stopped at a LeafNode, convert it to an ExtensionNode and add a new branch and a new LeafNode.
// - When stopped at an ExtensionNode, convert it to another ExtensionNode with shorter KeySlice and create a new BranchNode points to the ExtensionNode.
func (t *Trie) Put(key []byte, value []byte) {
	// need to use pointer, so that I can update Root in place without
	// keeping trace of the parent node
	nibbles := FromBytes(key)
	if IsEmptyNode(t.Root) {
		t.Root = NewLeafNodeFromNibbles(nibbles, value)
		return
	}
	node := t.Root
	for {
		if node.Flag == Leaf {
			matched := PrefixMatchedLen(node.KeySlice, nibbles)

			// if all matched, update value even if the value are equal
			if matched == len(nibbles) && matched == len(node.KeySlice) {
				*node = *NewLeafNodeFromNibbles(node.KeySlice, value)
				return
			}

			branch := NewBranchNode()
			// if matched some nibbles, check if matches either all remaining nibbles
			// or all leaf nibbles
			if matched == len(node.KeySlice) {
				if err := branch.SetValue(node.Value); err != nil {
					panic(err)
				}
			}

			if matched == len(nibbles) {
				if err := branch.SetValue(value); err != nil {
					panic(err)
				}
			}

			if matched < len(node.KeySlice) {
				// have dismatched
				// L 01020304 hello
				// + 010203   world

				// 01020304, 0, 4
				branchNibble, leafNibbles := node.KeySlice[matched], node.KeySlice[matched+1:]
				newLeaf := NewLeafNodeFromNibbles(leafNibbles, node.Value) // not :matched+1
				if err := branch.SetBranch(branchNibble, newLeaf); err != nil {
					panic(err)
				}
			}

			if matched < len(nibbles) {
				// L 01020304 hello
				// + 010203040 world

				// L 01020304 hello
				// + 010203040506 world
				branchNibble, leafNibbles := nibbles[matched], nibbles[matched+1:]
				newLeaf := NewLeafNodeFromNibbles(leafNibbles, value)
				if err := branch.SetBranch(branchNibble, newLeaf); err != nil {
					panic(err)
				}
			}

			// if there is matched nibbles, an extension node will be created
			if matched > 0 {
				// create an extension node for the shared nibbles
				ext := NewExtensionNode(node.KeySlice[:matched], branch)
				*node = *ext
			} else {
				// when there no matched nibble, there is no need to keep the extension node
				*node = *branch
			}

			return
		}

		if node.Flag == Branch {
			if len(nibbles) == 0 {
				if err := node.SetValue(value); err != nil {
					panic(err)
				}
				return
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			if node.Index[int(b)] == -1 {
				node.Branches = append(node.Branches, NewLeafNodeFromNibbles(nibbles, value))
				node.Index[int(b)] = len(node.Branches) - 1
				return
			}
			node = node.Branches[node.Index[int(b)]]
			continue
		}

		// E 01020304
		// B 0 hello
		// L 506 world
		// + 010203 good
		if node.Flag == Extension {
			matched := PrefixMatchedLen(node.KeySlice, nibbles)
			if matched < len(node.KeySlice) {
				// E 01020304
				// + 010203 good
				extNibbles, branchNibble, extRemainingnibbles := node.KeySlice[:matched], node.KeySlice[matched], node.KeySlice[matched+1:]
				nodeBranchNibble, nodeLeafNibbles := nibbles[matched], nibbles[matched+1:]
				branch := NewBranchNode()
				if len(extRemainingnibbles) == 0 {
					// E 0102030
					// + 010203 good
					if err := branch.SetBranch(branchNibble, node.Next); err != nil {
						panic(err)
					}
				} else {
					// E 01020304
					// + 010203 good
					newExt := NewExtensionNode(extRemainingnibbles, node.Next)
					if err := branch.SetBranch(branchNibble, newExt); err != nil {
						panic(err)
					}
				}

				remainingLeaf := NewLeafNodeFromNibbles(nodeLeafNibbles, value)
				if err := branch.SetBranch(nodeBranchNibble, remainingLeaf); err != nil {
					panic(err)
				}

				// if there is no shared extension nibbles any more, then we don't need the extension node
				// any more
				// E 01020304
				// + 1234 good
				if len(extNibbles) == 0 {
					*node = *branch
				} else {
					// otherwise create a new extension node
					*node = *NewExtensionNode(extNibbles, branch)
				}
				return
			}

			nibbles = nibbles[matched:]
			node = node.Next
			continue
		}

		panic("unknown type")
	}

}

func (t *Trie) Delete(address []byte) error {
	if _, ok := t.Get(address); !ok {
		return fmt.Errorf("%v not found", address)
	}
	key := FromBytes(address)
	var err error
	t.Root, err = t.Root.Delete(key)
	return err
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewTrieWithData(data [][]byte) *Trie {
	trie := NewTrie()

	for id, datum_bytes := range data {
		str_id := strconv.Itoa(id)
		trie.Put(datum_bytes, []byte(str_id))
	}

	return trie
}

func (t *Trie) PrintState() {
	node := t.Root
	if node == nil {
		fmt.Println("It's an empty tree!")
	}else {
		key := []byte{}
		node.PrintState(key)
	}	
}

func DeepCopy(dst, src *Trie) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(dst)
}