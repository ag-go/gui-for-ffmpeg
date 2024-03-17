package convertor

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/convertor/view"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
)

type ViewContract interface {
	Main(
		formConversion view.ConversionContract,
	)
	SelectFFPath(
		ffmpegPath string,
		ffprobePath string,
		save func(ffmpegPath string, ffprobePath string) error,
		cancel func(),
		donwloadFFmpeg func(progressBar *widget.ProgressBar, progressMessage *canvas.Text) error,
	)
}

type View struct {
	app kernel.AppContract
}

func NewView(app kernel.AppContract) *View {
	return &View{
		app: app,
	}
}

func (v View) Main(formConversion view.ConversionContract) {
	converterVideoFilesTitle := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "converterVideoFilesTitle",
	})
	v.app.GetWindow().SetContent(widget.NewCard(converterVideoFilesTitle, "", container.NewVScroll(formConversion.GetContent())))
	formConversion.AfterViewContent()
}

func setStringErrorStyle(text *canvas.Text) {
	text.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	text.Refresh()
}

func setStringSuccessStyle(text *canvas.Text) {
	text.Color = color.RGBA{R: 49, G: 127, B: 114, A: 255}
	text.Refresh()
}
