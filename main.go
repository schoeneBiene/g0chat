package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	LoginGui "github.com/schoeneBiene/g0chat/gui/login"
    MainScreenGui "github.com/schoeneBiene/g0chat/gui/mainscreen"
    Socket "github.com/schoeneBiene/g0chat/ws"
    State "github.com/schoeneBiene/g0chat/state"
)

func main() {
	app := app.NewWithID("me.goodbee.g0chat");

    window := app.NewWindow("G0Chat");

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
    window.ShowAndRun();
}
