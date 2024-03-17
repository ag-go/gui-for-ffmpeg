package mjpeg

import encoder2 "git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"

func NewEncoder() encoder2.EncoderContract {
	parameters := map[string]encoder2.ParameterContract{}
	getParams := func(parameters map[string]encoder2.ParameterContract) []string {
		return []string{"-c:v", "mjpeg"}
	}

	return encoder2.NewEncoder("mjpeg", parameters, getParams)
}

func NewData() encoder2.EncoderDataContract {
	title := "mjpeg"
	formats := []string{"jpg"}
	fileType := encoder2.FileType(encoder2.Image)
	return encoder2.NewData(title, formats, fileType, NewEncoder)
}
