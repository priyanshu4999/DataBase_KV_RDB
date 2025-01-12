import (
	"binary"
)
type Btree struct {
	data []byte 
}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2 
)

type BTree struct{
	root uint64
	get func(uint64) BNode
	new func(BNode) unit64
	del func(uint64) 
}

const HEADER = 4
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init(){
	nodelmax := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_KEY_SIZE 
	assert(nodelmax <= BTREE_PAGE_SIZE)
}
func (node BNode) btype() uint16{return binary.LittleEndian.Uint16}
func (node BNode) nkeys() uint16{return binary.LittleEndian.Uint16}
func (node BNode) setHeader(btype , nkeys int){binary.LittleEndian}
func (node BNode) getPtr(idx uint16)unint64{
	asset
}