package wirwl

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

//Stores key combinations mapped to actions. Every action that should be handled, should have a function bound to it
//which will get executed when key combination for that action gets pressed
type InputHandler struct {
	keymap  map[KeyCombination]Action
	actions map[Action]func()
	lastKey fyne.KeyName
}

func NewInputHandler(keymap map[string]Action) InputHandler {
	handler := InputHandler{keymap: map[KeyCombination]Action{}, actions: map[Action]func(){}}
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

func (handler *InputHandler) bindFunctionToAction(action Action, function func()) {
	handler.actions[action] = function
}

func (handler *InputHandler) handle(keyName fyne.KeyName) {
	for keyCombination, action := range handler.keymap {
		if (keyCombination.secondKey == keyName && keyCombination.firstKey == handler.lastKey) ||
			(keyCombination.firstKey == keyName && keyCombination.secondKey == "") {
			handler.actions[action]()
			break
		}
	}
	handler.lastKey = keyName
}
