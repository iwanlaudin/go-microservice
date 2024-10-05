package helpers

import "reflect"

func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(input)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i).Interface()
		result[field.Name] = fieldValue
	}

	return result
}
