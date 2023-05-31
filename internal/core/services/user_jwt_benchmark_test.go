package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkUserJwtToken(b *testing.B) {
	id := uuid.New()
	name := "John Doe"

	for n := 0; n < b.N; n++ {
		_, err := UserJwtToken(id, name)
		require.NoError(b, err)
	}
}
