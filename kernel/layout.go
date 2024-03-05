package kernel

import (
	"bufio"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type LayoutContract interface {
	SetContent(content fyne.CanvasObject) *fyne.Container
	NewProgressbar(queueId int, totalDuration float64) ProgressContract
	ChangeQueueStatus(queueId int, queue *Queue)
}

type Layout struct {
	layout            *fyne.Container
	queueLayoutObject QueueLayoutObjectContract
	localizerService  LocalizerContract
}

func NewLayout(queueLayoutObject QueueLayoutObjectContract, localizerService LocalizerContract) *Layout {
	layout := container.NewAdaptiveGrid(2, widget.NewLabel(""), container.NewVScroll(queueLayoutObject.GetCanvasObject()))

	return &Layout{
		layout:            layout,
		queueLayoutObject: queueLayoutObject,
		localizerService:  localizerService,
	}
}

func (l Layout) SetContent(content fyne.CanvasObject) *fyne.Container {
	l.layout.Objects[0] = content
	return l.layout
}

func (l Layout) NewProgressbar(queueId int, totalDuration float64) ProgressContract {
	progressbar := l.queueLayoutObject.GetProgressbar(queueId)
	return NewProgress(totalDuration, progressbar, l.localizerService)
}

func (l Layout) ChangeQueueStatus(queueId int, queue *Queue) {
	l.queueLayoutObject.ChangeQueueStatus(queueId, queue)
}

type QueueLayoutObjectContract interface {
	GetCanvasObject() fyne.CanvasObject
	GetProgressbar(queueId int) *widget.ProgressBar
	ChangeQueueStatus(queueId int, queue *Queue)
}

type QueueLayoutObject struct {
	QueueListContract QueueListContract

	queue                 QueueListContract
	container             *fyne.Container
	items                 map[int]QueueLayoutItem
	localizerService      LocalizerContract
	queueStatisticsFormat *queueStatisticsFormat
}

type QueueLayoutItem struct {
	CanvasObject  fyne.CanvasObject
	ProgressBar   *widget.ProgressBar
	StatusMessage *canvas.Text
	MessageError  *canvas.Text

	status *StatusContract
}

func NewQueueLayoutObject(queue QueueListContract, localizerService LocalizerContract) *QueueLayoutObject {
	title := widget.NewLabel(localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: "queue"}))
	title.TextStyle.Bold = true

	localizerService.AddChangeCallback("queue", func(text string) {
		title.Text = text
		title.Refresh()
	})

	items := map[int]QueueLayoutItem{}
	queueStatisticsFormat := newQueueStatisticsFormat(localizerService, &items)

	queueLayoutObject := &QueueLayoutObject{
		queue: queue,
		container: container.NewVBox(
			container.NewHBox(title, queueStatisticsFormat.completed.widget, queueStatisticsFormat.error.widget),
			container.NewHBox(queueStatisticsFormat.inProgress.widget, queueStatisticsFormat.waiting.widget, queueStatisticsFormat.total.widget),
		),
		items:                 items,
		localizerService:      localizerService,
		queueStatisticsFormat: queueStatisticsFormat,
	}

	queue.AddListener(queueLayoutObject)

	return queueLayoutObject
}

func (o QueueLayoutObject) GetCanvasObject() fyne.CanvasObject {
	return o.container
}

func (o QueueLayoutObject) GetProgressbar(queueId int) *widget.ProgressBar {
	if item, ok := o.items[queueId]; ok {
		return item.ProgressBar
	}
	return widget.NewProgressBar()
}

func (o QueueLayoutObject) Add(id int, queue *Queue) {
	progressBar := widget.NewProgressBar()
	statusMessage := canvas.NewText(o.getStatusTitle(queue.Status), theme.PrimaryColor())
	messageError := canvas.NewText("", theme.ErrorColor())

	content := container.NewVBox(
		container.NewHScroll(widget.NewLabel(queue.Setting.VideoFileInput.Name)),
		progressBar,
		container.NewHScroll(statusMessage),
		container.NewHScroll(messageError),
		canvas.NewLine(theme.FocusColor()),
		container.NewPadded(),
	)

	o.queueStatisticsFormat.addQueue()
	if o.queueStatisticsFormat.isChecked(queue.Status) == false {
		content.Hide()
	}

	o.items[id] = QueueLayoutItem{
		CanvasObject:  content,
		ProgressBar:   progressBar,
		StatusMessage: statusMessage,
		MessageError:  messageError,
		status:        &queue.Status,
	}
	o.container.Add(content)
}

func (o QueueLayoutObject) Remove(id int) {
	if item, ok := o.items[id]; ok {
		o.container.Remove(item.CanvasObject)
		o.queueStatisticsFormat.removeQueue(*item.status)
		o.items[id] = QueueLayoutItem{}
	}
}

