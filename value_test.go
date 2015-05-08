package config

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestResolveFlag(t *testing.T) {
	fv := 2
	v := value{Flag: &fv}
	assert.Equal(t, fv, v.resolve())
}

func TestResolveEnv(t *testing.T) {
	ev := 2
	v := value{Env: ev, UseEnv: true}
	assert.Equal(t, ev, v.resolve())
}
