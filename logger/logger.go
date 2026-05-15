package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	bold   = "\033[1m"
	dim    = "\033[2m"
)

type ProgressItem struct {
	ID    string
	Name  string
	Total int
}

type ProblemState struct {
	ID      string
	Name    string
	Total   int
	Current int
	Failed  int
	Status  string // "waiting", "uploading", "finished", "partial", "failed"
}

type Logger struct {
	mu       sync.Mutex
	updateCh chan struct{}
	states   []*ProblemState
	done     chan struct{}
	wg       sync.WaitGroup
}

func New() *Logger {
	return &Logger{
		updateCh: make(chan struct{}, 100),
		done:     make(chan struct{}),
	}
}

func (l *Logger) flush(data string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprint(os.Stdout, data)
}

func (l *Logger) Fatal(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprint(os.Stderr, red+"fatal: "+reset+msg+"\n")
	os.Exit(1)
}

func (l *Logger) Header(server, mode, source string) {
	l.flush(
		bold + "╔═══════════════════════════════════════════════════════════════╗" + reset + "\n" +
			bold + "║                 DOMjudge Test Case Uploader                   ║" + reset + "\n" +
			bold + "╚═══════════════════════════════════════════════════════════════╝" + reset + "\n" +
			dim + "  Server  : " + reset + server + "\n" +
			dim + "  Mode    : " + reset + mode + "\n" +
			dim + "  Source  : " + reset + source + "\n" +
			"\n",
	)
}

func (l *Logger) Info(msg string) {
	l.flush(msg + "\n")
}

func (l *Logger) Success(msg string) {
	l.flush(green + "  + " + reset + msg + "\n")
}

func (l *Logger) Warn(msg string) {
	l.flush(yellow + "  ? " + reset + msg + "\n")
}

func (l *Logger) InitProgress(items []ProgressItem) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.states = make([]*ProblemState, len(items))
	for i, p := range items {
		l.states[i] = &ProblemState{
			ID:     p.ID,
			Name:   p.Name,
			Total:  p.Total,
			Status: "waiting",
		}
	}
}

func (l *Logger) StartRenderer() {
	l.wg.Add(1)
	go l.renderLoop()
}

func (l *Logger) renderLoop() {
	defer l.wg.Done()
	for {
		select {
		case <-l.updateCh:
			l.drawUpdate()
		case <-l.done:
			return
		}
	}
}

func (l *Logger) DrawInitial() {
	l.mu.Lock()
	defer l.mu.Unlock()
	var b strings.Builder
	l.buildBlock(&b)
	fmt.Fprint(os.Stdout, b.String())
}

func (l *Logger) drawUpdate() {
	l.mu.Lock()
	defer l.mu.Unlock()

	linesCount := 3 + (len(l.states) * 2)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\033[%dA", linesCount))
	b.WriteString("\033[J")

	l.buildBlock(&b)

	fmt.Fprint(os.Stdout, b.String())
}

func (l *Logger) buildBlock(b *strings.Builder) {
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("\n")

	for _, s := range l.states {
		b.WriteString(fmt.Sprintf("  %s[%s][%s]%s\n", bold, s.Name, s.ID, reset))

		var status string
		switch s.Status {
		case "waiting":
			status = dim + "  Waiting..." + reset
		case "uploading":
			status = "  Uploading test case " + yellow + fmt.Sprintf("%d/%d", s.Current, s.Total) + reset
		case "finished":
			status = green + "  ALL FINISHED (" + fmt.Sprintf("%d/%d)", s.Total, s.Total) + reset
		case "partial":
			status = yellow + "  PARTIAL (" + fmt.Sprintf("%d OK, %d FAILED)", s.Total-s.Failed, s.Failed) + reset
		case "failed":
			status = red + "  ALL FAILED (" + fmt.Sprintf("%d/%d)", s.Failed, s.Total) + reset
		}
		b.WriteString(status + "\n")
	}

	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func (l *Logger) UpdateProgress(probID string, current int, status string, failed int) {
	l.mu.Lock()
	for _, s := range l.states {
		if s.ID == probID {
			s.Current = current
			s.Status = status
			s.Failed = failed
			break
		}
	}
	l.mu.Unlock()

	select {
	case l.updateCh <- struct{}{}:
	default:
	}
}

func (l *Logger) StopRenderer() {
	close(l.done)
	l.wg.Wait()
}

func (l *Logger) Interrupt(msg string) {
	l.StopRenderer()
	l.flush("\n" + yellow + "  ! " + reset + msg + "\n")
}

func (l *Logger) SummaryStart() {
	l.flush(
		"\n" +
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
			bold + "  SUMMARY" + reset + "\n" +
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
	)
}

func (l *Logger) SummaryItem(name, id string, total, failed int) {
	prefix := "  + "
	color := green
	if failed > 0 && failed < total {
		prefix = "  ! "
		color = yellow
	} else if failed == total {
		prefix = "  - "
		color = red
	}
	l.flush(fmt.Sprintf("%s%s[%s][%s] %d/%d uploaded%s\n", color, prefix, name, id, total-failed, total, reset))
}

func (l *Logger) SummaryStats(totalProbs, probsFailed, totalTC, totalFailed int) {
	l.flush(
		"\n" +
			"──────────────────────────────────────────────────────────────────\n" +
			"\n" +
			fmt.Sprintf("%s  Problems  : %s%d total, %d OK, %d with failures\n", dim, reset, totalProbs, totalProbs-probsFailed, probsFailed) +
			fmt.Sprintf("%s  Test Cases: %s%d uploaded, %d failed\n", dim, reset, totalTC-totalFailed, totalFailed) +
			"\n" +
			"──────────────────────────────────────────────────────────────────\n",
	)
}

func (l *Logger) SummaryEnd(totalFailed int) {
	if totalFailed == 0 {
		l.flush("\n" + green + "  All uploads completed successfully." + reset + "\n\n")
	} else {
		l.flush("\n" + yellow + "  Some uploads failed. Review output above." + reset + "\n\n")
	}
}
