package main

import (
	"encoding/json"
	"fmt"
)

//把请求包定义成一个结构体
type Requestbody struct {
	req string
}

//以指针的方式传入，但在使用时却可以不用关心
// result 是函数内的临时变量，作为返回值可以直接返回调用层
func (r *Requestbody) Json2map() (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
		return nil, err
	}
	return result, nil
}
func main() {
	//json转map
	var r Requestbody
	r.req = `{"name": "xym","sex": "male"}`
	if req2map, err := r.Json2map(); err == nil {
		fmt.Println(req2map["name"])
		fmt.Println(req2map)
	} else {
		fmt.Println(err)
	}
}
