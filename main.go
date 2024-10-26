package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	LoginGui "github.com/schoeneBiene/g0chat/gui/login"
	MainScreenGui "github.com/schoeneBiene/g0chat/gui/mainscreen"
	"github.com/schoeneBiene/g0chat/gui/settings"
	State "github.com/schoeneBiene/g0chat/state"
	Socket "github.com/schoeneBiene/g0chat/ws"
)

func main() {
	app := app.NewWithID("me.goodbee.g0chat");

    token := app.Preferences().StringWithFallback("token", "");
    username := app.Preferences().StringWithFallback("username", "");

    window := app.NewWindow("G0Chat");

    // load settings
    prefs := app.Preferences();
    
    if(prefs.StringWithFallback("theme", "dark") == "light") {
        settings.SetThemeVariant(theme.VariantLight)
    } else {
        settings.SetThemeVariant(theme.VariantDark);
    }

    // auth
    if(token == "" && username == "") {
        loginContent := container.NewAppTabs(
            container.NewTabItem("Guest", LoginGui.MakeGuestLogin(func(name string) {
                State.Login_Anon = true;
                State.Login_Username = name;

                window.SetContent(MainScreenGui.MakeMainScreen())
                go Socket.MakeSocketConnection();
            })),
            container.NewTabItem("Registered", LoginGui.MakeRegisteredLogin(func(email, password string) {
                State.Login_Anon = false;
                State.Login_Email = email;
                State.Login_Password = password;

                window.SetContent(MainScreenGui.MakeMainScreen());
                go Socket.MakeSocketConnection();
            })),
        )

        window.SetContent(loginContent);
    } else if(token != "") {
        State.Login_Anon = false;
        State.Login_Token = token;

        window.SetContent(MainScreenGui.MakeMainScreen());
        go Socket.MakeSocketConnection();
    } else {
        State.Login_Anon = true;
        State.Login_Username = username;

        window.SetContent(MainScreenGui.MakeMainScreen());
        go Socket.MakeSocketConnection();
    }

    window.SetCloseIntercept(func() {
        os.Exit(0);
    })

    window.Resize(fyne.NewSize(960, 690))
    State.MainWindow = window;
    window.ShowAndRun();
}
