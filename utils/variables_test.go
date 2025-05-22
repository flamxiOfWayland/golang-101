package utils

import (
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, IString, " ", "string comparison")
}

func TestInt(t *testing.T) {
	assert.Equal(t, IInt, IUint, "int vs uint comparison")
}

func TestFloat(t *testing.T) {
	assert.Equal(t, IFloat32, IFloat64, "float32 vs float64 comparison")
}

// TODO: Homework
func TestComplex(t *testing.T) {
	assert.Equal(t, IComplex64, IComplex128, "complex numbers comparison")
}

func TestIntPointer(t *testing.T) {
	assert.Equal(t, IPInt, IInt)
	// IInt = 123
	// assert.Equal(t, IPInt, IInt)
	// assert.Equal(t, IPInt, 123)
}

func TestIArrayInt(t *testing.T) {
	assert.Equal(t, IArrayInt, []int{1, 2, 3, 4})
}

// TODO: Homework
func TestIArraySlice(t *testing.T) {
	assert.Equal(t, ISliceInt, []int{1, 2, 3, 4})
	ISliceInt = append(ISliceInt, -1)
	ISliceInt[len(ISliceInt)-1] = -2
	sort.Slice(ISliceInt, func(i, j int) bool {
		return i > j
	})
	// assert.Equal(t, ISliceInt, []int{-2, 1, 2, 3, 4})
	// assert.Equal(t, ISliceInt, []int{-2, 4, 3, 2, 1})
	// assert.Equal(t, ISliceInt, []int{4, 3, 2, 1, -2})
	ISliceInt = append(ISliceInt, ISliceInt...)
}

func TestIMapIntString(t *testing.T) {
	assert.Equal(t, IMapIntString[1], "one")
	IMapIntString[123] = "123"
	for key, value := range IMapIntString {
		log.Printf("key: %d, value: %s\n", key, value)
	}
}
