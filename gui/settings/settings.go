package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/schoeneBiene/g0chat/state"
)

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func SetThemeVariant(variant fyne.ThemeVariant) {
    fyne.CurrentApp().Settings().SetTheme(&forcedVariant{ Theme: theme.DefaultTheme(), variant: variant })
}

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
}

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

    appearanceTab := container.NewVBox();

    themeSelector := widget.NewRadioGroup([]string{"Light", "Dark"}, func(s string) {
        if(s == "Light") {
            SetThemeVariant(theme.VariantLight)
            prefs.SetString("theme", "light");
        } else {
            SetThemeVariant(theme.VariantDark)
            prefs.SetString("theme", "dark");
        }
    })
    
    if(prefs.StringWithFallback("theme", "dark") == "light") {
        themeSelector.SetSelected("Light");
    } else {
        themeSelector.SetSelected("Dark");
    }

    themeSelector.Required = true;

    appearanceTab.Add(themeSelector);

    content := container.NewAppTabs(
        container.NewTabItemWithIcon("General", theme.SettingsIcon(), generalTab),
        container.NewTabItemWithIcon("Appearance", theme.ColorChromaticIcon(), appearanceTab),
    )
    
    window.SetContent(content);

    return window;
}
