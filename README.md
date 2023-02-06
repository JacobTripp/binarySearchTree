![Binary Search Tree Logo](/assets/bst_logo.svg)
Binary Search Tree
==================

If you use this package in the next five minutes I'll include a free toaster!\*

This is more of a fun exercise in playing with a data structure that I normally
don't mess with in such a manual way. This is a pretty standard BST but so far
it doesn't implement and traversing methods. It's just an insert/search
structure.  The only unique aspect of this implementation is that you can set it
to be unique based on either a string/int type field in a struct, or you can
define a custom function to determine uniqueness.

<sub>*\*offer only valid to time travelers, must provide proof of time travel by
submitting future jackpot winning Powerball numbers for a drawing occurring
within 3 months after submission date.*</sub>

## Go Doc

package binarySearchTree // import "github.com/JacobTripp/binarySearchTree"

Package binarySearchTree implements a BST with the ability to set an
arbitrary struct field as a uniq key.

This implementation of a BST data structure as one extra ability, you can
search by the auto-generated leaf key or by an arbitrary field of the struct
you store in a leaf's value. Right now the field must be a string or an int.
The leaf added to a tree must be unique by the auto- generated key and by
the defined searchable field value.

So far there are no traversing methods provided since this is intended to be
a search and store only type of structure. Perhaps traversal methods will be
added in future releases.

This should create a balanced tree no matter if you pass in sorted values
since it randomly generates leaf keys when leaves are created.

See the examples in the bst_test.go file.

VARIABLES

var DuplicateLeafError = errors.New("Duplicate leaf")
    duplicates are not allowed


FUNCTIONS

func WithCustomSearchFn(fn KeyFn) bstOpt
    if you want a more customized key function you can provide one with this
    option.

func WithSearchable(attributeName string) bstOpt
    An easy helper option so you just need to provide the name of the struct
    field you want to set as a unique key. The value of the key must be a string
    or an int


TYPES

type BinarySearchTree struct {
	// Has unexported fields.
}
    in order to have struct fields searchable it maintains a map of the field
    and it's associated key, it then uses that key for seraching.

func NewBST(opts ...bstOpt) BinarySearchTree
    Start a new tree.

func (bst BinarySearchTree) FindByKey(key uint) *Leaf
    Given a key, return a leaf or nil if it doesn't exists

func (bst BinarySearchTree) FindByValue(v any) *Leaf
    Find a leaf by the defined searchable or using the default of just the
    value. First it gets the key from the leaf map and the searches by key.

func (bst *BinarySearchTree) Insert(leaf *Leaf) error
    Insert a new leaf into the tree. The choice to have it accept only the Leaf
    type instead of any is because I want the developer to be aware of
    duplicates and to keep track of the leaf keys.

type KeyFn func(any) any // This needs more specific types as well

type Leaf struct {
	Value any

	// Has unexported fields.
}
    Leaf is the basic node of the BST, the key is the auto-generated uuid.ID
    used for leaf placement

func NewLeaf(v any) *Leaf
    Make a new leaf pointer with the id initalized

func (l Leaf) Key() uint
    the key of the leave isn't exported so it cannot be set arbitraritly set by
    developers using this package. But since you can search by a leaf key this
    provide access to it's value.

