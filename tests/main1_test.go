package tests

import "testing"

func TestAdd2(t *testing.T) {
	n := Add(1, 2)
	if n == 3 {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}

func TestAdd3(t *testing.T) {
	n := Add(1, 2)
	if n == 3 {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}
