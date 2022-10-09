package utils

import (
	"errors"
	"fmt"
	"reflect"
	// "strings"
	"runtime/debug"

	"log"
	"os"
	"time"
)

func TypeName(obj interface{}) string {
	str := fmt.Sprintf("%#v", obj)
	beginIndex := 0
	endIndex := 0
	for i, c := range str {
		if c == '.' {
			beginIndex = i + 1
			continue
		}

		if c == '{' {
			endIndex = i
			break
		}
	}

	return str[beginIndex:endIndex]
}

// CallMethod 通过reflect实现调用对象成员方法
func CallMethod(obj interface{}, fname string, args []interface{}) (result []interface{}, e error) {
	defer func() {
		if err := recover(); err != nil {
			errorInfo := fmt.Sprint(err)
			e = errors.New(fmt.Sprintf("CallMethod: %s::%s [%s].", TypeName(obj), fname, errorInfo))
			stack := debug.Stack()

			info := make([]string, 0)
			info = append(info, e.Error())
			info = append(info, string(stack))
			StackLog(info)
		}
	}()

	if obj == nil {
		e = errors.New("CallMethod: obj为空.")
		return
	}

	if len(fname) == 0 {
		e = errors.New(fmt.Sprintf("CallMethod: fname为空."))
		return
	}

	targetObject := reflect.ValueOf(obj)

	argsValue := make([]reflect.Value, len(args))
	for i, v := range args {
		//if v != nil {
		//	argsValue[i] = reflect.ValueOf(v)
		//} else {
		//	var nil_interfack interface{} = nil
		//	argsValue[i] = reflect.ValueOf(nil_interfack)
		//}
		argsValue[i] = reflect.ValueOf(v)
	}

	method := targetObject.MethodByName(fname)

	resultValue := method.Call(argsValue)
	result = make([]interface{}, len(resultValue))
	for i, v := range resultValue {
		result[i] = v.Interface()
	}

	return
}

func ExpandResult(result []interface{}, args ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			errorInfo := fmt.Sprint(err)
			e = errors.New(fmt.Sprintf("ExpandResult: %s.", errorInfo))
			return
		}
	}()

	if len(result) != len(args) {
		return errors.New("ExpandResult: 参数不匹配.")
	}

	var resultValue reflect.Value
	var argValue reflect.Value
	var elemValue reflect.Value

	for i, v := range result {
		argValue = reflect.ValueOf(args[i])
		if argValue.Kind() != reflect.Ptr {
			errorInfo := fmt.Sprintf("ExpandResult: 第%d个参数前必须加&.", i)
			return errors.New(errorInfo)
		}

		elemValue = argValue.Elem()
		if argValue.IsNil() {
			errorInfo := fmt.Sprintf("ExpandResult: 第%d个参数为空指针.", i)
			return errors.New(errorInfo)
		}
		resultValue = reflect.ValueOf(v)

		if elemValue.Type() != resultValue.Type() {
			errorInfo := fmt.Sprintf("ExpandResult: 第%d个参数[%s无法给%s赋值].", i, resultValue.Type(), elemValue.Type())
			return errors.New(errorInfo)
		}

		elemValue.Set(resultValue)
	}

	return nil
}

func StackLog(info []string) {
	_, err := os.Stat("errors")
	if err != nil {
		os.MkdirAll("errors", 0700)
	}

	filepath := fmt.Sprintf("errors/err_stack_%s", time.Now().Format("2006-01-02"))
	logfile, openErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr != nil {
		panic(openErr)
	}
	defer logfile.Close()
	logger := log.New(logfile, "", log.Ldate|log.Ltime)

	logger.Println("===============================================================")
	for _, param := range info {
		logger.Println(param)
		log.Println(param)
	}
}

func StackLogContext(handle uint32, name string, info []string) {
	_, err := os.Stat("errors")
	if err != nil {
		os.MkdirAll("errors", 0700)
	}

	filepath := fmt.Sprintf("errors/err_stack_%s[%s]", time.Now().Format("2006-01-02"), name)
	logfile, openErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr != nil {
		panic(openErr)
	}
	defer logfile.Close()
	logger := log.New(logfile, "", log.Ldate|log.Ltime)

	logger.Println("===============================================================")
	for _, param := range info {
		logger.Println(param)
		log.Println(param)
	}
}

func LogFuncCall(ctx_name string, func_name string, args ...interface{}) {
	_, err := os.Stat("calls")
	if err != nil {
		os.MkdirAll("calls", 0700)
	}

	filepath := fmt.Sprintf("calls/%s-%s", time.Now().Format("2006-01-02"), ctx_name)
	logfile, openErr := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr != nil {
		panic(openErr)
	}
	defer logfile.Close()
	logger := log.New(logfile, "", log.Ldate|log.Ltime)
	logger.Println(fmt.Sprintf("call function:%s", func_name), args)
}
