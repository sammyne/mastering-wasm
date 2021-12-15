package tools

import (
	"github.com/sammyne/mastering-wasm/wavm"
	"github.com/sammyne/mastering-wasm/wavm/validator"
)

func Check(m *wavm.Module) error {
	return validator.Validate(*m)
}
