package mpt

// 简单的merkle树

type Node []byte

type MerkleTree struct {
	nodes []Node
}

func (mt *MerkleTree) AddNode(n Node) {
	mt.nodes = append(mt.nodes, n)
}

func (mt *MerkleTree) Root() {

}
