package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortLink(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want int
	}{
		{
			name: "GenerateShortLink 5",
			arg:  5,
			want: 5,
		},
		{
			name: "GenerateShortLink 0",
			arg:  0,
			want: 0,
		},
		{
			name: "GenerateShortLink -5",
			arg:  -5,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, len(GenerateShortLink(tt.arg)))
		})
	}
}

func TestGenerateShortLink_Contain_Letters(t *testing.T) {
	t.Run("only valid characters", func(t *testing.T) {
		s := GenerateShortLink(1000)
		for _, c := range s {
			assert.Contains(t, letters, string(c),
				"символ %q не из алфавита", string(c))
		}
	})
}

func TestGenerateShortLink_Collisions(t *testing.T) {
	t.Run("are two links distinguish", func(t *testing.T) {
		a := GenerateShortLink(10)
		b := GenerateShortLink(10)

		assert.NotEqual(t, a, b)
	})
}
