// testValidator project testValidator.go
package main

import (
	"log"
	. "testValidator/validator"
)

type testValidator struct {
	ItemNum int `valid:"PosNO" name:"items" ` // 验证是否是正数
	IsCount int `valid:"PosNO"`
}

func main() {
	// validator := NewValidator()
	testObj := &testValidator{
		ItemNum: -1,
		IsCount: 1,
	}

	err := Validate.ValidateData(testObj)

	log.Println(err, 44)
}
