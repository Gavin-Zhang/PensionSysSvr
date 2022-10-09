package utils

import (
	"testing"
)

type testMethodClass struct {
	callMethodA  bool
	callMethodA1 bool
	callMethodA2 bool
	callMethodB  bool
	methodBArg1  int
	callMethodB1 bool
	methodB1Arg1 int
	methodB1Arg2 int
}

func (mc *testMethodClass) MethodA() {
	mc.callMethodA = true
}

func (mc *testMethodClass) MethodA1() int {
	mc.callMethodA1 = true
	return 100
}

func (mc *testMethodClass) MethodA2() (r1 int, r2 int) {
	mc.callMethodA2 = true
	r1 = 101
	r2 = 102
	return
}

func (mc *testMethodClass) MethodB(arg1 int) {
	mc.callMethodB = true
	mc.methodBArg1 = arg1
}

func (mc *testMethodClass) MethodB1(arg1 int, arg2 int) int {
	mc.callMethodB1 = true
	mc.methodB1Arg1 = arg1
	mc.methodB1Arg2 = arg2
	return 300
}

func TestCallMethod(t *testing.T) {
	obj := &testMethodClass{}
	var emptyArgs []interface{}
	var result []interface{}
	var err error

	// 测试obj参数为空
	result, err = CallMethod(nil, "MethodA", emptyArgs)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试obj为空 失败.")
	}

	// 测试fname参数为空
	result, err = CallMethod(obj, "", emptyArgs)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试fname为空 失败.")
	}

	// 测试不存在的函数调用
	result, err = CallMethod(obj, "MethodNil", emptyArgs)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试不存在的函数 失败.")
	}

	// 测试空参数
	result, err = CallMethod(obj, "MethodB", emptyArgs)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试空参数 失败.")
	}

	// 测试参数数量不匹配
	args1 := []interface{}{uint32(0)}
	result, err = CallMethod(obj, "MethodB", args1)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试参数数量不匹配 失败.")
	}

	// 测试参数类型不匹配
	args2 := []interface{}{float32(1.0), int(0)}
	result, err = CallMethod(obj, "MethodB", args2)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试参数类型不匹配 失败.")
	}

	// 测试给无参函数传入参数
	args3 := []interface{}{float32(1.0), int(0)}
	result, err = CallMethod(obj, "MethodA", args3)
	if (len(result) != 0) || (err == nil) {
		t.Error("测试给无参函数传入参数 失败.")
	}

	// 测试调用无参无返回值函数
	result, err = CallMethod(obj, "MethodA", emptyArgs)
	if (len(result) != 0) || (err != nil) {
		t.Error("测试调用无参无返回值函数 失败.")
	}

	if obj.callMethodA != true {
		t.Error("测试调用无参无返回值函数 失败.")
	}

	// 测试调用无参有1返回值函数
	result, err = CallMethod(obj, "MethodA1", emptyArgs)
	if (len(result) != 1) || (err != nil) {
		t.Error("测试调用无参有1返回值函数 失败.")
	}

	if (obj.callMethodA1 != true) || (result[0] != 100) {
		t.Error("测试调用无参有1返回值函数 失败.")
	}

	// 测试调用无参有2返回值函数
	result, err = CallMethod(obj, "MethodA2", emptyArgs)
	if (len(result) != 2) || (err != nil) {
		t.Error("测试调用无参有2返回值函数 失败.")
	}

	if (obj.callMethodA2 != true) || (result[0] != 101) || (result[1] != 102) {
		t.Error("测试调用无参有2返回值函数 失败.")
	}

	// 测试调用1参无返回值函数
	args4 := []interface{}{int(200)}
	result, err = CallMethod(obj, "MethodB", args4)
	if (len(result) != 0) || (err != nil) {
		t.Error("测试调用无参有1返回值函数 失败.")
	}

	if (obj.callMethodB != true) || (obj.methodBArg1 != 200) {
		t.Error("测试调用无参有1返回值函数 失败.")
	}

	// 测试调用2参1返回值函数
	args5 := []interface{}{int(210), int(220)}
	result, err = CallMethod(obj, "MethodB1", args5)
	if (len(result) != 1) || (err != nil) {
		t.Error("测试调用2参1返回值函数 失败.")
	}

	if (obj.callMethodB1 != true) ||
		(obj.methodB1Arg1 != 210) ||
		(obj.methodB1Arg2 != 220) ||
		(result[0] != 300) {
		t.Error("测试调用2参1返回值函数 失败.")
	}
}

func TestExpandResult(t *testing.T) {
	result := []interface{}{1000, uint32(204), "Hello"}
	var arg1 int
	var arg2 uint32
	var arg3 string

	var err error

	// 测试长度不匹配
	err = ExpandResult(result, &arg1, &arg2)
	if err == nil {
		t.Error("测试长度不匹配 失败")
	}

	// 测试参数未加&
	err = ExpandResult(result, arg1, &arg2, &arg3)
	if err == nil {
		t.Error("测试参数未加& 失败")
	}

	// 测试参数为空指针
	resultElem := 2000
	result1 := []interface{}{&resultElem, uint32(204)}
	var arg4 *int
	err = ExpandResult(result1, arg4, &arg2)
	if err == nil {
		t.Error("测试参数为空指针 失败")
	}

	// 测试类型不匹配
	var arg5 uint32
	err = ExpandResult(result, &arg5, &arg2, &arg3)
	if err == nil {
		t.Error("测试类型不匹配 失败")
	}

	// 测试参数赋值
	err = ExpandResult(result, &arg1, &arg2, &arg3)
	if err != nil {
		t.Error("测试参数赋值 失败")
	}

	if (arg1 != result[0].(int)) || (arg2 != result[1].(uint32)) || (arg3 != result[2].(string)) {
		t.Error("测试参数赋值 失败")
	}
}
