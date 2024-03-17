package libwebp_anim

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "libwebp_anim"}
	}

	return encoder2.NewEncoder("libwebp_anim", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "libwebp_anim"
	formats := []string{"webp"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
