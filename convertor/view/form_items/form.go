package form_items

import (
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor/view/form_items/h264_nvenc"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor/view/form_items/libx264"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor/view/form_items/libx265"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
)

var Views = map[string]func(encoder encoder.EncoderContract, app kernel.AppContract) []*widget.FormItem{
	"libx264":    libx264.View,
	"h264_nvenc": h264_nvenc.View,
	"libx265":    libx265.View,
}
