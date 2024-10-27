package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type RightclickableLabel struct {
	widget.Label
    OnRightClick func(*fyne.PointEvent)
}

func (label *RightclickableLabel) Tapped(_ *fyne.PointEvent) {}
func (label *RightclickableLabel) TappedSecondary(e *fyne.PointEvent) {
    label.OnRightClick(e);
}

func NewRightclickableLabel(text string) *RightclickableLabel {
    label := &RightclickableLabel{};
    label.ExtendBaseWidget(label);
    label.SetText(text);
    
    return label;
}
