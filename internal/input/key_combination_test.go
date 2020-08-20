package input

import (
	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatKeyCombinationStringFunctionWorksCorrectlyForSingleKeys(t *testing.T) {
	keyCombination := SingleKeyCombination(fyne.KeyN)
	assert.Equal(t, "N", keyCombination.String())
}

func TestThatKeyCombinationStringFunctionWorksCorrectlyForDoubleKeys(t *testing.T) {
	keyCombination := TwoKeyCombination(fyne.KeyR, fyne.KeyL)
	assert.Equal(t, "R,L", keyCombination.String())
}
