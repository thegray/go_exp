package structs

import "github.com/go-playground/validator"

type AppConfig struct {
	Port       string `json:"port"`
	Secret     string `json:"secret"`
	ClientPort string `json:"clport"`
}

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
