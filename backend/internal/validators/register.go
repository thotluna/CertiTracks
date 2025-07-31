// Package validators provides a centralized way to register all custom validators.
package validators

import (
	"sync"

	"certitrack/internal/services/auth"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidatorRegisterer interface {
	Register(validate *validator.Validate) error
}

var (
	registerers  = make([]ValidatorRegisterer, 0)
	registerOnce sync.Once
)

func Register(vr ValidatorRegisterer) {
	registerers = append(registerers, vr)
}

func registerAuthValidators() {
	Register(auth.NewAuthValidators())
}

func RegisterAll() error {
	var err error
	registerOnce.Do(func() {
		registerAuthValidators()

		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}

		for _, r := range registerers {
			if e := r.Register(v); e != nil {
				err = e
				return
			}
		}
	})

	return err
}