func (o QueueLayoutObject) ChangeQueueStatus(queueId int, queue *Queue) {
	if item, ok := o.items[queueId]; ok {
		statusColor := o.getStatusColor(queue.Status)
		item.StatusMessage.Text = o.getStatusTitle(queue.Status)
		item.StatusMessage.Color = statusColor
		item.StatusMessage.Refresh()
		if queue.Error != nil {
			item.MessageError.Text = queue.Error.Error()
			item.MessageError.Color = statusColor
			item.MessageError.Refresh()
		}
		if o.queueStatisticsFormat.isChecked(queue.Status) == false && item.CanvasObject.Visible() == true {
			item.CanvasObject.Hide()
		} else if item.CanvasObject.Visible() == false {
			item.CanvasObject.Show()
		}
		o.queueStatisticsFormat.changeQueue(queue.Status)
	}
}

func (o QueueLayoutObject) getStatusColor(status StatusContract) color.Color {
	if status == StatusType(Error) {
		return theme.ErrorColor()
	}

	if status == StatusType(Completed) {
		return color.RGBA{R: 49, G: 127, B: 114, A: 255}
	}

	return theme.PrimaryColor()
}

func (o QueueLayoutObject) getStatusTitle(status StatusContract) string {
	return o.localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: status.Name() + "Queue"})
}

type Progress struct {
	totalDuration    float64
	progressbar      *widget.ProgressBar
	protocol         string
	localizerService LocalizerContract
}

func NewProgress(totalDuration float64, progressbar *widget.ProgressBar, localizerService LocalizerContract) Progress {
	return Progress{
		totalDuration:    totalDuration,
		progressbar:      progressbar,
		protocol:         "pipe:",
		localizerService: localizerService,
	}
}

func (p Progress) GetProtocole() string {
	return p.protocol
}

func (p Progress) Run(stdOut io.ReadCloser, stdErr io.ReadCloser) error {
	isProcessCompleted := false
	var errorText string

	p.progressbar.Value = 0
	p.progressbar.Max = p.totalDuration
	p.progressbar.Refresh()

	progress := 0.0

	go func() {
		scannerErr := bufio.NewReader(stdErr)
		for {
			line, _, err := scannerErr.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}
			data := strings.TrimSpace(string(line))
			errorText = data
		}
	}()

	scannerOut := bufio.NewReader(stdOut)
	for {
		line, _, err := scannerOut.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		data := strings.TrimSpace(string(line))
		if strings.Contains(data, "progress=end") {
			p.progressbar.Value = p.totalDuration
			p.progressbar.Refresh()
			isProcessCompleted = true
			break
		}

		re := regexp.MustCompile(`frame=(\d+)`)
		a := re.FindAllStringSubmatch(data, -1)

		if len(a) > 0 && len(a[len(a)-1]) > 0 {
			c, err := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
			if err != nil {
				continue
			}
			progress = float64(c)
		}
		if p.progressbar.Value != progress {
			p.progressbar.Value = progress
			p.progressbar.Refresh()
		}
	}

	if isProcessCompleted == false {
		if len(errorText) == 0 {
			errorText = p.localizerService.GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorConverter",
			})
		}
		return errors.New(errorText)
	}

	return nil
}

type queueStatistics struct {
	widget *widget.Check
	title  string
	count  *int64
}
type queueStatisticsFormat struct {
	waiting    *queueStatistics
	inProgress *queueStatistics
	completed  *queueStatistics
	error      *queueStatistics
	total      *queueStatistics
}

func newQueueStatisticsFormat(localizerService LocalizerContract, queueItems *map[int]QueueLayoutItem) *queueStatisticsFormat {
	checkWaiting := newQueueStatistics("waitingQueue", localizerService)
	checkInProgress := newQueueStatistics("inProgressQueue", localizerService)
	checkCompleted := newQueueStatistics("completedQueue", localizerService)
	checkError := newQueueStatistics("errorQueue", localizerService)
	checkTotal := newQueueStatistics("total", localizerService)

	queueStatisticsFormat := &queueStatisticsFormat{
		waiting:    checkWaiting,
		inProgress: checkInProgress,
		completed:  checkCompleted,
		error:      checkError,
		total:      checkTotal,
	}

	checkTotal.widget.OnChanged = func(b bool) {
		if b == true {
			queueStatisticsFormat.allCheckboxChecked()
		} else {
			queueStatisticsFormat.allUnCheckboxChecked()
		}
		queueStatisticsFormat.redrawingQueueItems(queueItems)
	}

	queueStatisticsFormat.waiting.widget.OnChanged = func(b bool) {
		if b == true {
			queueStatisticsFormat.checkboxChecked()
		} else {
			queueStatisticsFormat.unCheckboxChecked()
		}
		queueStatisticsFormat.redrawingQueueItems(queueItems)
	}

	queueStatisticsFormat.inProgress.widget.OnChanged = func(b bool) {
		if b == true {
			queueStatisticsFormat.checkboxChecked()
		} else {
			queueStatisticsFormat.unCheckboxChecked()
		}
		queueStatisticsFormat.redrawingQueueItems(queueItems)
	}

	queueStatisticsFormat.completed.widget.OnChanged = func(b bool) {
		if b == true {
			queueStatisticsFormat.checkboxChecked()
		} else {
			queueStatisticsFormat.unCheckboxChecked()
		}
		queueStatisticsFormat.redrawingQueueItems(queueItems)
	}

	queueStatisticsFormat.error.widget.OnChanged = func(b bool) {
		if b == true {
			queueStatisticsFormat.checkboxChecked()
		} else {
			queueStatisticsFormat.unCheckboxChecked()
		}
		queueStatisticsFormat.redrawingQueueItems(queueItems)
	}

	return queueStatisticsFormat
}

