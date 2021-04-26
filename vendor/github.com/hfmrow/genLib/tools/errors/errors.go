// errors.go

package errors

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// WarnCallerMess: return information about caller that can be interpreted and
// clickable inside LiteIde. Output format: "./main.go:114 Warning!,"
func WarnCallerMess(skip ...int) string {
	var s int = 2
	if len(skip) > 0 {
		s = s + skip[0]
	}
	_, file, line, _ := runtime.Caller(s) // error callback type handling
	return fmt.Sprintf("./%s:%d Warning!,", filepath.Base(file), line)
}

// Check: Display errors in convenient way with stack display
// "options" set to true -> exit on error.
// NOTICE: exit option must not be used if a "defer" function
// is initiated, otherwise, defer will never be applied !
func Check(err error, options ...bool) (isError bool) {
	if err != nil {
		type errorInf struct {
			fn string // function
			f  string // file
		}
		var stacked []errorInf
		var outStrErr string
		isError = true
		stack := strings.Split(string(debug.Stack()), "\n")
		for errIdx := 5; errIdx < len(stack)-1; errIdx++ {
			stacked = append(stacked, errorInf{fn: stack[errIdx], f: strings.TrimSpace(stack[errIdx+1])})
			errIdx++
		}
		baseMessage := strings.Split(err.Error(), "\n\n")
		for _, mess := range baseMessage {
			outStrErr += fmt.Sprintf("[%s]\n", mess)
		}
		for errIdx := 1; errIdx < len(stacked); errIdx++ {
			outStrErr += fmt.Sprintf("[%s]*[%s]\n", strings.SplitN(stacked[errIdx].fn, "(", 2)[0], stacked[errIdx].f)
		}
		fmt.Print(outStrErr)
		if len(options) > 0 {
			if options[0] {
				os.Exit(1)
			}
		}
	}
	return
}

// Get all information on runtime from callers functions
// include memory, processor, Goroutine ...
type FuncTracer struct {
	filename string
	file     *os.File
	writer   *bufio.Writer
	frames   *runtime.Frames
	delay    time.Duration
	lastNow  time.Time
	FiltersInclude,
	FiltersExclude []string
}

func FuncTracerNew(filename string, delay float64) (*FuncTracer, error) {
	var err error
	ft := new(FuncTracer)

	ft.delay = time.Duration(float64(time.Second) * delay)
	ft.lastNow = time.Now()

	ft.FiltersExclude = []string{"/usr/local/go/src/"}

	if ft.file, err = os.Create(filename); err == nil {
		ft.writer = bufio.NewWriter(ft.file)
		err = ft.Write(fmt.Sprintf("Time;Alloc;SysMem;NumGc;nCgoCall;nGortn;Cpu%%;Rss;Line;Func;File;\n"))
	}
	return ft, err
}

func (ft *FuncTracer) memStats() *runtime.MemStats {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m
}

func (ft *FuncTracer) Close() (err error) {
	if err = ft.writer.Flush(); err == nil {
		err = ft.file.Close()
	}
	return
}

func (ft *FuncTracer) Write(line string) error {
	count, err := ft.writer.WriteString(line)
	if count != len(line) {
		err = fmt.Errorf("[Missing written bytes], Recieved bytes: %d, Written: %d", len(line), count)
	}
	return err
}

func (ft *FuncTracer) match(name string) bool {
	for _, excl := range ft.FiltersExclude {
		if strings.Contains(name, excl) {
			return true
		}
	}
	return false
}

// DebugFunc: Retrieve information about the caller function
// include filename, function name and line number.
type DebugFunc struct {
	Frames *runtime.Frames

	File,
	Caller string
	Line,

	Skip int
}

func DebugFuncNew(skip ...int) (d *DebugFunc) {
	d = new(DebugFunc)
	d.Skip = 1
	if len(skip) > 0 {
		d.Skip = skip[0]
	}
	d.details()
	return
}

func (d *DebugFunc) details() (file, fnct string, line int, ok bool) {

	var pc uintptr
	pc, d.File, d.Line, ok = runtime.Caller(d.Skip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		d.Caller = details.Name()
	}
	return d.File, d.Caller, d.Line, ok
}

func (d *DebugFunc) GetCurrentFunctionName() string {
	// Skip GetCurrentFunctionName
	return d.getFrame(d.Skip).Function
}

func (d *DebugFunc) GetCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return d.getFrame(d.Skip + 1).Function
}

func (d *DebugFunc) getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		d.Frames = runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = d.Frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

var nonAlNum = regexp.MustCompile(`[[:punct:][:alpha:]]`)

// Get current timestamp
func timeStamp() string {
	date := nonAlNum.ReplaceAllString(time.Now().Format(time.RFC3339), "")[:14]
	return date[:8] + "-" + date[8:]
}

// humanReadableSize:
func humanReadableSize(size interface{}) string {

	return fmt.Sprintf("%.2fMiB", float64(size.(uint64))/1024e3)
}
