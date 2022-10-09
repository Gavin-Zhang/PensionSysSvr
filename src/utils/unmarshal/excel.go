package unmarshal

import (
	"errors"
	"fmt"
	"reflect"
)

// ExcelUnmarshal 按excel表字段解码成对象
func ExcelUnmarshal(data []string, keys []string, ptr interface{}, classptr interface{}) error {
	val := reflect.ValueOf(ptr)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return errors.New("Is not a Struct!")
	}

	findIndex := func(keyName string) (int, error) {
		for i := 0; i < len(keys); i++ {
			if keys[i] == keyName {
				return i, nil
			}
		}
		return -1, fmt.Errorf("keyName[%s] not exit", keyName)
	}

	other := func(classptr interface{}, str string, name string) error {
		refVal := reflect.ValueOf(classptr)
		refTyp := refVal.Type()
		if refTyp.Kind() == reflect.Ptr {
			refTyp = refTyp.Elem()
		}

		className := refTyp.Name()
		funcName := fmt.Sprintf("%s%s", "AnalyzeStruct", name)

		mothodVal := refVal.MethodByName(funcName)
		if !mothodVal.IsValid() {
			return fmt.Errorf("%s don't have function %s()", className, funcName)
		}
		retVal := mothodVal.Call([]reflect.Value{reflect.ValueOf(ptr), reflect.ValueOf(str)})[0]
		if retVal.Interface() != nil {
			err := retVal.Interface().(error)
			return err
		}
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {
		kind := val.Field(i).Kind()
		name := typ.Field(i).Name
		val := val.FieldByName(name)

		index, err := findIndex(name)
		if err != nil {
			return err
		}
		str := data[index]

		if kind == reflect.String {
			val.SetString(str)
		} else if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
			s, _ := StringToInt(str)
			val.SetInt(s)
		} else if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
			s, _ := StringToUint(str)
			val.SetUint(s)
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			s, _ := StringToFloat(str)
			val.SetFloat(s)
		} else if kind == reflect.Bool {
			s, _ := StringToBool(str)
			val.SetBool(s)
		} else if kind == reflect.Array || kind == reflect.Map || kind == reflect.Slice || kind == reflect.Struct || kind == reflect.Ptr {
			if classptr == nil {
				return fmt.Errorf("field [%s] not set!%v %v \n", name, kind, typ.Field(i).Type)
			}
			err := other(classptr, str, name)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("field [%s] not set!%v %v \n", name, kind, typ.Field(i).Type)
		}

	}

	return nil
}

// ExcelUnmarshalKeyValue 按excel表首字段的值解码成对象
func ExcelUnmarshalKeyValue(data [][]string, ptr interface{}, classptr interface{}) error {
	val := reflect.ValueOf(ptr)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return errors.New("Is not a Struct!")
	}

	finder := func(keyname string) string {
		for i := 0; i < len(data); i++ {
			if data[i][0] == keyname {
				return data[i][1]
			}
		}
		return ""
	}

	other := func(classptr interface{}, str string, name string) error {
		refVal := reflect.ValueOf(classptr)
		refTyp := refVal.Type()
		if refTyp.Kind() == reflect.Ptr {
			refTyp = refTyp.Elem()
		}

		className := refTyp.Name()
		funcName := fmt.Sprintf("AnalyzeStruct%s", name)

		mothodVal := refVal.MethodByName(funcName)
		if !mothodVal.IsValid() {
			return fmt.Errorf("%s don't have function %s()", className, funcName)
		}
		retVal := mothodVal.Call([]reflect.Value{reflect.ValueOf(ptr), reflect.ValueOf(str)})[0]
		if retVal.Interface() != nil {
			err := retVal.Interface().(error)
			return err
		}
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {
		kind := val.Field(i).Kind()
		keyName := typ.Field(i).Name
		val := val.FieldByName(keyName)

		value := finder(keyName)
		if value == "" {
			//return fmt.Errorf("Table don't have key[%s]", keyName)
			fmt.Println("Table don't have key[%s]", keyName)
		}
		if kind == reflect.String {
			val.SetString(value)
		} else if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
			s, _ := StringToInt(value)
			val.SetInt(s)
		} else if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
			s, _ := StringToUint(value)
			val.SetUint(s)
		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			s, _ := StringToFloat(value)
			val.SetFloat(s)
		} else if kind == reflect.Bool {
			s, _ := StringToBool(value)
			val.SetBool(s)
		} else if kind == reflect.Array ||
			kind == reflect.Map ||
			kind == reflect.Slice ||
			kind == reflect.Struct ||
			kind == reflect.Ptr {
			if classptr == nil {
				return fmt.Errorf("field [%s] not set!%v %v \n", keyName, kind, typ.Field(i).Type)
			}
			err := other(classptr, value, keyName)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("field [%s] not set!%v %v \n", keyName, kind, typ.Field(i).Type)
		}
	}
	return nil
}
