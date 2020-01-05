package nonull

import (
	"encoding/json"
	"fmt"
	"testing"
)

type T struct {
	Name  string
	Name2 *string
	Arr   []int
	Arr2  *[]int
	Map1  map[string]int
	Map2  *map[string]int
	Sub1  T0
	Sub2  *T0
}

type T0 struct {
	_Name string
	Name3 string
	_Arr  []int
	Arr3  []int
}

func Test_a(t *testing.T) {
	var obj = T{
		Name: "",
		Arr:  nil,
	}
	d, _ := json.Marshal(obj)
	fmt.Printf("%s\n", d)
	Make(&obj)
	d, _ = json.Marshal(obj)
	fmt.Printf("%s\n", d)
}
