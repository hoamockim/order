package util

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func Test_Retry(t *testing.T) {
	count := 1
	if err := Delay(1 * time.Second).Run(func() error {
		var err error
		if count < 3 {
			count++
			err = errors.New("cannot run function at this time")
		} else {
			fmt.Println("Hello world")
			err = errors.New("done")
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
}