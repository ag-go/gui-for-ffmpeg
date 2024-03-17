package libx264

import (
	"errors"
	encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
)

var Presets = []string{
	"ultrafast",
	"superfast",
	"veryfast",
	"faster",
	"fast",
	"medium",
	"slow",
	"slower",
	"veryslow",
	"placebo",
}

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{
		"preset": newParameterPreset(),
	}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		params := []string{"-c:v", "libx264"}

		if parameters["preset"] != nil && parameters["preset"].IsEnabled() {
			params = append(params, "-preset", parameters["preset"].Get())
		}

		return params
	}

	return encoder2.NewEncoder("libx264", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "libx264"
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
	return encoder2.NewParameter("preset", false, "medium", setParameter)
}
