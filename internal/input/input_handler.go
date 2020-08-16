package input

import (
	"fyne.io/fyne"
	"strings"
)

//Represents an action that should be executed when certain keys are pressed
type Action string

type KeyCombination struct {
	firstKey  fyne.KeyName
	secondKey fyne.KeyName
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
	keymap  map[KeyCombination]Action
	actions map[CallerActionPair]func()
	lastKey fyne.KeyName
}

func NewInputHandler(keymap map[string]Action) InputHandler {
	handler := InputHandler{keymap: map[KeyCombination]Action{}, actions: map[CallerActionPair]func(){}}
	handler.createActualKeymap(keymap)
	return handler
}

func (handler *InputHandler) createActualKeymap(keymap map[string]Action) {
	for key, action := range keymap {
		keyCombination := getKeyCombinationFromStringKey(key)
		handler.keymap[keyCombination] = action
	}
}

func getKeyCombinationFromStringKey(key string) KeyCombination {
	if strings.Contains(key, ",") {
		keys := strings.Split(key, ",")
		return KeyCombination{
			firstKey:  fyne.KeyName(keys[0]),
			secondKey: fyne.KeyName(keys[1]),
		}
	}
	return KeyCombination{
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
	for keyCombination, action := range handler.keymap {
		if (keyCombination.secondKey == keyName && keyCombination.firstKey == handler.lastKey) ||
			(keyCombination.firstKey == keyName && keyCombination.secondKey == "") {
			callerActionPair := CallerActionPair{
				caller: caller,
				action: action,
			}
			function := handler.actions[callerActionPair]
			if function != nil {
				function()
			}
			break
		}
	}
	handler.lastKey = keyName
}
