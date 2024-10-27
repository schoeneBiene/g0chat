package mainscreen

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/schoeneBiene/g0chat/gui/settings"
	"github.com/schoeneBiene/g0chat/gui/widgets"

	State "github.com/schoeneBiene/g0chat/state"
)

type Message struct {
    authorName string
    authorId string
    authorRole string

    content string
    timstamp int64
}

var messages = []Message{};
var messageList *widget.List;

func AddMessage(authorName, authorId, authorRole, content string, timestamp int64) {
    messages = append(messages, Message{
        authorName: authorName,
        authorId: authorId,
        authorRole: authorRole,

        content: content,
        timstamp: timestamp,
    })

    messageList.Refresh();
    messageList.ScrollToBottom();
}

var members = []string{};
var memberList *widget.List;

func UpdateMembers(newMembers []string) {
    log.Println("Updating members list")

    members = newMembers;
    memberList.Refresh();
}

func MakeMainScreen() fyne.CanvasObject {
    messageList = widget.NewList(
        func() int {
            return len(messages);
        },

        func() fyne.CanvasObject {
            return widget.NewLabel("template");
        },

        func(lii widget.ListItemID, co fyne.CanvasObject) {
            data := messages[lii];
            label := co.(*widget.Label);
            
            label.SetText(fmt.Sprintf("%s%s", func() string {
                if(data.authorName == "System") {
                    return "";
                } else {
                    return fmt.Sprintf("[%s] %s: ", data.authorRole, data.authorName);
                }
            }(), data.content));
            messageList.SetItemHeight(lii, label.MinSize().Height)
        },
    )
    messageList.OnSelected = func(_ widget.ListItemID) {
        messageList.UnselectAll();
    }

    memberList = widget.NewList(
        func() int {
            return len(members);
        },

        func() fyne.CanvasObject {
            return widget.NewLabel("template");
        },

        func(lii widget.ListItemID, co fyne.CanvasObject) {
            co.(*widget.Label).SetText(members[lii]);
        },
    );
    memberList.OnSelected = func(_ widget.ListItemID) {
        memberList.UnselectAll();
    }

    messageInput := widgets.NewCustomEntry();
    messageInput.MultiLine = true;
    messageInput.SetMinRowsVisible(5);

    onSubmit := func() {
        State.SendMessage(messageInput.Text);
        messageInput.SetText("");
    }

    messageInput.OnReturn = func() {
        if(fyne.CurrentDevice().IsMobile()) { return; }

        onSubmit();
    }

    messageInputForm := &widget.Form{
        Items: []*widget.FormItem{
            { Text: "", Widget: messageInput },
        },

        OnSubmit: onSubmit,
        SubmitText: "Send",
    }

    toolbar := widget.NewToolbar(
        widget.NewToolbarAction(theme.SettingsIcon(), func() {
            settingsWindow := settings.MakeSettingsWindow();

            settingsWindow.CenterOnScreen();
            settingsWindow.Resize(fyne.NewSize(640, 460))

            settingsWindow.Show();
        }),
    )

    mainContent := container.NewBorder(toolbar, messageInputForm, nil, nil, messageList);
    content := container.NewHSplit(mainContent, memberList);
    content.SetOffset(fyne.CurrentApp().Preferences().FloatWithFallback("member_list_split", 0.9));

    go func() {
        prevOffset := content.Offset;

        for range time.Tick(5 * time.Second) {
            if(content.Offset != prevOffset) {
                prevOffset = content.Offset;

                fyne.CurrentApp().Preferences().SetFloat("member_list_split", content.Offset);
            }
        }
    }()

    return content;
}
