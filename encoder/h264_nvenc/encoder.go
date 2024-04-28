package h264_nvenc

import (
	"errors"
	encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
)

var Presets = []string{
	"default",
	"slow",
	"medium",
	"fast",
	"hp",
	"hq",
	"bd",
	"ll",
	"llhq",
	"llhp",
	"lossless",
	"losslesshp",
}

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{
		"preset": newParameterPreset(),
	}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		params := []string{"-c:v", "h264_nvenc"}

		if parameters["preset"] != nil && parameters["preset"].IsEnabled() {
			params = append(params, "-preset", parameters["preset"].Get())
		}

		return params
	}

	return encoder2.NewEncoder("h264_nvenc", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "h264_nvenc"
	formats := []string{"mp4"}
	fileType := encoder2.FileType(encoder2.Video)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}

func newParameterPreset() encoder2.ParameterContract {
	setParameter := func(s string) (string, error) {
		for _, value := range Presets {
			if value == s {
				return value, nil
			}
		}

		return "", errors.New("preset not found")
	}
	return encoder2.NewParameter("preset", false, "default", setParameter)
}
