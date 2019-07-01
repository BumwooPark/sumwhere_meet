package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	hashed := HashAndSalt([]byte("1q2w3e4r"))
	assert.True(t, ComparePasswords(hashed, []byte("1q2w3e4r")))
}
