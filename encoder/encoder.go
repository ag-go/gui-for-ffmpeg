package encoder

import "errors"

type EncoderContract interface {
	GetName() string
	GetParams() []string
	GetParameter(name string) (ParameterContract, error)
}

type ParameterContract interface {
	GetName() string
	Set(string) error
	Get() string
	IsEnabled() bool
	SetEnable()
	SetDisable()
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

type Encoder struct {
	name       string
	parameters map[string]ParameterContract
	getParams  func(parameters map[string]ParameterContract) []string
}

func NewEncoder(name string, parameters map[string]ParameterContract, getParams func(parameters map[string]ParameterContract) []string) *Encoder {
	return &Encoder{
		name:       name,
		parameters: parameters,
		getParams:  getParams,
	}
}

func (e *Encoder) GetName() string {
	return e.name
}

func (e *Encoder) GetParams() []string {
	return e.getParams(e.parameters)
}

func (e *Encoder) GetParameter(name string) (ParameterContract, error) {
	if e.parameters[name] == nil {
		return nil, errors.New("parameter not found")
	}

	return e.parameters[name], nil
}

type Parameter struct {
	name         string
	isEnabled    bool
	parameter    string
	setParameter func(string) (string, error)
}

func NewParameter(name string, isEnabled bool, defaultParameter string, setParameter func(string) (string, error)) *Parameter {
	return &Parameter{
		name:         name,
		isEnabled:    isEnabled,
		parameter:    defaultParameter,
		setParameter: setParameter,
	}
}

func (p *Parameter) GetName() string {
	return p.name
}

func (p *Parameter) Set(s string) (err error) {
	if p.setParameter != nil {
		s, err = p.setParameter(s)
		if err != nil {
			return err
		}
	}
	p.parameter = s
	return nil
}

func (p *Parameter) Get() string {
	return p.parameter
}

func (p *Parameter) IsEnabled() bool {
	return p.isEnabled
}

func (p *Parameter) SetEnable() {
	p.isEnabled = true
}

func (p *Parameter) SetDisable() {
	p.isEnabled = false
}
