// Package validator provides a centralized way to register all custom validators.
package validators

import (
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidatorRegisterer interface {
	Register(validate *validator.Validate) error
}

var (
	registerers  []ValidatorRegisterer
	registerOnce sync.Once
)

func Register(vr ValidatorRegisterer) {
	registerers = append(registerers, vr)
}

func RegisterAll() error {
	var err error
	registerOnce.Do(func() {
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
