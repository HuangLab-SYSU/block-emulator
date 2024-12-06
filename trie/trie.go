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

	enc := gob.NewEncoder(&buff) // 不能有空指针，变量首字母需大写
	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeStateTree(to_decode []byte) *Trie {
	var tree Trie

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tree) // 如何decode出tree？
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
			matched := PrefixMatchedLen(node.KeySlice, nibbles) // key words 1
			if matched != len(node.KeySlice) || matched != len(nibbles) {
				return nil, false
			}
			return node.Value, true
		}

		if node.Flag == Branch {
			if len(nibbles) == 0 { // 若有元素在该分支节点终止，返回分支节点value值
				if yes, err := node.HasValue(); err != nil {
					panic(err)
				} else {
					return node.Value, yes // 终止地址对应balance
				}
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			if node.Index[int(b)] == -1 {
				return nil, false
			}
			node = node.Branches[node.Index[int(b)]] // node.Index[int(b)] = len(node.Branches) - 1
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

// 加入节点 or 更新节点value
// Put adds a key value pair to the Trie
// In general the rule is:
// - When stopped at an EmptyNode, replace it with a new LeafNode with the remaining KeySlice.
// - When stopped at a LeafNode, convert it to an ExtensionNode and add a new branch and a new LeafNode.
// - When stopped at an ExtensionNode, convert it to another ExtensionNode with shorter KeySlice and create a new BranchNode points to the ExtensionNode.
func (t *Trie) Put(key []byte, value []byte) {
	// need to use pointer, so that I can update Root in place without
	// keeping trace of the parent node
	// key: account, value: state
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
			// 1、完全匹配原leaf_node，直接更新value值
			if matched == len(nibbles) && matched == len(node.KeySlice) {
				*node = *NewLeafNodeFromNibbles(node.KeySlice, value)
				return
			}

			// 2、处理新加入的branch_node及下面的leaf_node
			// 若有其中一个完全匹配，则作为分支节点的value
			branch := NewBranchNode()
			// if matched some nibbles, check if matches either all remaining nibbles
			// or all leaf nibbles
			// 新节点长度大于原叶子节点，原叶子节点变为分支节点
			if matched == len(node.KeySlice) {
				if err := branch.SetValue(node.Value); err != nil {
					panic(err)
				}
			}

			// 新节点长度小于原叶子节点，新节点变为分支节点
			if matched == len(nibbles) {
				if err := branch.SetValue(value); err != nil {
					panic(err)
				}
			}

			// 原叶子节点部分匹配，123456，12378，branchNibble=4，leafNibbles=5
			if matched < len(node.KeySlice) {
				// have dismatched
				// L 01020345 hello
				// + 010203   world

				// 01020345, 4, 5
				branchNibble, leafNibbles := node.KeySlice[matched], node.KeySlice[matched+1:]
				newLeaf := NewLeafNodeFromNibbles(leafNibbles, node.Value) // not :matched+1
				if err := branch.SetBranch(branchNibble, newLeaf); err != nil {
					panic(err)
				}
			}

			// 新节点部分匹配，123456，12378，branchNibble=4，leafNibbles=5
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

			// 3、若有匹配位，则增加extension_node
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
			// 若在branch停下，保存value
			if len(nibbles) == 0 {
				if err := node.SetValue(value); err != nil {
					panic(err)
				}
				return
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			// 若此分支下无节点，直接添加一个leaf_node
			if node.Index[int(b)] == -1 {
				node.Branches = append(node.Branches, NewLeafNodeFromNibbles(nibbles, value))
				node.Index[int(b)] = len(node.Branches) - 1
				return
			}
			// 否则向下移动到子节点，下一轮for循环中对应地处理
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
				// 处理原状态树结构
				if len(extRemainingnibbles) == 0 {
					// 在原extension node 和 branch 间多一个new_branch
					// new_branch的node.KeySlice[matched]指向原branch，原extension node减少一位
					// E 0102030
					// + 010203 good
					if err := branch.SetBranch(branchNibble, node.Next); err != nil {
						panic(err)
					}
				} else {
					// 在原extension node 和 branch 间多一个new_branch和new_extension_node
					// new_branch的node.KeySlice[matched]指向new_extension node, new_extention_node指向原branch
					// E 01020304
					// + 010203 good
					newExt := NewExtensionNode(extRemainingnibbles, node.Next)
					if err := branch.SetBranch(branchNibble, newExt); err != nil {
						panic(err)
					}
				}

				// 处理新加入节点，为 new_branch指向的new_leaf_node
				remainingLeaf := NewLeafNodeFromNibbles(nodeLeafNibbles, value)
				if err := branch.SetBranch(nodeBranchNibble, remainingLeaf); err != nil {
					panic(err)
				}

				// 处理原extension node
				// if there is no shared extension nibbles any more, then we don't need the extension node any more
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

func (t *Trie) Delete(address string) error {
	if _, ok := t.Get([]byte(address)); !ok {
		return fmt.Errorf("%v not found", address)
	}
	key := FromString(address)
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
	} else {
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
