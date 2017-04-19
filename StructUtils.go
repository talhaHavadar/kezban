package kezban

import (
	"reflect"
	"fmt"
)

func ReadTagAndValue(val interface{}) {
	value := reflect.ValueOf(val).Elem()
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag
		field := value.Field(i)
		fmt.Printf("%v\t%v\t\n", tag.Get("kezban"), field.String())
	}
}

func GetFields(model interface{}, search string) map[string]interface{} {
	res := make(map[string]interface{})
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag
		field := value.Field(i)
		if tagStr := tag.Get("kezban"); tagStr != "" {
			if tagStr == search {
				res[value.Type().Field(i).Name] = field.Interface()
			}
		}
	}
	return res
}