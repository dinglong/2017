package main

import (
	"bytes"
	"encoding/json"
	"log"
    "fmt"
)

// ---------------------------------------------------------------------------------------------------
// json的 struct field's tag 规范定义如下：
// go语言里，StructTag是一个标记字符串，此字符串可跟随在Struct中字段定义的后面。
// StructTag就是一系列的 key:”value” 形式的组合，其中key是一个不可为空的字符串，key-value组合可以有多个，空格分隔。
// 在StructTag中加入”omitempty”, 标识该字段的数据可忽略。
// ---------------------------------------------------------------------------------------------------

type X struct {
	A string `json:"a,omitempty"`
	B string `json:",omitempty"`
	C string `json:"c"`
	D string
}

func main() {
	x := X{
		A: "aaa",
		B: "bbb",
		C: "ccc",
		D: "ddd",
	}

	b, err := json.Marshal(x)
	if err != nil {
		log.Fatalf("Marshal struct X error: %v\n", err)
	}

	var buff bytes.Buffer
	if err = json.Indent(&buff, b, "", "  "); err != nil {
		log.Fatalf("Marshal struct X error: %v\n", err)
	}

	fmt.Printf("--- json --- \n%s\n-----------\n", buff.String())
}