func (f queueStatisticsFormat) redrawingQueueItems(queueItems *map[int]QueueLayoutItem) {
	for _, item := range *queueItems {
		if f.isChecked(*item.status) == true && item.CanvasObject.Visible() == false {
			item.CanvasObject.Show()
			continue
		}
		if f.isChecked(*item.status) == false && item.CanvasObject.Visible() == true {
			item.CanvasObject.Hide()
		}
	}
}

func (f queueStatisticsFormat) isChecked(status StatusContract) bool {
	if status == StatusType(InProgress) {
		return f.inProgress.widget.Checked
	}
	if status == StatusType(Completed) {
		return f.completed.widget.Checked
	}
	if status == StatusType(Error) {
		return f.error.widget.Checked
	}
	if status == StatusType(Waiting) {
		return f.waiting.widget.Checked
	}

	return true
}

func (f queueStatisticsFormat) addQueue() {
	f.waiting.add()
	f.total.add()
}

func (f queueStatisticsFormat) changeQueue(status StatusContract) {
	if status == StatusType(InProgress) {
		f.waiting.remove()
		f.inProgress.add()
		return
	}

	if status == StatusType(Completed) {
		f.inProgress.remove()
		f.completed.add()
		return
	}

	if status == StatusType(Error) {
		f.inProgress.remove()
		f.error.add()
		return
	}
}

func (f queueStatisticsFormat) removeQueue(status StatusContract) {
	f.total.remove()

	if status == StatusType(Completed) {
		f.completed.remove()
		return
	}

	if status == StatusType(Error) {
		f.error.remove()
		return
	}

	if status == StatusType(InProgress) {
		f.inProgress.remove()
		return
	}

	if status == StatusType(Waiting) {
		f.waiting.remove()
		return
	}
}

func (f queueStatisticsFormat) checkboxChecked() {
	if f.total.widget.Checked == true {
		return
	}

	if f.waiting.widget.Checked == false {
		return
	}

	if f.inProgress.widget.Checked == false {
		return
	}

	if f.completed.widget.Checked == false {
		return
	}

	if f.error.widget.Checked == false {
		return
	}

	f.total.widget.Checked = true
	f.total.widget.Refresh()
}

func (f queueStatisticsFormat) unCheckboxChecked() {
	if f.total.widget.Checked == false {
		return
	}

	f.total.widget.Checked = false
	f.total.widget.Refresh()
}

func (f queueStatisticsFormat) allCheckboxChecked() {
	f.waiting.widget.Checked = true
	f.waiting.widget.Refresh()
	f.inProgress.widget.Checked = true
	f.inProgress.widget.Refresh()
	f.completed.widget.Checked = true
	f.completed.widget.Refresh()
	f.error.widget.Checked = true
	f.error.widget.Refresh()
}

func (f queueStatisticsFormat) allUnCheckboxChecked() {
	f.waiting.widget.Checked = false
	f.waiting.widget.Refresh()
	f.inProgress.widget.Checked = false
	f.inProgress.widget.Refresh()
	f.completed.widget.Checked = false
	f.completed.widget.Refresh()
	f.error.widget.Checked = false
	f.error.widget.Refresh()
}

func newQueueStatistics(messaigeID string, localizerService LocalizerContract) *queueStatistics {
	checkbox := widget.NewCheck("", nil)
	checkbox.Checked = true

	count := int64(0)

	title := localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: messaigeID})
	queueStatistics := &queueStatistics{
		widget: checkbox,
		title:  strings.ToLower(title),
		count:  &count,
	}

	queueStatistics.formatText(false)

	localizerService.AddChangeCallback(messaigeID, func(text string) {
		queueStatistics.title = strings.ToLower(text)
		queueStatistics.formatText(true)
		queueStatistics.widget.Refresh()
	})

	return queueStatistics
}

func (s queueStatistics) add() {
	*s.count += 1
	s.formatText(true)
}

func (s queueStatistics) remove() {
	if *s.count == 0 {
		return
	}
	*s.count -= 1
	s.formatText(true)
}

func (s queueStatistics) formatText(refresh bool) {
	s.widget.Text = s.title + ": " + strconv.FormatInt(*s.count, 10)
	if refresh == true {
		s.widget.Refresh()
	}
}
