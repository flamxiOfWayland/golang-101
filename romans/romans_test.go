package romans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var input map[string]int = map[string]int{"I": 1, "II": 2, "V": 5}

func TestRomansToIntegers(t *testing.T) {
	for roman, integer := range input {
		assert.Equal(t, integer, RomanToInteger(roman))
	}
}
