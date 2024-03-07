package encoder

import (
	"errors"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/encoder"
)

type ConvertorFormatContract interface {
	GetTitle() string
	AddEncoder(encoder encoder.EncoderDataContract)
	GetFileType() encoder.FileTypeContract
	GetEncoders() map[int]encoder.EncoderDataContract
}

type ConvertorFormat struct {
	title    string
	fileType encoder.FileTypeContract
	encoders map[int]encoder.EncoderDataContract
}

func NewConvertorFormat(title string, fileType encoder.FileTypeContract) *ConvertorFormat {
	return &ConvertorFormat{
		title:    title,
		fileType: fileType,
		encoders: map[int]encoder.EncoderDataContract{},
	}
}

func (f ConvertorFormat) GetTitle() string {
	return f.title
}

func (f ConvertorFormat) AddEncoder(encoder encoder.EncoderDataContract) {
	f.encoders[len(f.encoders)] = encoder
}

func (f ConvertorFormat) GetEncoders() map[int]encoder.EncoderDataContract {
	return f.encoders
}

func (f ConvertorFormat) GetFileType() encoder.FileTypeContract {
	return f.fileType
}

type ConvertorFormatsContract interface {
	NewEncoder(encoderName string) bool
	GetFormats() map[string]ConvertorFormatContract
	GetFormat(format string) (ConvertorFormatContract, error)
}

type ConvertorFormats struct {
	formats map[string]ConvertorFormatContract
}

func NewConvertorFormats() *ConvertorFormats {
	return &ConvertorFormats{
		formats: map[string]ConvertorFormatContract{},
	}
}

func (f ConvertorFormats) NewEncoder(encoderName string) bool {
	if supportEncoders[encoderName] == nil {
		return false
	}
	data := supportEncoders[encoderName]()
	for _, format := range data.GetFormats() {
		if f.formats[format] == nil {
			f.formats[format] = NewConvertorFormat(format, data.GetFileType())
		}
		f.formats[format].AddEncoder(data)
	}
	return true
}

func (f ConvertorFormats) GetFormats() map[string]ConvertorFormatContract {
	return f.formats
}

func (f ConvertorFormats) GetFormat(format string) (ConvertorFormatContract, error) {
	if f.formats[format] == nil {
		return ConvertorFormat{}, errors.New("not found ConvertorFormat")
	}
	return f.formats[format], nil
}
