package input

import (
	"fyne.io/fyne"
	"time"
)

//Caller should be anything that allows to unambiguously find the correct action for that caller
//That is, if there are various objects that use the same struct (e.g. many copies of certain widget), it's best
//to pass the object itself as this guarantees unambiguity
//If struct only ever has one instance, it might be fine to pass something else, e.g. a string (like in App's case)
type callerActionPair struct {
	caller interface{}
	action Action
}

//Stores key combinations mapped to actions. Every action that should be handled, should have a function bound to it
//which will get executed when key combination for that action gets pressed.
//There are two modes in which handler operates:
//First, normal mode where user has to either type a key or a combination of two keys for the action to be executed.
//	Keys are stored as they are input e.g. pressing 'K' then 'G' will either execute action for 'K' if it matches first or action for 'KG' if 'K' doesn't match
//	Pressing another key resets the combination e.g. Pressing 'L' after above combination will make the current combination 'L'
//Second, input mode where on every key press, second key of current combination replaces it's first key and the input key becomes the second
//	In practice this looks like this e.g. press 'Q', current combination is 'Q', then press 'H', combination is 'QH'
//	Press another key, 'U'. Combination becomes 'HU'.
//	Also, an action can be executed based on either key of current combination (single key actions) or the combination itself.
//	That is, in 'HU' both actions for 'H', 'U' or 'HU' can execute but order of precedence is 'HU', 'H', 'U'.
type Handler struct {
	keymap                map[KeyCombination][]Action
	actions               map[callerActionPair]func()
	currentKeyCombination KeyCombination
	lastKeyPressTime      time.Time
	onKeyPressedCallback  func(KeyCombination)
}

func NewHandler(actionKeyMap map[Action]KeyCombination) Handler {
	keyActionMap := convertActionKeyKeymapToKeyCombinationActionKeymap(actionKeyMap)
	handler := Handler{
		keymap:                keyActionMap,
		actions:               map[callerActionPair]func(){},
		currentKeyCombination: KeyCombination{},
		lastKeyPressTime:      time.Now(),
		onKeyPressedCallback: func(combination KeyCombination) {
		},
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
	callerActionPair := callerActionPair{
		caller: caller,
		action: action,
	}
	handler.actions[callerActionPair] = function
}

func (handler *Handler) HandleInNormalMode(caller interface{}, keyName fyne.KeyName) {
	handler.currentKeyCombination.press(keyName)
	handler.onKeyPressedCallback(handler.currentKeyCombination)
	handler.tryExecutingFunctionForCallerAndKeyCombination(caller, handler.currentKeyCombination)
}

func (handler *Handler) HandleInInputMode(caller interface{}, keyName fyne.KeyName) bool {
	if handler.currentKeyCombination.bothKeysPressed() {
		handler.currentKeyCombination.press(handler.currentKeyCombination.secondKey)
	}
	handler.currentKeyCombination.press(keyName)
	handler.onKeyPressedCallback(handler.currentKeyCombination)
	functionExecuted := handler.tryExecutingFunctionForCallerAndKeyCombination(caller, handler.currentKeyCombination)
	if functionExecuted {
		return true
	} else {
		return handler.tryExecutingFunctionForCallerAndKeyCombination(caller, SingleKeyCombination(handler.currentKeyCombination.secondKey))
	}
}

func (handler *Handler) tryExecutingFunctionForCallerAndKeyCombination(caller interface{}, keyCombination KeyCombination) bool {
	timeNow := time.Now()
	defer func() { handler.lastKeyPressTime = timeNow }()
	timeSinceLastKeyPress := timeNow.Sub(handler.lastKeyPressTime)
	for _, action := range handler.keymap[keyCombination] {
		if handler.currentKeyCombination.oneKeyPressed() ||
			(handler.currentKeyCombination.bothKeysPressed() && timeSinceLastKeyPress < time.Second) {
			callerActionPair := callerActionPair{
				caller: caller,
				action: action,
			}
			function := handler.actions[callerActionPair]
			if function != nil {
				function()
				handler.currentKeyCombination.releaseKeys()
				return true
			}
		}
	}
	return false
}

func (handler *Handler) SetOnKeyPressedCallbackFunction(function func(KeyCombination)) {
	handler.onKeyPressedCallback = function
}
