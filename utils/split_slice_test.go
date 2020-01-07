package utils

import (
	"fmt"
	"testing"
)

func Test_split_slice(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := NewSplitIter(len(arr), 3)
	for {
		start, end, hasNext := iter.Next()
		fmt.Println(arr[start:end])
		if !hasNext {
			break
		}
	}
}
