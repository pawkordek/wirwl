package input

import (
	"fyne.io/fyne"
	"strings"
	"time"
)

//Represents an action that should be executed when certain keys are pressed
type Action string

type keyCombination struct {
	firstKey  fyne.KeyName
	secondKey fyne.KeyName
}

func (keyCombination *keyCombination) press(key fyne.KeyName) {
	if keyCombination.firstKey == "" {
		keyCombination.firstKey = key
	} else if keyCombination.secondKey == "" {
		keyCombination.secondKey = key
	} else {
		keyCombination.firstKey = key
		keyCombination.secondKey = ""
	}
}

func (keyCombination *keyCombination) oneKeyPressed() bool {
	return keyCombination.firstKey != "" && keyCombination.secondKey == ""
}

func (keyCombination *keyCombination) bothKeysPressed() bool {
	return keyCombination.firstKey != "" && keyCombination.secondKey != ""
}

func (keyCombination *keyCombination) releaseKeys() {
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
type InputHandler struct {
	keymap                map[keyCombination]Action
	actions               map[CallerActionPair]func()
	currentKeyCombination keyCombination
	lastKeyPressTime      time.Time
}

func NewInputHandler(keymap map[string]Action) InputHandler {
	handler := InputHandler{keymap: map[keyCombination]Action{}, actions: map[CallerActionPair]func(){}, lastKeyPressTime: time.Now()}
	handler.createActualKeymap(keymap)
	return handler
}

func (handler *InputHandler) createActualKeymap(keymap map[string]Action) {
	for key, action := range keymap {
		keyCombination := getKeyCombinationFromStringKey(key)
		handler.keymap[keyCombination] = action
	}
}

func getKeyCombinationFromStringKey(key string) keyCombination {
	if strings.Contains(key, ",") {
		keys := strings.Split(key, ",")
		return keyCombination{
			firstKey:  fyne.KeyName(keys[0]),
			secondKey: fyne.KeyName(keys[1]),
		}
	}
	return keyCombination{
		firstKey:  fyne.KeyName(key),
		secondKey: "",
	}
}

func (handler *InputHandler) BindFunctionToAction(caller interface{}, action Action, function func()) {
	callerActionPair := CallerActionPair{
		caller: caller,
		action: action,
	}
	handler.actions[callerActionPair] = function
}

func (handler *InputHandler) Handle(caller interface{}, keyName fyne.KeyName) {
	handler.currentKeyCombination.press(keyName)
	actionToExecute := handler.keymap[handler.currentKeyCombination]
	timeNow := time.Now()
	timeSinceLastKeyPress := timeNow.Sub(handler.lastKeyPressTime)
	if handler.currentKeyCombination.oneKeyPressed() ||
		(handler.currentKeyCombination.bothKeysPressed() && timeSinceLastKeyPress < time.Second) {
		callerActionPair := CallerActionPair{
			caller: caller,
			action: actionToExecute,
		}
		function := handler.actions[callerActionPair]
		if function != nil {
			function()
			handler.currentKeyCombination.releaseKeys()
			return
		}
	}
	handler.lastKeyPressTime = timeNow
}
