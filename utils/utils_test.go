package utils

import (
	"fmt"
	"testing"
)

func TestConvertMap(t *testing.T) {
	type T struct {
		Name string `sd:"name"`
		Age  int    `sd:"age"`
		sex  int    `sd:"sex"`
	}
	tT := &T{
		Name: "test",
		Age:  18,
		sex:  1,
	}
	res, err := ConvertMap(tT)
	if err != nil {
		t.Fatal("convertMap failed", err)
		return
	}
	fmt.Println(res["$name"])
	fmt.Println(res["$age"])
	fmt.Println(res["$sex"])
}
