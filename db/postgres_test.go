package db

import (
	"fmt"
	"testing"
)

func TestSetDSN(t *testing.T) {
	dsn, err := setDSN()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dsn)
}
