package wirwl

import (
	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"wirwl/internal/input"
	"wirwl/internal/log"
)

func TestThatErrorGetsReturnedIfConfigFileDoesNotExist(t *testing.T) {
	config := NewConfig("/nonsensepath")
	err := config.load()
	assert.NotNil(t, err)
}

func TestThatAppDataDirPathDefaultsToXDG_DATA_HOMEIfItIsSet(t *testing.T) {
	err := os.Setenv("XDG_DATA_HOME", "some path")
	if err != nil {
		log.Fatal(err)
	}
	config := NewConfig("")
	err = config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join("some path", "wirwl")
	assert.Equal(t, expectedPath, config.AppDataDirPath)
	err = os.Unsetenv("XDG_DATA_HOME")
	if err != nil {
		log.Fatal(err)
	}
}

func TestThatAppDirDefaultsToLocalShareIfXDG_DATA_HOMEIsNotSet(t *testing.T) {
	config := NewConfig("")
	err := config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join(currentUser.HomeDir, ".local", "share", "wirwl")
	assert.Equal(t, expectedPath, config.AppDataDirPath)
}

func TestThatConfigDirDefaultsToUserConfigDir(t *testing.T) {
	userConfigDirPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	config := NewConfig("")
	err = config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	expectedPath := filepath.Join(userConfigDirPath, "wirwl")
	assert.Equal(t, expectedPath, config.ConfigDirPath)
}

func TestThatConfigGetsLoadedIfItExists(t *testing.T) {
	createCorrectWirwlConfigFileForLoadingInPath(testConfigDirPath)
	config := NewConfig(testConfigDirPath)
	config.load()
	assert.Equal(t, "some db path", config.AppDataDirPath)
	assert.Equal(t, testConfigDirPath, config.ConfigDirPath)
}

func TestThatConfigFilePathGetterReturnsCorrectPath(t *testing.T) {
	config := NewConfig(testConfigDirPath)
	actualPath := config.ConfigFilePath()
	expectedPath := filepath.Join(config.ConfigDirPath, appName+".cfg")
	assert.Equal(t, expectedPath, actualPath)
}

func TestThatDefaultConfigHasCorrectKeymap(t *testing.T) {
	config := NewConfig(testConfigDirPath)
	err := config.loadDefaults()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyL), config.Keymap[input.SelectNextTabAction])
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyH), config.Keymap[input.SelectPreviousTabAction])
	assert.Equal(t, input.TwoKeyCombination(fyne.KeyS, fyne.KeyY), config.Keymap[input.SaveChangesAction])
	assert.Equal(t, input.TwoKeyCombination(fyne.KeyT, fyne.KeyI), config.Keymap[input.DisplayDialogForAddingNewEntryTypAction])
	assert.Equal(t, input.TwoKeyCombination(fyne.KeyT, fyne.KeyD), config.Keymap[input.RemoveEntryTypeAction])
	assert.Equal(t, input.TwoKeyCombination(fyne.KeyT, fyne.KeyE), config.Keymap[input.EditCurrentEntryTypeAction])
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyJ), config.Keymap[input.MoveDownAction])
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyK), config.Keymap[input.MoveUpAction])
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyI), config.Keymap[input.EnterInputModeAction])
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyReturn), config.Keymap[input.ConfirmAction])
	assert.Equal(t, input.SingleKeyCombination(fyne.KeyEscape), config.Keymap[input.CancelAction])

}
