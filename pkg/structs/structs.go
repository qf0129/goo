package structs

import (
	"fmt"
	"reflect"
)

func GetFields(structObj any) []string {
	t := reflect.TypeOf(structObj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		fmt.Println("Check type error not Struct")
		return nil
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}

	return result
}
