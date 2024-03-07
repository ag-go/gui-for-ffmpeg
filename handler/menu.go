package handler

import (
	"fyne.io/fyne/v2"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/menu"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type MenuHandlerContract interface {
	GetMainMenu() *fyne.MainMenu
	LanguageSelection()
}

type MenuHandler struct {
	app                 kernel.AppContract
	convertorHandler    ConvertorHandlerContract
	menuView            menu.ViewContract
	localizerView       localizer.ViewContract
	localizerRepository localizer.RepositoryContract
}

func NewMenuHandler(
	app kernel.AppContract,
	convertorHandler ConvertorHandlerContract,
	menuView menu.ViewContract,
	localizerView localizer.ViewContract,
	localizerRepository localizer.RepositoryContract,
) *MenuHandler {
	return &MenuHandler{
		app:                 app,
		convertorHandler:    convertorHandler,
		menuView:            menuView,
		localizerView:       localizerView,
		localizerRepository: localizerRepository,
	}
}

func (h MenuHandler) GetMainMenu() *fyne.MainMenu {
	settings := h.getMenuSettings()
	help := h.getMenuHelp()

	return fyne.NewMainMenu(settings, help)
}

func (h MenuHandler) getMenuSettings() *fyne.Menu {
	quit := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "exit",
	}), nil)
	quit.IsQuit = true
	h.app.GetLocalizerService().AddChangeCallback("exit", func(text string) {
		quit.Label = text
	})

	languageSelection := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "changeLanguage",
	}), h.LanguageSelection)
	h.app.GetLocalizerService().AddChangeCallback("changeLanguage", func(text string) {
		languageSelection.Label = text
	})

	ffPathSelection := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "changeFFPath",
	}), h.convertorHandler.FfPathSelection)
	h.app.GetLocalizerService().AddChangeCallback("changeFFPath", func(text string) {
		ffPathSelection.Label = text
	})

	settings := fyne.NewMenu(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "settings",
	}), languageSelection, ffPathSelection, quit)
	h.app.GetLocalizerService().AddChangeCallback("settings", func(text string) {
		settings.Label = text
		settings.Refresh()
	})

	return settings
}

func (h MenuHandler) getMenuHelp() *fyne.Menu {
	about := fyne.NewMenuItem(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "about",
	}), h.openAbout)
	h.app.GetLocalizerService().AddChangeCallback("about", func(text string) {
		about.Label = text
	})

	help := fyne.NewMenu(h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "help",
	}), about)
	h.app.GetLocalizerService().AddChangeCallback("help", func(text string) {
		help.Label = text
		help.Refresh()
	})

	return help
}

func (h MenuHandler) openAbout() {
	ffmpeg, err := h.convertorHandler.GetFfmpegVersion()
	if err != nil {
		ffmpeg = h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFmpegVersion",
		})
	}
	ffprobe, err := h.convertorHandler.GetFfprobeVersion()
	if err != nil {
		ffprobe = h.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
			MessageID: "errorFFprobeVersion",
		})
	}

	h.menuView.About(ffmpeg, ffprobe)
}

func (h MenuHandler) LanguageSelection() {
	h.localizerView.LanguageSelection(func(lang kernel.Lang) {
		_, _ = h.localizerRepository.Save(lang.Code)
		h.convertorHandler.MainConvertor()
	})
}
