package loginGui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

func MakeGuestLogin(submit func(name string)) fyne.CanvasObject {
    name := widget.NewEntry();

    form := &widget.Form{
        Items: []*widget.FormItem{
            { Text: "Name", Widget: name },
        },

        OnSubmit: func() {
            submit(name.Text);
        },
    }

    return form;
}

func MakeRegisteredLogin(submit func(email string, password string)) fyne.CanvasObject {
    email := widget.NewEntry();
    email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "Not a valid email");

    password := widget.NewPasswordEntry();

    form := &widget.Form{
        Items: []*widget.FormItem{
            { Text: "Email", Widget: email },
            { Text: "Password", Widget: password },
        },

        OnSubmit: func() {
            submit(email.Text, password.Text)
        },
    }

    return form;
}
