package binarySearchTree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Examples

This example will create a list of employees in two ways:
1) default searchable
2) searchable by employee ID
3) custom searchable by full name.
*/

type employee struct {
	first      string
	last       string
	employeeID uint32
}

var employeeData = []employee{
	{
		first:      "john",
		last:       "doe",
		employeeID: 0,
	},
	{
		first:      "jane",
		last:       "doe",
		employeeID: 1,
	},
	{
		first:      "linus",
		last:       "torvalds",
		employeeID: 2,
	},
	{
		first:      "alan",
		last:       "turing",
		employeeID: 3,
	},
	{
		first:      "bill",
		last:       "gates",
		employeeID: 4,
	},
}

func ExampleF_withDefaultSearchable() {
	bst := NewBST()
	keys := make([]uint, len(employeeData))
	for i, employee := range employeeData {
		leaf := NewLeaf(employee)
		keys[i] = leaf.Key()
		bst.Insert(leaf)
	}
	john := bst.FindByKey(keys[0])
	fmt.Println(john.Value.(employee).employeeID)

	// Output: 0
}

func ExampleF_withSearchable() {
	bst := NewBST(WithSearchable("employeeID"))
	for _, employee := range employeeData {
		bst.Insert(NewLeaf(employee))
	}
	found := bst.FindByValue(4)

	fmt.Println(found.Value.(employee).first)

	// Output: bill
}

func ExampleF_withCustomSearchable() {
	fn := func(s any) any {
		fullName := s.(employee).first + " " + s.(employee).last
		return fullName
	}
	bst := NewBST(WithCustomSearchFn(fn))
	for _, employee := range employeeData {
		bst.Insert(NewLeaf(employee))
	}
	found := bst.FindByValue("john doe")

	fmt.Println(found.Value.(employee).employeeID)

	// Output: 0
}

/*
End Examples
*/

func setUp() BinarySearchTree {
	bst := NewBST()
	leafs := []struct {
		value string
		key   uint
	}{{"foo", 3}, {"bar", 2}, {"baz", 5}, {"qux", 4}}

	for _, leaf := range leafs {
		bst.Insert(&Leaf{key: leaf.key, Value: leaf.value})
	}
	return bst
}
func TestNewLeaf(t *testing.T) {
	l := NewLeaf("foo")
	assert.IsType(t, &Leaf{}, l)
}

func TestSearchableTypes(t *testing.T) {
	bst := NewBST(WithSearchable("num"))
	bst.Insert(NewLeaf(struct{ num int32 }{1}))
	assert.NotNil(t, bst.FindByValue(1))
	bst.Insert(NewLeaf(struct{ num int64 }{2}))
	assert.NotNil(t, bst.FindByValue(2))
	bst.Insert(NewLeaf(struct{ num uint }{3}))
	assert.NotNil(t, bst.FindByValue(3))
	bst.Insert(NewLeaf(struct{ num uint64 }{4}))
	assert.NotNil(t, bst.FindByValue(4))
	bst.Insert(NewLeaf(struct{ num float64 }{4.0}))
	assert.NotNil(t, bst.FindByValue(struct{ num float64 }{4.0}))
}
func TestWithSearchableInt(t *testing.T) {
	type testCase struct {
		name string
		age  int
	}
	values := []testCase{
		{name: "xxxx", age: 0},
		{name: "xx", age: 5},
		{name: "xxx", age: 1},
		{name: "x", age: 2},
		{name: "xxxxx", age: 4},
	}
	bst := NewBST(WithSearchable("age"))
	for _, val := range values {
		bst.Insert(NewLeaf(val))
	}
	found := bst.FindByValue(2)
	assert.Equal(t, "x", found.Value.(testCase).name)
}

func TestWithSearchableString(t *testing.T) {
	type testCase struct {
		name string
		age  int
	}
	values := []testCase{
		{name: "xxxx", age: 0},
		{name: "xx", age: 5},
		{name: "xxx", age: 1},
		{name: "x", age: 2},
		{name: "xxxxx", age: 4},
	}
	bst := NewBST(WithSearchable("name"))
	for _, val := range values {
		bst.Insert(NewLeaf(val))
	}
	found := bst.FindByValue("x")
	assert.Equal(t, "x", found.Value.(testCase).name)
}

func TestNewBSTwithCustomFn(t *testing.T) {
	type testCase struct {
		name string
		age  int
	}
	values := []testCase{
		{name: "xxxx", age: 0},
		{name: "xx", age: 0},
		{name: "xxx", age: 0},
		{name: "x", age: 0},
		{name: "xxxxx", age: 0},
	}
	valueFn := func(s any) any {
		return s.(testCase).name + "s"
	}
	bst := NewBST(WithCustomSearchFn(valueFn))
	for _, val := range values {
		bst.Insert(NewLeaf(val))
	}
	found := bst.FindByValue("xxxs")
	assert.Equal(t, 0, found.Value.(testCase).age)
}

func TestInsert(t *testing.T) {
	bst := NewBST(WithSearchable("value"))
	leafVal := struct {
		value string
		key   uint
	}{"dup", 3}
	leaf := NewLeaf(leafVal)
	err := bst.Insert(leaf)
	assert.NoError(t, err)
	err = bst.Insert(leaf)
	assert.ErrorIs(t, err, DuplicateLeafError)
}

func TestDuplicateKeys(t *testing.T) {
	bst := NewBST()
	leaf1 := &Leaf{key: 1, Value: "foo1"}
	leaf2 := &Leaf{key: 1, Value: "foo2"}
	err := bst.Insert(leaf1)
	assert.NoError(t, err)
	err = bst.Insert(leaf2)
	assert.ErrorIs(t, err, DuplicateLeafError)
}

func TestFindByKey(t *testing.T) {
	bst := setUp()
	assert.Equal(t, "bar", bst.FindByKey(2).Value)
	assert.Equal(t, "foo", bst.FindByKey(3).Value)
	assert.Equal(t, "qux", bst.FindByKey(4).Value)
	assert.Nil(t, bst.FindByKey(0))
}

func TestFindByValue(t *testing.T) {
	bst := setUp()
	assert.Equal(t, uint(2), bst.FindByValue("bar").Key())
	assert.Equal(t, uint(3), bst.FindByValue("foo").Key())
	assert.Equal(t, uint(4), bst.FindByValue("qux").Key())
	assert.Nil(t, bst.FindByValue("nope"))
}

func TestRemove(t *testing.T) {
	bst := setUp()
	assert.NotEmpty(t, bst)
}

func BenchmarkInsertAndNewLeaf(b *testing.B) {
	bst := NewBST()
	for i := 0; i < b.N; i++ {
		bst.Insert(NewLeaf(i))
	}
}

func BenchmarkFindByValue(b *testing.B) {
	bst := NewBST()
	for i := 0; i < 50_000; i++ {
		bst.Insert(NewLeaf(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.FindByValue(i)
	}
}
