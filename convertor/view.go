package convertor

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel/encoder"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
	"path/filepath"
)

type ViewContract interface {
	Main(
		runConvert func(setting HandleConvertSetting),
		formats encoder.ConvertorFormatsContract,
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

type HandleConvertSetting struct {
	VideoFileInput       kernel.File
	DirectoryForSave     string
	OverwriteOutputFiles bool
	Format               string
	Encoder              encoder2.EncoderContract
}

type enableFormConversionStruct struct {
	fileVideoForConversion *widget.Button
	buttonForSelectedDir   *widget.Button
	form                   *widget.Form
}

func NewView(app kernel.AppContract) *View {
	return &View{
		app: app,
	}
}

func (v View) Main(
	runConvert func(setting HandleConvertSetting),
	formats encoder.ConvertorFormatsContract,
) {
	form := &widget.Form{}

	conversionMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	conversionMessage.TextSize = 16
	conversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	fileVideoForConversion, fileVideoForConversionMessage, fileInput := v.getButtonFileVideoForConversion(form, conversionMessage)
	buttonForSelectedDir, buttonForSelectedDirMessage, pathToSaveDirectory := v.getButtonForSelectingDirectoryForSaving()

	isOverwriteOutputFiles := false
	checkboxOverwriteOutputFilesTitle := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "checkboxOverwriteOutputFilesTitle",
	})
	checkboxOverwriteOutputFiles := widget.NewCheck(checkboxOverwriteOutputFilesTitle, func(b bool) {
		isOverwriteOutputFiles = b
	})

	selectEncoder := v.getSelectFormat(formats)

	form.Items = []*widget.FormItem{
		{
			Text:   v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "fileForConversionTitle"}),
			Widget: fileVideoForConversion,
		},
		{
			Widget: container.NewHScroll(fileVideoForConversionMessage),
		},
		{
			Text:   v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "buttonForSelectedDirTitle"}),
			Widget: buttonForSelectedDir,
		},
		{
			Widget: container.NewHScroll(buttonForSelectedDirMessage),
		},
		{
			Widget: checkboxOverwriteOutputFiles,
		},
		{
			Widget: selectEncoder.SelectFileType,
		},
		{
			Text:   v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "selectFormat"}),
			Widget: selectEncoder.SelectFormat,
		},
		{
			Text:   v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "selectEncoder"}),
			Widget: selectEncoder.SelectEncoder,
		},
	}
	form.SubmitText = v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "converterVideoFilesSubmitTitle",
	})

	enableFormConversionStruct := enableFormConversionStruct{
		fileVideoForConversion: fileVideoForConversion,
		buttonForSelectedDir:   buttonForSelectedDir,
		form:                   form,
	}

	form.OnSubmit = func() {
		if len(*pathToSaveDirectory) == 0 {
			showConversionMessage(conversionMessage, errors.New(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorSelectedFolderSave",
			})))
			enableFormConversion(enableFormConversionStruct)
			return
		}
		if len(selectEncoder.SelectFormat.Selected) == 0 {
			showConversionMessage(conversionMessage, errors.New(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorSelectedFormat",
			})))
			return
		}
		if selectEncoder.Encoder == nil {
			showConversionMessage(conversionMessage, errors.New(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorSelectedEncoder",
			})))
			return
		}
		conversionMessage.Text = ""

		fileVideoForConversion.Disable()
		buttonForSelectedDir.Disable()
		form.Disable()

		setting := HandleConvertSetting{
			VideoFileInput:       *fileInput,
			DirectoryForSave:     *pathToSaveDirectory,
			OverwriteOutputFiles: isOverwriteOutputFiles,
			Format:               selectEncoder.SelectFormat.Selected,
			Encoder:              selectEncoder.Encoder,
		}
		runConvert(setting)
		enableFormConversion(enableFormConversionStruct)

		fileVideoForConversionMessage.Text = ""
		form.Disable()
	}

	converterVideoFilesTitle := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "converterVideoFilesTitle",
	})
	v.app.GetWindow().SetContent(widget.NewCard(converterVideoFilesTitle, "", container.NewVBox(form, conversionMessage)))
	form.Disable()
}

