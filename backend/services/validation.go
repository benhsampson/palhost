package services

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type StructValidation[S any] struct {
	callback func(*S) error
	t        *testing.T
}

type ValidationMap map[string][]string

func (a *StructValidation[S]) Test(s *S, rules ValidationMap) {
	err := a.callback(s)
	assert.Error(a.t, err)
	valErrs := err.(validator.ValidationErrors)
	assert.Equal(a.t, len(valErrs), len(rules), fmt.Sprintf("got %d errors, expected %d", len(valErrs), len(rules)))
	for _, valErr := range valErrs {
		field := valErr.Field()
		tag := valErr.Tag()
		assert.Contains(a.t, rules[field], tag, fmt.Sprintf("got unexpected error %s for field %s", tag, field))
	}
}
