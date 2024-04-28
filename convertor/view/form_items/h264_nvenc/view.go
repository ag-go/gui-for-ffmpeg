package h264_nvenc

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder/h264_nvenc"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func View(encoder encoder.EncoderContract, app kernel.AppContract) []*widget.FormItem {
	items := []*widget.FormItem{}

	items = append(items, presetParameter(encoder, app)...)

	return items
}

func presetParameter(encoder encoder.EncoderContract, app kernel.AppContract) []*widget.FormItem {
	parameter, err := encoder.GetParameter("preset")
	if err != nil {
		return nil
	}

	presets := map[string]string{}
	presetsForSelect := []string{}
	presetDefault := ""

	for _, name := range h264_nvenc.Presets {
		title := name
		presetsForSelect = append(presetsForSelect, name)
		presets[title] = name
		if name == parameter.Get() {
			presetDefault = title
		}
	}

	elementSelect := widget.NewSelect(presetsForSelect, func(s string) {
		if presets[s] == "" {
			return
		}
		parameter.Set(presets[s])
	})
	elementSelect.SetSelected(presetDefault)
	elementSelect.Hide()

	checkboxTitle := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "parameterCheckbox"})
	elementCheckbox := widget.NewCheck(checkboxTitle, func(b bool) {
		if b == true {
			parameter.SetEnable()
			elementSelect.Show()
			return
		}
		parameter.SetDisable()
		elementSelect.Hide()
	})

	return []*widget.FormItem{
		{
			Text:   app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "formPreset"}),
			Widget: container.NewVBox(elementCheckbox, elementSelect),
		},
	}
}
