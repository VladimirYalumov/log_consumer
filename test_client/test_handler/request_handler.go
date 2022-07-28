package test_handler

import (
	"log_consumer/logger"
	"fmt"
)

func AddNumbers(a int, b int) int {
	result := a + b
	datas := make(map[string]interface{})
	datas["msg"] = fmt.Sprintf("result sum of %d and %d is %d", a, b, result)
	_ = logger.GetInstance().InfoPush(datas)
	return result
}

func SubtractNumbers(a int, b int) int {
	result := a - b
	datas := make(map[string]interface{})
	datas["msg"] = fmt.Sprintf("result sum of %d and %d is %d", a, b, result)
	_ = logger.GetInstance().InfoPush(datas)
	return result
}
