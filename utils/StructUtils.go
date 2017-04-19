package utils

import (
	"reflect"
	"fmt"
)

func ReadTagAndValue(val interface{}) {
	value := reflect.ValueOf(val)
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag
		field := value.Field(i)
		fmt.Printf("%v\t%v\t\n", tag.Get("kezban"), field.String())
	}
}