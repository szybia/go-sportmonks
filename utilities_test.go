package sportmonks

import "testing"

func TestIntSliceToSepString(t *testing.T) {
	ret := IntSliceToSepString([]int{}, "")
	if ret != "" {
		t.Errorf("IntSliceToSepString of []int {} with \"\" separator was incorrect. Got: \"%v\". Want: \"\"", ret)
	}

	ret = IntSliceToSepString([]int{}, ",")
	if ret != "" {
		t.Errorf("IntSliceToSepString of []int {} with \",\" separator was incorrect. Got: \"%v\". Want: \"\"", ret)
	}

	ret = IntSliceToSepString([]int{1, 2, 3}, "")
	if ret != "123" {
		t.Errorf("IntSliceToSepString of []int {1, 2, 3} with \"\" separator was incorrect. Got: \"%v\". Want: \"%v\"", ret, "123")
	}

	ret = IntSliceToSepString([]int{1, 2, 3}, ",")
	if ret != "1,2,3" {
		t.Errorf("IntSliceToSepString of []int {1, 2, 3} with \",\" separator was incorrect. Got: \"%v\". Want: \"%v\"", ret, "1,2,3")
	}

	ret = IntSliceToSepString([]int{1, 2, 3}, "-")
	if ret != "1-2-3" {
		t.Errorf("IntSliceToSepString of []int {1, 2, 3} with \"-\" separator was incorrect. Got: \"%v\". Want: \"%v\"", ret, "1-2-3")
	}
}
