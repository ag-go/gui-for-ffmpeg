package view

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

type ConversionContract interface {
	GetContent() fyne.CanvasObject
	AfterViewContent()
}

type Conversion struct {
	app                  kernel.AppContract
	form                 *widget.Form
	conversionMessage    *canvas.Text
	fileForConversion    *fileForConversion
	directoryForSaving   *directoryForSaving
	overwriteOutputFiles *overwriteOutputFiles
	selectEncoder        *selectEncoder
	runConvert           func(setting HandleConvertSetting)
}

type HandleConvertSetting struct {
	FileInput            kernel.File
	DirectoryForSave     string
	OverwriteOutputFiles bool
	Format               string
	Encoder              encoder2.EncoderContract
}

func NewConversion(app kernel.AppContract, formats encoder.ConvertorFormatsContract, runConvert func(setting HandleConvertSetting)) *Conversion {
	conversionMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	conversionMessage.TextSize = 16
	conversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	fileForConversion := newFileForConversion(app)
	directoryForSaving := newDirectoryForSaving(app)
	overwriteOutputFiles := newOverwriteOutputFiles(app)
	selectEncoder := newSelectEncoder(app, formats)

	form := widget.NewForm()
	form.Items = []*widget.FormItem{
		{
			Text:   app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "fileForConversionTitle"}),
			Widget: fileForConversion.button,
		},
		{
			Widget: container.NewHScroll(fileForConversion.message),
		},
		{
			Text:   app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "buttonForSelectedDirTitle"}),
			Widget: directoryForSaving.button,
		},
		{
			Widget: container.NewHScroll(directoryForSaving.message),
		},
		{
			Widget: overwriteOutputFiles.checkbox,
		},
		{
			Widget: selectEncoder.SelectFileType,
		},
		{
			Text:   app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "selectFormat"}),
			Widget: selectEncoder.SelectFormat,
		},
		{
			Text:   app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "selectEncoder"}),
			Widget: selectEncoder.SelectEncoder,
		},
	}
	form.SubmitText = app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "converterVideoFilesSubmitTitle",
	})

	return &Conversion{
		app:                  app,
		form:                 form,
		conversionMessage:    conversionMessage,
		fileForConversion:    fileForConversion,
		directoryForSaving:   directoryForSaving,
		overwriteOutputFiles: overwriteOutputFiles,
		selectEncoder:        selectEncoder,
		runConvert:           runConvert,
	}
}

func (c Conversion) GetContent() fyne.CanvasObject {
	c.form.OnSubmit = c.submit
	c.fileForConversion.AddChangeCallback(c.selectFileForConversion)

	return container.NewVBox(
		c.form,
		c.conversionMessage,
	)
}

func (c Conversion) AfterViewContent() {
	c.form.Disable()
}

func (c Conversion) selectFileForConversion(err error) {
	c.conversionMessage.Text = ""
	if err != nil {
		c.form.Disable()
		return
	}

	c.form.Enable()
}

func (c Conversion) submit() {
	if len(c.directoryForSaving.path) == 0 {
		showConversionMessage(c.conversionMessage, errors.New(c.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorSelectedFolderSave",
		})))
		c.enableFormConversion()
		return
	}
	if len(c.selectEncoder.SelectFormat.Selected) == 0 {
		showConversionMessage(c.conversionMessage, errors.New(c.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorSelectedFormat",
		})))
		return
	}
	if c.selectEncoder.Encoder == nil {
		showConversionMessage(c.conversionMessage, errors.New(c.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorSelectedEncoder",
		})))
		return
	}
	c.conversionMessage.Text = ""

	c.fileForConversion.button.Disable()
	c.directoryForSaving.button.Disable()
	c.form.Disable()

	setting := HandleConvertSetting{
		FileInput:            *c.fileForConversion.file,
		DirectoryForSave:     c.directoryForSaving.path,
		OverwriteOutputFiles: c.overwriteOutputFiles.IsChecked(),
		Format:               c.selectEncoder.SelectFormat.Selected,
		Encoder:              c.selectEncoder.Encoder,
	}
	c.runConvert(setting)
	c.enableFormConversion()

	c.fileForConversion.message.Text = ""
	c.form.Disable()
}

func (c Conversion) enableFormConversion() {
	c.fileForConversion.button.Enable()
	c.directoryForSaving.button.Enable()
	c.form.Enable()
}

type fileForConversion struct {
	button  *widget.Button
	message *canvas.Text
	file    *kernel.File

	changeCallbacks map[int]func(err error)
}

