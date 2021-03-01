package ptrutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	i := 1
	iPtr := Int(i)
	assert.Equal(t, i, *iPtr)
}

func TestString(t *testing.T) {
	i := "a"
	iPtr := String(i)
	assert.Equal(t, i, *iPtr)
}

func TestTime(t *testing.T) {
	i := time.Now()
	iPtr := Time(i)
	assert.Equal(t, i, *iPtr)
}

func TestBool(t *testing.T) {
	i := true
	iPtr := Bool(i)
	assert.Equal(t, i, *iPtr)
}
