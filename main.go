package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/azumads/selenium/app"
)

func main() {
	mux := http.NewServeMux()
	app.Admin.MountTo("/admin", mux)

	for _, path := range []string{"system", "javascripts", "stylesheets", "images"} {
		mux.Handle(fmt.Sprintf("/%s/", path), http.FileServer(http.Dir("public")))
	}

	RunAutoTest()

	fmt.Printf("Listening on: %v\n", app.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", app.Config.Port), mux); err != nil {
		panic(err)
	}
}

const (
	STATUS_PROCESSING = "processing"
	STATUS_DONE       = "done"
)

func RunAutoTest() {
	Go(func() {
		stats := ""
		t := time.Tick(time.Minute)
		for now := range t {
			if stats == STATUS_PROCESSING {
				continue
			}
			log.Println("start processing: " + now.String())
			stats = STATUS_PROCESSING
			var tests []app.ScheduledTest
			app.DB.Find(&tests)
			if len(tests) == 0 {
				stats = STATUS_DONE
				continue
			}
			for _, test := range tests {
				if time.Now().Before(test.NextRun) {
					continue
				}
				log.Println("run job id: " + test.JobId)
				response, _ := http.Post(app.HostUrl()+"/admin/workers/"+test.JobId+"/run", "application/x-www-form-urlencoded", nil)
				if response == nil || response.Body == nil {
					continue
				}
				if response.StatusCode >= 300 {
					log.Println(response.Status)
					continue
				}
				response.Body.Close()
				app.DB.Unscoped().Delete(&test)
			}
			stats = STATUS_DONE
			log.Println("done ")
		}
	})
}

func Go(f func()) {
	go func() {
		defer GoRoutineRecover()
		f()
	}()
}

func GoRoutineRecover() {
	if err := recover(); err != nil {
		stack := stack(3)
		log.Printf("Panic recovery from goroutine -> %s\n%s\n", err, stack)
		// Airbrake.Notify(fmt.Sprintf("Panic recovery -> %s\n%s\n", err, stack), nil)

	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
