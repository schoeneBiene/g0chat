package widgets

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CustomEntry struct {
	widget.Entry
	OnReturn func()
    LastShiftPress int64
}

func (p *CustomEntry) TypedKey(ev *fyne.KeyEvent) {
    if(ev.Name == "LeftShift" || ev.Name == "RightShift") {
        p.LastShiftPress = time.Now().UnixNano() / int64(time.Millisecond);
    }
    
    timeNow := time.Now().UnixNano() / int64(time.Millisecond);

        
    if(ev.Name == fyne.KeyReturn) {
        if(timeNow - p.LastShiftPress <= 200) {
            p.Entry.TypedKey(ev);
            return;
        }

        p.OnReturn()
    } else {
        p.Entry.TypedKey(ev);
    }
}

func NewCustomEntry() *CustomEntry {
    p := &CustomEntry{}
    p.ExtendBaseWidget(p)
    return p
}
