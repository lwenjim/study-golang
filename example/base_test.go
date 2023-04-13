package example

import (
	"fmt"
	"strconv"
	"testing"
)

func qTestInteger(t *testing.T) {
	a := struct {
		Name string
		Age  int
		Desc string
	}{
		Name: "lwenjim",
		Age:  32,
		Desc: "abc",
	}
	//%T flag 表示打印go语言结构
	fmt.Printf("%T\n", a)
	b := []string{"abc"}
	fmt.Printf("%T\n", b)

	//%% flag 打印百分号
	fmt.Printf("%%\n")

	//%t 打印 true or false
	c := true
	fmt.Printf("%t\n", c)

	//%b 打印二进制码
	e := 2
	fmt.Printf("%b\n", e)

	//打印Unicode 码点
	f := 65
	fmt.Printf("%c\n", f)

	//打印十进制
	fmt.Printf("%d\n", 0o343)

	//%o 打印八进制
	fmt.Printf("%o\n", 0x9)
	fmt.Printf("%O\n", 0x9)

	//%q 码点转字符
	fmt.Printf("%q\n", 65)

	//%x 打印十六进制
	fmt.Printf("%x\n", 651)
	fmt.Printf("%X\n", 651)
	fmt.Printf("%U\n", 651)
	fmt.Printf("U+%04X\n", 651)
}
func TestFloat(t *testing.T) {
	fmt.Printf("%b\n", 1230.98)
	fmt.Printf("%s\n", strconv.FormatFloat(0.112, 'b', 1, 32))
	fmt.Printf("%e\n", 123.12)
	fmt.Printf("%E\n", 123.12)
	fmt.Printf("%f\n", 123.12)
	fmt.Printf("%g\n", 123.12)
	fmt.Printf("%G\n", 123.12)
	fmt.Printf("%x\n", 123.12)
	fmt.Printf("%X\n", 123.12)

	fmt.Printf("%012.3f\n", 123.128098)

	fmt.Printf("%010.2e\n", 123.128098)
	fmt.Printf("%09.2g\n", 123.128098)
	fmt.Printf("%#g\n", 123.128098)

	fmt.Printf("%+022f\n", -123.128098)
	fmt.Printf("%#b\n", 0b11111)
	fmt.Printf("%#o\n", 0b11111)
	fmt.Printf("%#x\n", 0b11111)
	fmt.Printf("%#X\n", 0b11111)
	p := []byte("abc")
	c := &p
	fmt.Printf("%#p\n", &c)
	fmt.Printf("%p\n", c)

	fmt.Printf("%e\n", 244340.424234)
	fmt.Printf("%#e\n", 244340.424234)

	fmt.Printf("%0f\n", 40.9)

	// fmt.Printf("% x\n", "abc")
	// fmt.Printf("% x\n", []byte("abc"))
	// fmt.Printf("% d\n", []byte("abc"))

}
func aTestString(t *testing.T) {
	fmt.Printf("%s\n", "abc")
	fmt.Printf("%q\n", "abc\nab\tafsf\"")
	fmt.Printf("%x\n", "ab")
	fmt.Printf("%X\n", "ab")
	fmt.Printf("%d\n", int('a'))
	fmt.Printf("%x\n", 97)

	fmt.Printf("%q\n", `\a`)
	fmt.Printf("%q\n", "\ta")
	fmt.Printf("%q\n", 'a')
}

func aTestSlice(t *testing.T) {
	p := []string{"a"}
	fmt.Printf("%p\n", p)
	fmt.Printf("%T\n", &p)
}

func aTestPointer(t *testing.T) {
	p := &[]string{"a"}
	fmt.Printf("%T\n", p)
	fmt.Printf("%b\n", &p)
	fmt.Printf("%d\n", &p)
	fmt.Printf("%x\n", &p)
	fmt.Printf("%o\n", &p)
}
