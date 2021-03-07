package input

//Represents an action that should be executed when certain keys are pressed
type Action string

const (
	SelectNextTabAction        Action = "SELECT_NEXT_TAB"
	SelectPreviousTabAction    Action = "SELECT_PREVIOUS_TAB"
	SaveChangesAction          Action = "SAVE_CHANGES"
	AddEntryTypeAction         Action = "ADD_ENTRY_TYPE"
	RemoveEntryTypeAction      Action = "REMOVE_ENTRY_TYPE"
	EditCurrentEntryTypeAction Action = "EDIT_CURRENT_ENTRY_TYPE"
	MoveDownAction             Action = "MOVE_DOWN"
	MoveUpAction               Action = "MOVE_UP"
	EnterInputModeAction       Action = "ENTER_INPUT_MODE"
	ExitInputModeAction        Action = "EXIT_INPUT_MODE"
	ExitTableAction            Action = "EXIT_TABLE"
	ConfirmAction              Action = "CONFIRM"
	CancelAction               Action = "CANCEL"
)
