package db

import (
	"fmt"
	"testing"
)

func Test_GetAll(t *testing.T) {

}

func Test_Create(t *testing.T) {

}

func Test_Update(t *testing.T) {

}

func Test_Delete(t *testing.T) {

}

func Test_Filter(t *testing.T) {

}

func Test_MakeFilter(t *testing.T) {
	var f []F
	f = append(f, F{Key: "name", Value: "Wuy"})
	f = append(f, F{Key: "age", Value: 36})
	ft := Filter{f}
	m, err := makeFilter(ft)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("m: ", m)

}
