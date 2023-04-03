package main

import "fmt"

type demo struct {
	Url string `json:"url" validate:"required,ipv4"`
}

func main() {
	data := make(map[string]string)
	val := data["SSSSg"]
	fmt.Println("value:=", val)
	// s := string(`{"operation": "get", "key": "example"}`)
	// err := validator.New().Struct(&demo{
	// 	Url: "192.168.1.11",
	// })
	// if err != nil {
	// 	println(err.Error())
	// }
}
