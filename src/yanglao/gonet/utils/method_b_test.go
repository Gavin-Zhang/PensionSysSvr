package utils

import (
	"testing"
)

type testBenchmarkMethodClass struct {
}

func (bmc *testBenchmarkMethodClass) Test(a1 int, a2 int, a3 int) int {
	return 100
}

func BenchmarkCallMethod(b *testing.B) {
	obj := &testBenchmarkMethodClass{}
	args := []interface{}{100, 200, 300}

	for i := 0; i < b.N; i++ {
		result, err := CallMethod(obj, "Test", args)
		if (err != nil) || (result[0] != 100) {
			b.Error("CallMethod压力测试失败")
			break
		}
	}
}

func BenchmarkExpandResult(b *testing.B) {
	result := []interface{}{2000, 3000, 4000}
	var arg1, arg2, arg3 int

	for i := 0; i < b.N; i++ {
		err := ExpandResult(result, &arg1, &arg2, &arg3)
		if (err != nil) || (arg1 != 2000) || (arg2 != 3000) || (arg3 != 4000) {
			b.Error("ExpandResult压力测试失败")
			break
		}
	}
}
