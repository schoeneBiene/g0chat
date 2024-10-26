package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	LoginGui "github.com/schoeneBiene/g0chat/gui/login"
	MainScreenGui "github.com/schoeneBiene/g0chat/gui/mainscreen"
	State "github.com/schoeneBiene/g0chat/state"
	Socket "github.com/schoeneBiene/g0chat/ws"
)

func main() {
	app := app.NewWithID("me.goodbee.g0chat");

    token := app.Preferences().StringWithFallback("token", "");
    username := app.Preferences().StringWithFallback("username", "");

    window := app.NewWindow("G0Chat");
    
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

    window.Resize(fyne.NewSize(640, 460))
    window.ShowAndRun();
}
