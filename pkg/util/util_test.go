package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	FirstName string `json:"first_name" sensitive:"true"`
	LastName  string
	Address   Address `json:"address"`
	T         [1]string
}

type Address struct {
	City string `json:"city"`
	Pr   string `json:"pr"  sensitive:"true"`
}

func Test_Clone(t *testing.T) {

	user := &User{
		FirstName: "Quy",
		LastName:  "Kim",
		Address: Address{
			City: "Ha noi",
		},
	}
	var t1 [1]string
	t1[0] = "Hello"
	user.T = t1
	u2 := Clone(user)
	u3, ok := u2.(User)
	if ok {
		u3.LastName = "No"
	}
	user.LastName = "Test"

	js, _ := json.Marshal(u2)
	fmt.Println("new User: ", string(js))
}

func Test_Mask(t *testing.T) {
	user := User{
		FirstName: "Quy",
		LastName:  "Kim",
		Address: Address{
			City: "Ha noi",
			Pr:   "truong lam",
		},
	}
	Mask(&user)
	js, _ := json.Marshal(&user)
	fmt.Println("mask User: ", string(js))
}
