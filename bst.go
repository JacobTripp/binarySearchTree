// Package binarySearchTree implements a BST with the ability to set an
// arbitrary struct field as a uniq key.
//
// This implementation of a BST data structure as one extra ability,
// you can search by the auto-generated leaf key or by an arbitrary field
// of the struct you store in a leaf's value. Right now the field must be a
// string or an int. The leaf added to a tree must be unique by the auto-
// generated key and by the defined searchable field value.
//
// So far there are no traversing methods provided since this is intended to
// be a search and store only type of structure. Perhaps traversal methods will
// be added in future releases.
//
// This should create a balanced tree no matter if you pass in sorted values
// since it randomly generates leaf keys when leaves are created.
//
// See the examples in the bst_test.go file.
package binarySearchTree

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/google/uuid"
)

// Leaf is the basic node of the BST, the key is the auto-generated
// uuid.ID used for leaf placement
type Leaf struct {
	key   uint
	Value any
	left  *Leaf
	right *Leaf
}

// Make a new leaf pointer with the id initalized
func NewLeaf(v any) *Leaf {
	return &Leaf{
		key:   uint(uuid.New().ID()),
		Value: v,
	}
}

// the key of the leave isn't exported so it cannot be set arbitraritly set
// by developers using this package. But since you can search by a leaf key
// this provide access to it's value.
func (l Leaf) Key() uint {
	return l.key
}

// in order to have struct fields searchable it maintains a map of the
// field and it's associated key, it then uses that key for seraching.
type BinarySearchTree struct {
	root    *Leaf
	leafMap map[any]uint // at some point this should have better typing
	mapfn   KeyFn
	lock    sync.RWMutex
}

type KeyFn func(any) any // This needs more specific types as well

type bstOpt func(*BinarySearchTree)

// if you want a more customized key function you can provide one with this
// option.
func WithCustomSearchFn(fn KeyFn) bstOpt {
	return func(bst *BinarySearchTree) {
		bst.mapfn = fn
	}
}

// An easy helper option so you just need to provide the name of the struct
// field you want to set as a unique key. The value of the key must be a string
// or an int
func WithSearchable(attributeName string) bstOpt {
	fn := func(s any) any {
		val := reflect.ValueOf(s)
		name := val.FieldByName(attributeName)
		if name.Type().String() == "string" {
			return name.String()
		}
		switch name.Type().String() {
		case "string":
			return name.String()
		case "int":
			return int(name.Int())
		case "int32":
			return int(name.Int())
		case "int64":
			return int(name.Int())
		case "uint":
			return int(name.Int())
		case "uint64":
			return int(name.Int())
		case "uint32":
			return int(name.Uint()) // maps don't like anything but ints
		}
		return s
	}
	return func(bst *BinarySearchTree) {
		bst.mapfn = fn
	}
}

// Start a new tree.
func NewBST(opts ...bstOpt) BinarySearchTree {
	bst := BinarySearchTree{
		leafMap: map[any]uint{},
		mapfn:   func(v any) any { return v },
	}
	for _, opt := range opts {
		opt(&bst)
	}
	return bst
}

// duplicates are not allowed
var DuplicateLeafError = errors.New("Duplicate leaf")

// Insert a new leaf into the tree.
// The choice to have it accept only the Leaf type instead of any is because
// I want the developer to be aware of duplicates and to keep track of the
// leaf keys.
func (bst *BinarySearchTree) Insert(leaf *Leaf) error {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	_, found := bst.leafMap[bst.mapfn(leaf.Value)]
	if found {
		return fmt.Errorf(
			"%w: '%v' value is already in the tree",
			DuplicateLeafError,
			leaf.Value,
		)
	}
	bst.leafMap[bst.mapfn(leaf.Value)] = leaf.key
	if bst.root == nil {
		bst.root = leaf
	} else {
		return insertLeaf(bst.root, leaf)
	}
	return nil
}

// The meat of the insert, standard BST algo where left is less than and right
// is greater than.
func insertLeaf(leaf, toInsert *Leaf) error {
	if toInsert.key == leaf.key {
		return fmt.Errorf(
			"%w: the leaf key '%d' already exists",
			DuplicateLeafError,
			leaf.key,
		)
	}
	if toInsert.key < leaf.key {
		if leaf.left == nil {
			leaf.left = toInsert
		} else {
			return insertLeaf(leaf.left, toInsert)
		}
	} else {
		if leaf.right == nil {
			leaf.right = toInsert
		} else {
			return insertLeaf(leaf.right, toInsert)
		}
	}
	return nil
}

// Given a key, return a leaf or nil if it doesn't exists
func (bst BinarySearchTree) FindByKey(key uint) *Leaf {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	if bst.root.key == key {
		return bst.root
	}
	if key < bst.root.key {
		return findByKey(bst.root.left, key)
	}
	if key > bst.root.key {
		return findByKey(bst.root.right, key)
	}
	return nil
}

func findByKey(leaf *Leaf, key uint) *Leaf {
	if leaf == nil {
		return nil
	}
	if key == leaf.key {
		return leaf
	}
	if key < leaf.key {
		return findByKey(leaf.left, key)
	}
	if key > leaf.key {
		return findByKey(leaf.right, key)
	}
	return nil
}

// Find a leaf by the defined searchable or using the default of just the
// value.
// First it gets the key from the leaf map and the searches by key.
func (bst BinarySearchTree) FindByValue(v any) *Leaf {
	key, ok := bst.leafMap[v]
	if !ok {
		return nil
	}
	return findByKey(bst.root, key)
}