func newFileForConversion(app kernel.AppContract) *fileForConversion {
	fileForConversion := &fileForConversion{
		file:            &kernel.File{},
		changeCallbacks: map[int]func(err error){},
	}

	buttonTitle := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	fileForConversion.message = canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	fileForConversion.message.TextSize = 16
	fileForConversion.message.TextStyle = fyne.TextStyle{Bold: true}

	var locationURI fyne.ListableURI

	fileForConversion.button = widget.NewButton(buttonTitle, func() {
		app.GetWindow().NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				fileForConversion.message.Text = err.Error()
				setStringErrorStyle(fileForConversion.message)
				fileForConversion.eventSelectFile(err)
				return
			}
			if r == nil {
				return
			}

			fileForConversion.file.Path = r.URI().Path()
			fileForConversion.file.Name = r.URI().Name()
			fileForConversion.file.Ext = r.URI().Extension()

			fileForConversion.message.Text = r.URI().Path()
			setStringSuccessStyle(fileForConversion.message)

			fileForConversion.eventSelectFile(nil)

			listableURI := storage.NewFileURI(filepath.Dir(r.URI().Path()))
			locationURI, err = storage.ListerForURI(listableURI)
		}, locationURI)
	})

	return fileForConversion
}

func (c fileForConversion) AddChangeCallback(callback func(err error)) {
	c.changeCallbacks[len(c.changeCallbacks)] = callback
}

func (c fileForConversion) eventSelectFile(err error) {
	for _, changeCallback := range c.changeCallbacks {
		changeCallback(err)
	}
}

type directoryForSaving struct {
	button  *widget.Button
	message *canvas.Text
	path    string
}

func newDirectoryForSaving(app kernel.AppContract) *directoryForSaving {
	directoryForSaving := &directoryForSaving{
		path: "",
	}

	directoryForSaving.message = canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	directoryForSaving.message.TextSize = 16
	directoryForSaving.message.TextStyle = fyne.TextStyle{Bold: true}

	buttonTitle := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	var locationURI fyne.ListableURI

	directoryForSaving.button = widget.NewButton(buttonTitle, func() {
		app.GetWindow().NewFolderOpen(func(r fyne.ListableURI, err error) {
			if err != nil {
				directoryForSaving.message.Text = err.Error()
				setStringErrorStyle(directoryForSaving.message)
				return
			}
			if r == nil {
				return
			}

			directoryForSaving.path = r.Path()

			directoryForSaving.message.Text = r.Path()
			setStringSuccessStyle(directoryForSaving.message)
			locationURI, _ = storage.ListerForURI(r)

		}, locationURI)
	})

	return directoryForSaving
}

type overwriteOutputFiles struct {
	checkbox  *widget.Check
	isChecked bool
}

func newOverwriteOutputFiles(app kernel.AppContract) *overwriteOutputFiles {
	overwriteOutputFiles := &overwriteOutputFiles{
		isChecked: false,
	}
	checkboxOverwriteOutputFilesTitle := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "checkboxOverwriteOutputFilesTitle",
	})
	overwriteOutputFiles.checkbox = widget.NewCheck(checkboxOverwriteOutputFilesTitle, func(b bool) {
		overwriteOutputFiles.isChecked = b
	})

	return overwriteOutputFiles
}

func (receiver overwriteOutputFiles) IsChecked() bool {
	return receiver.isChecked
}

type selectEncoder struct {
	SelectFileType *widget.RadioGroup
	SelectFormat   *widget.Select
	SelectEncoder  *widget.Select
	Encoder        encoder2.EncoderContract
}

func newSelectEncoder(app kernel.AppContract, formats encoder.ConvertorFormatsContract) *selectEncoder {
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
			encoderOptions = append(encoderOptions, app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "encoder_" + e.GetTitle()}))
		}
		selectEncoder.SelectEncoder.SetOptions(encoderOptions)
		selectEncoder.SelectEncoder.SetSelectedIndex(0)
	})

	fileTypeOptions := []string{}
	for _, fileType := range encoder2.GetListFileType() {
		fileTypeOptions = append(fileTypeOptions, fileType.Name())
	}

	encoderGroupVideo := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "encoderGroupVideo"})
	encoderGroupAudio := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "encoderGroupAudio"})
	encoderGroupImage := app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{MessageID: "encoderGroupImage"})
	encoderGroup := map[string]string{
		encoderGroupVideo: "video",
		encoderGroupAudio: "audio",
		encoderGroupImage: "image",
	}
	selectEncoder.SelectFileType = widget.NewRadioGroup([]string{encoderGroupVideo, encoderGroupAudio, encoderGroupImage}, func(s string) {
		groupCode := encoderGroup[s]

		formatOptions := []string{}
		for _, f := range formats.GetFormats() {
			if groupCode != f.GetFileType().Name() {
				continue
			}
			formatOptions = append(formatOptions, f.GetTitle())
		}
		selectEncoder.SelectFormat.SetOptions(formatOptions)
		if groupCode == encoder2.FileType(encoder2.Video).Name() {
			selectEncoder.SelectFormat.SetSelected("mp4")
		} else {
			selectEncoder.SelectFormat.SetSelectedIndex(0)
		}
	})
	selectEncoder.SelectFileType.Horizontal = true
	selectEncoder.SelectFileType.SetSelected(encoderGroupVideo)

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
