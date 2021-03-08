package input

import (
	"fyne.io/fyne/v2"
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

func TestThatKeyCombinationGetsCreatedCorrectlyFromStringWithOneKey(t *testing.T) {
	keyCombination := KeyCombinationFromString("V")
	assert.Equal(t, fyne.KeyV, keyCombination.firstKey)
	assert.Empty(t, keyCombination.secondKey)
}

func TestThatKeyCombinationGetsCreatedCorrectlyFromStringWithTwoKeys(t *testing.T) {
	keyCombination := KeyCombinationFromString("H,O")
	assert.Equal(t, fyne.KeyH, keyCombination.firstKey)
	assert.Equal(t, fyne.KeyO, keyCombination.secondKey)
}

func TestThatKeyCombinationCorrectlyKnowsWhetherItHasOneKeyPressed(t *testing.T) {
	keyCombination := SingleKeyCombination(fyne.KeyR)
	assert.True(t, keyCombination.OneKeyPressed())
}

func TestThatKeyCombinationCorrectlyKnowsWhetherItHasTwoKeysPressed(t *testing.T) {
	keyCombination := TwoKeyCombination(fyne.KeyE, fyne.KeyN)
	assert.True(t, keyCombination.BothKeysPressed())
}
