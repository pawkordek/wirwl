package input

import (
	"fyne.io/fyne"
	"strings"
)

/*Represents pressed keys. Max two keys can be pressed at once.
This is handled so that after pressing first key, it stays pressed, then after pressing second
key, both are pressed. Pressing a third key releases both keys and third key stays pressed.
Therefore user can press a key or two keys in succession (combination) which is what system can handle.
The third key has to be therefore handled as a new press.
*/
type KeyCombination struct {
	firstKey  fyne.KeyName
	secondKey fyne.KeyName
}

func SingleKeyCombination(key fyne.KeyName) KeyCombination {
	return KeyCombination{
		firstKey:  key,
		secondKey: "",
	}
}

func TwoKeyCombination(firstKey fyne.KeyName, secondKey fyne.KeyName) KeyCombination {
	return KeyCombination{
		firstKey:  firstKey,
		secondKey: secondKey,
	}
}

//Input string format should be either:
//e.g. "J" for singular key
//e.g. "J,Q" for two key combination
func KeyCombinationFromString(keyCombination string) KeyCombination {
	keys := strings.Split(keyCombination, ",")
	if len(keys) > 1 {
		return KeyCombination{
			firstKey:  fyne.KeyName(keys[0]),
			secondKey: fyne.KeyName(keys[1]),
		}
	} else{
		return KeyCombination{
			firstKey:  fyne.KeyName(keys[0]),
			secondKey: "",
		}
	}
}

func (keyCombination *KeyCombination) press(key fyne.KeyName) {
	if keyCombination.firstKey == "" {
		keyCombination.firstKey = key
	} else if keyCombination.secondKey == "" {
		keyCombination.secondKey = key
	} else {
		keyCombination.firstKey = key
		keyCombination.secondKey = ""
	}
}

func (keyCombination *KeyCombination) oneKeyPressed() bool {
	return keyCombination.firstKey != "" && keyCombination.secondKey == ""
}

func (keyCombination *KeyCombination) bothKeysPressed() bool {
	return keyCombination.firstKey != "" && keyCombination.secondKey != ""
}

func (keyCombination *KeyCombination) releaseKeys() {
	keyCombination.firstKey = ""
	keyCombination.secondKey = ""
}

func (keyCombination *KeyCombination) String() string {
	if keyCombination.secondKey != "" {
		return string(keyCombination.firstKey + "," + keyCombination.secondKey)
	} else {
		return string(keyCombination.firstKey)
	}
}
