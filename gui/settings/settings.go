package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/schoeneBiene/g0chat/state"
)

func MakeSettingsWindow() fyne.Window {
    window := fyne.CurrentApp().NewWindow("Settings");
    prefs := fyne.CurrentApp().Preferences();

    generalTab := container.NewVBox();
    generalTab.Add(widget.NewButtonWithIcon("Logout", theme.LogoutIcon(), func() {
        prefs.SetString("username", "");
        prefs.SetString("token", "");

        window.Close();
        state.MainWindow.Close();
    }))

    content := container.NewAppTabs(
        container.NewTabItem("General", generalTab),
    )
    
    window.SetContent(content);

    return window;
}
