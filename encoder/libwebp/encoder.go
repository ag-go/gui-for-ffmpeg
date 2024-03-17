package libwebp

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "libwebp"}
	}

	return encoder2.NewEncoder("libwebp", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "libwebp"
	formats := []string{"webp"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
