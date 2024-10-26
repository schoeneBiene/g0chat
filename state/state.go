package state

import "fyne.io/fyne/v2"

var Login_Anon bool
var Login_Username string
var Login_Email string
var Login_Password string
var Login_Token = ""

var SendMessage func(content string)

var MainWindow fyne.Window;

var Debug_WS bool;
