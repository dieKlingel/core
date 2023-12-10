package config_test

import (
	"os"
	"testing"

	"github.com/dieklingel/core/config"
)

func TestNew(t *testing.T) {
	println(os.Getwd())

	env := config.New()

	if env == nil {
		t.Fatalf("expected not nil, got nil")
	}
}
