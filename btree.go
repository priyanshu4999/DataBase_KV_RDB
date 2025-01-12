package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

type BNode struct {
	data []byte
}
type BnodeInterface interface {
	btype() uint16
	nkeys() uint16
	setHeader(datatype uint16, numberofkeys uint16)

	getPtr(pos uint16) uint64
	setPtr(pos uint16, value uint64)

	offsetPos(pos uint16) uint16 //helper to getOffset and setOffset
	getOffset(pos uint16) uint16
	setOffset(pos uint16, value uint16)

	kvPos(pos uint16) uint16 //helper to getKey and setKey
	getKey(pos uint16) []byte
	getVal(pos uint16 )[] byte

}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

type BTree struct {
	root uint64
	get  func(uint64) BNode
	new  func(BNode) uint64
	del  func(uint64)
}

const HEADER = 4
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

// SOME general rules:
// 	To access 2byte values binary littleendian uint16 function is used
//  To access 8byte values binary littleendian uint64 function is used
//  Slice range indexing(slicing) is avoided for fixed length values
// For variable length storage range indexing cannot be avoided

func assert(condition bool, message string) {
	if !condition {
		log.Fatal(message)
		fmt.Println(message)
	}
}
func init() {
	nodelmax := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_KEY_SIZE

	assert(nodelmax <= BTREE_PAGE_SIZE, "nodelength greater than BTREE PAGE")
}
func (node BNode) btype() uint16 { return binary.LittleEndian.Uint16(node.data[0:2]) }
func (node BNode) nkeys() uint16 { return binary.LittleEndian.Uint16(node.data[2:4]) }
func (node BNode) setHeader(btype, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// ////////// POINTER RELATED
func (node BNode) getPtr(idx uint16) uint64 { /// FIRST METHID TO GET KEY VALUE PAIR
	assert(node.nkeys() > idx, "index more than keys length")
	ptr := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[ptr:])
}
func (node BNode) setPtr(idx uint16, value uint64) {
	assert(node.nkeys() > idx, "index more than keys length")
	ptr := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[ptr:], value)
}

// ///////// OFFSET RELATED
func (node BNode) offsetPos(idx uint16) uint16 {
	assert(node.nkeys() >= idx && idx >= 1, "REQUESTED OFFSET OUT OF BOUNDS")
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 { /// SECOND METHOD TO GET KEY VALUE PAIR
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[node.offsetPos(idx):])
}

func (node BNode) setOffset(idx, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[node.offsetPos(idx):], offset)
}

// ////////// Key Value Pair
// || (HEADER) TYPE 2B | nKEYS 2B || (POINTERS) PTR1 8B | PTR2 8B ...EACH 8B || (OFFSETS)  offset1 2B offset2 2B ....||
// ||(KV-PAIRS) klen1 2B vlen1 2B | key1 (klen) val1 (vlen) ..... ||
// || type , nKeys | 8*nkeys | 2*nkeys | (offset1 = 0)(offset2 = 4 + klen1 + vlen1) .......   |
// e.g. | 2bxINTEGER | 2bx256  || 8bx65536 ..... (256th pointer) || 2bx001 ..... (256th) ||
// /
func (node BNode) kvPos(idx uint16) uint16 {
	//kv pos = getOffsetKv(@offsetPos(n-1thOffset))
	assert(idx <= node.nkeys(), "INDEX GREATER THAN NKEYS")
	// return // total-header + total-pointers + total-offset + value-of-nth-offset
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	assert(idx <= node.nkeys(), "INDEX GREATER THAN NKEYS")
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}
func (node BNode) getVal(idx uint16) []byte {
	assert(idx <= node.nkeys(), "INDEX GREATER THAN NKEYS")
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}
