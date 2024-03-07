package encoder

type EncoderContract interface {
	GetParams() []string
}

type EncoderDataContract interface {
	GetTitle() string
	GetFormats() []string
	GetFileType() FileTypeContract
	NewEncoder() EncoderContract
}

type Data struct {
	title    string
	formats  []string
	fileType FileTypeContract
	encoder  func() EncoderContract
}

func NewData(title string, formats []string, fileType FileTypeContract, encoder func() EncoderContract) *Data {
	return &Data{
		title:    title,
		formats:  formats,
		fileType: fileType,
		encoder:  encoder,
	}
}

func (data Data) GetTitle() string {
	return data.title
}

func (data Data) GetFormats() []string {
	return data.formats
}

func (data Data) NewEncoder() EncoderContract {
	return data.encoder()
}

func (data Data) GetFileType() FileTypeContract {
	return data.fileType
}

type FileTypeContract interface {
	Name() string
	Ordinal() int
}

const (
	Video = iota
	Audio
	Image
)

type FileType uint

var fileTypeStrings = []string{
	"video",
	"audio",
	"image",
}

func (fileType FileType) Name() string {
	return fileTypeStrings[fileType]
}

func (fileType FileType) Ordinal() int {
	return int(fileType)
}

func GetListFileType() []FileTypeContract {
	return []FileTypeContract{
		FileType(Video),
		FileType(Audio),
		FileType(Image),
	}
}