func (v View) getButtonFileVideoForConversion(form *widget.Form, conversionMessage *canvas.Text) (*widget.Button, *canvas.Text, *kernel.File) {
	fileInput := &kernel.File{}

	fileVideoForConversionMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	fileVideoForConversionMessage.TextSize = 16
	fileVideoForConversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	buttonTitle := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	var locationURI fyne.ListableURI

	button := widget.NewButton(buttonTitle, func() {
		v.app.GetWindow().NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				fileVideoForConversionMessage.Text = err.Error()
				setStringErrorStyle(fileVideoForConversionMessage)
				return
			}
			if r == nil {
				return
			}

			fileInput.Path = r.URI().Path()
			fileInput.Name = r.URI().Name()
			fileInput.Ext = r.URI().Extension()

			fileVideoForConversionMessage.Text = r.URI().Path()
			setStringSuccessStyle(fileVideoForConversionMessage)

			form.Enable()
			conversionMessage.Text = ""

			listableURI := storage.NewFileURI(filepath.Dir(r.URI().Path()))
			locationURI, err = storage.ListerForURI(listableURI)
		}, locationURI)
	})

	return button, fileVideoForConversionMessage, fileInput
}

func (v View) getButtonForSelectingDirectoryForSaving() (button *widget.Button, buttonMessage *canvas.Text, dirPath *string) {
	buttonMessage = canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	buttonMessage.TextSize = 16
	buttonMessage.TextStyle = fyne.TextStyle{Bold: true}

	path := ""
	dirPath = &path

	buttonTitle := v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	var locationURI fyne.ListableURI

	button = widget.NewButton(buttonTitle, func() {
		v.app.GetWindow().NewFolderOpen(func(r fyne.ListableURI, err error) {
			if err != nil {
				buttonMessage.Text = err.Error()
				setStringErrorStyle(buttonMessage)
				return
			}
			if r == nil {
				return
			}

			path = r.Path()

			buttonMessage.Text = r.Path()
			setStringSuccessStyle(buttonMessage)
			locationURI, _ = storage.ListerForURI(r)

		}, locationURI)
	})

	return
}

type selectEncoder struct {
	SelectFileType *widget.RadioGroup
	SelectFormat   *widget.Select
	SelectEncoder  *widget.Select
	Encoder        encoder2.EncoderContract
}

func (v View) getSelectFormat(formats encoder.ConvertorFormatsContract) *selectEncoder {
	selectEncoder := &selectEncoder{}

	encoders := map[int]encoder2.EncoderDataContract{}
	selectEncoder.SelectEncoder = widget.NewSelect([]string{}, func(s string) {
		if encoders[selectEncoder.SelectEncoder.SelectedIndex()] == nil {
			return
		}
		selectEncoder.Encoder = encoders[selectEncoder.SelectEncoder.SelectedIndex()].NewEncoder()
	})

	formatSelected := ""
	selectEncoder.SelectFormat = widget.NewSelect([]string{}, func(s string) {
		if formatSelected == s {
			return
		}
		formatSelected = s
		format, err := formats.GetFormat(s)
		if err != nil {
			return
		}
		encoderOptions := []string{}
		encoders = format.GetEncoders()
		for _, e := range encoders {
			encoderOptions = append(encoderOptions, v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "encoder_" + e.GetTitle()}))
		}
		selectEncoder.SelectEncoder.SetOptions(encoderOptions)
		selectEncoder.SelectEncoder.SetSelectedIndex(0)
	})

	fileTypeOptions := []string{}
	for _, fileType := range encoder2.GetListFileType() {
		fileTypeOptions = append(fileTypeOptions, fileType.Name())
	}
	selectEncoder.SelectFileType = widget.NewRadioGroup([]string{"video", "audio", "image"}, func(s string) {
		formatOptions := []string{}
		for _, f := range formats.GetFormats() {
			if s != f.GetFileType().Name() {
				continue
			}
			formatOptions = append(formatOptions, f.GetTitle())
		}
		selectEncoder.SelectFormat.SetOptions(formatOptions)
		if s == encoder2.FileType(encoder2.Video).Name() {
			selectEncoder.SelectFormat.SetSelected("mp4")
		} else {
			selectEncoder.SelectFormat.SetSelectedIndex(0)
		}
	})
	selectEncoder.SelectFileType.Horizontal = true
	selectEncoder.SelectFileType.SetSelected("video")

	return selectEncoder
}

func setStringErrorStyle(text *canvas.Text) {
	text.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	text.Refresh()
}

func setStringSuccessStyle(text *canvas.Text) {
	text.Color = color.RGBA{R: 49, G: 127, B: 114, A: 255}
	text.Refresh()
}

func showConversionMessage(conversionMessage *canvas.Text, err error) {
	conversionMessage.Text = err.Error()
	setStringErrorStyle(conversionMessage)
}

func enableFormConversion(enableFormConversionStruct enableFormConversionStruct) {
	enableFormConversionStruct.fileVideoForConversion.Enable()
	enableFormConversionStruct.buttonForSelectedDir.Enable()
	enableFormConversionStruct.form.Enable()
}
