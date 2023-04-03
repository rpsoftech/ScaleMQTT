package main

import "github.com/go-playground/validator/v10"

type demo struct {
	Url string `json:"url" validate:"required,ipv4"`
}

func main() {
	// s := string(`{"operation": "get", "key": "example"}`)
	err := validator.New().Struct(&demo{
		Url: "192.168.1.11",
	})
	if err != nil {
		println(err.Error())
	}
}
