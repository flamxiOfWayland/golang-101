package customerrors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRich(t *testing.T) {
	result := AreWeRich()
	if result != nil {
		fmt.Println("we are not rich")
		fmt.Println(result.Error())
	}
	assert.False(t, true)
}
