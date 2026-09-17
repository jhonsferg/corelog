package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	instance *validator.Validate
)

func Get() *validator.Validate {
	once.Do(func() {
		instance = validator.New(validator.WithRequiredStructEnabled())
	})
	return instance
}

func Struct(s any) error {
	return Get().Struct(s)
}
