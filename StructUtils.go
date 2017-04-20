package kezban

import (
	"reflect"
	"fmt"
	"errors"
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

func FillStruct(s interface{}, m map[string]interface{}) error {
	for k,v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(s interface{}, name string, value interface{}) error {
	s_value := reflect.ValueOf(s)
	if s_value.Kind() == reflect.Ptr {
		s_value = s_value.Elem()
	}

	sFieldVal := s_value.FieldByName(name)

	if !sFieldVal.IsValid() {
		return errors.New("No such field named: " + name)
	}
	if !sFieldVal.CanSet() {
		return errors.New("Field("+name+") cannot set.")
	}
	sFieldType := sFieldVal.Type()
	val := reflect.ValueOf(value)
	if sFieldType != val.Type() {
		return errors.New("Field and value types didnt match.")
	}
	sFieldVal.Set(val)
	return nil

}

func createEmptyStruct(s interface{}) interface{} {
	svalue := reflect.ValueOf(s)
	if svalue.Kind() == reflect.Ptr {
		svalue = svalue.Elem()
	}
	stype := svalue.Type()
	return reflect.New(stype).Interface()
}