package wirwl

import (
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisplayingGUI(t *testing.T) {
	app := NewApp()
	app.LoadAndDisplayGUI(test.NewApp())
	assert.Equal(t, app.currentTabLabel.Text, "No data loaded")
}
