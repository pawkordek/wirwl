package input

import (
	"fyne.io/fyne"
	"time"
)

//Represents an action that should be executed when certain keys are pressed
type Action string

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

//Caller should be anything that allows to unambiguously find the correct action for that caller
//That is, if there are various objects that use the same struct (e.g. many copies of certain widget), it's best
//to pass the object itself as this guarantees unambiguity
//If struct only ever has one instance, it might be fine to pass something else, e.g. a string (like in App's case)
type CallerActionPair struct {
	caller interface{}
	action Action
}

//Stores key combinations mapped to actions. Every action that should be handled, should have a function bound to it
//which will get executed when key combination for that action gets pressed
type Handler struct {
	keymap                map[KeyCombination][]Action
	actions               map[CallerActionPair]func()
	currentKeyCombination KeyCombination
	lastKeyPressTime      time.Time
}

func NewInputHandler(actionKeyMap map[Action]KeyCombination) Handler {
	keyActionMap := convertActionKeyKeymapToKeyCombinationActionKeymap(actionKeyMap)
	handler := Handler{
		keymap:                keyActionMap,
		actions:               map[CallerActionPair]func(){},
		currentKeyCombination: KeyCombination{},
		lastKeyPressTime:      time.Now(),
	}
	return handler
}

func convertActionKeyKeymapToKeyCombinationActionKeymap(actionKeyMap map[Action]KeyCombination) map[KeyCombination][]Action {
	keyActionMap := make(map[KeyCombination][]Action)
	for action, keyCombination := range actionKeyMap {
		keyActionMap[keyCombination] = append(keyActionMap[keyCombination], action)
	}
	return keyActionMap
}

func (handler *Handler) BindFunctionToAction(caller interface{}, action Action, function func()) {
	callerActionPair := CallerActionPair{
		caller: caller,
		action: action,
	}
	handler.actions[callerActionPair] = function
}

func (handler *Handler) Handle(caller interface{}, keyName fyne.KeyName) {
	timeNow := time.Now()
	defer func() { handler.lastKeyPressTime = timeNow }()
	timeSinceLastKeyPress := timeNow.Sub(handler.lastKeyPressTime)
	handler.currentKeyCombination.press(keyName)
	for _, action := range handler.keymap[handler.currentKeyCombination] {
		if handler.currentKeyCombination.oneKeyPressed() ||
			(handler.currentKeyCombination.bothKeysPressed() && timeSinceLastKeyPress < time.Second) {
			callerActionPair := CallerActionPair{
				caller: caller,
				action: action,
			}
			function := handler.actions[callerActionPair]
			if function != nil {
				function()
				handler.currentKeyCombination.releaseKeys()
				return
			}
		}
	}
}
