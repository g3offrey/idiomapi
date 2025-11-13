package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		sep      string
		expected string
	}{
		{
			name:     "empty slice",
			strs:     []string{},
			sep:      ", ",
			expected: "",
		},
		{
			name:     "single element",
			strs:     []string{"foo"},
			sep:      ", ",
			expected: "foo",
		},
		{
			name:     "multiple elements with comma",
			strs:     []string{"foo", "bar", "baz"},
			sep:      ", ",
			expected: "foo, bar, baz",
		},
		{
			name:     "multiple elements with space",
			strs:     []string{"a", "b", "c"},
			sep:      " ",
			expected: "a b c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinStrings(tt.strs, tt.sep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrNotFound(t *testing.T) {
	assert.NotNil(t, ErrNotFound)
	assert.Equal(t, "todo not found", ErrNotFound.Error())
}
