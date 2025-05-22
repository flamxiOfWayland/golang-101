package utils

import "fmt"

func IfFoo(a int, b, c float32) bool {
	if a > 0 {
		return true
	} else if b > 0 {
		return true
	} else if c > 0 {
		return true
	}

	return false
}

func SwtichFoo(a int, b, c float32) bool {
	switch result := a > 0 && b > 0 && c > 0; result {
	case true:
		return true
	case false:
		return false
	default:
		fmt.Println("WTF")
	}
	return false
}

func CreateSlice(s int) ([]int, error) {
	if s <= 0 {
		return nil, fmt.Errorf("invalid size request")
	}
	data := make([]int, s)
	for idx := range data {
		data[idx] = 0
	}

	return data, nil
}
