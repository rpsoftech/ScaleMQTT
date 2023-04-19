package validator

import (
	v "github.com/go-playground/validator/v10"
)

var Validator = v.New()

func init() {
	println("Registerd")
	Validator.RegisterValidation("port", ValidatePort)
	// Validator.Re("port", ValidatePort)
}

// func addTranslation(tag string, errMessage string) {
// 	registerFn := func(ut ut.Translator) error {
// 		return ut.Add(tag, errMessage, false)
// 	}

// 	transFn := func(ut ut.Translator, fe validator.FieldError) string {
// 		param := fe.Param()
// 		tag := fe.Tag()

// 		t, err := ut.T(tag, fe.Field(), param)
// 		if err != nil {
// 			return fe.(error).Error()
// 		}
// 		return t
// 	}

// 	_ = Validator.RegisterTranslation(tag, trans, registerFn, transFn)
// }
