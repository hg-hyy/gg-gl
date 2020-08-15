package handler

import (
	"fmt"
	"reflect"
)

// Student ...
type Student struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Content string
}

func (s Student) hellogo() {
	fmt.Println(s.Name)
}

func get(i interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	num := v.NumField()
	fmt.Println(num)
	for i := 0; i < num; i++ {
		tagval := t.Field(i)
		if tagval.Tag.Get("json") != "" {
			fmt.Println(tagval.Tag.Get("json"))
		}
	}
	numofmethod := v.NumMethod()
	fmt.Println(numofmethod)
}

// MyStruct ...
type MyStruct struct {
	N int
}

// Printreflect ...
func Printreflect() {
	var Students Student
	get(Students)

	n := MyStruct{1}

	// get
	immutable := reflect.ValueOf(n)
	val := immutable.FieldByName("N").Int()
	fmt.Printf("N=%d\n", val) // prints 1

	// set
	mutable := reflect.ValueOf(&n).Elem()
	mutable.FieldByName("N").SetInt(7)
	fmt.Printf("N=%d\n", n.N) // prints 7
}
