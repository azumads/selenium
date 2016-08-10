package app

import (
	"bytes"
	"os/exec"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/qor/admin"
	"github.com/qor/media_library"
	"github.com/qor/worker"
)

var Loops = [][]string{
	{"", "Never"},
	{"1", "1H"},
	{"2", "2H"},
	{"3", "3H"},
	{"6", "6H"},
	{"12", "12H"},
	{"24", "24H"},
}

type AutoTest struct {
	gorm.Model
	Name     string
	JobId    string
	LoopHour string
	NextRun  time.Time
}

type AutoTestingArgument struct {
	Name     string
	Loop     string
	TestFile media_library.FileSystem
	CsvFile  media_library.FileSystem
}

const ONCE_WRITE_COUNT = 1000

func AddWorker() *worker.Worker {
	Worker := worker.New()

	autoTestingRes := Admin.NewResource(&AutoTestingArgument{})
	autoTestingRes.Meta(&admin.Meta{
		Name:   "Loop",
		Config: &admin.SelectOneConfig{Collection: Loops},
	})

	Worker.RegisterJob(&worker.Job{
		Name:     "Auto Testing",
		Resource: autoTestingRes,
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) (err error) {
			AutoTestingArgument := arg.(*AutoTestingArgument)
			loop := AutoTestingArgument.Loop
			intloop, _ := strconv.ParseInt(loop, 10, 0)
			if intloop != 0 {
				DB.Create(&AutoTest{
					Name:     AutoTestingArgument.Name,
					JobId:    qorJob.GetJobID(),
					LoopHour: loop,
					NextRun:  time.Now().Add(time.Duration(intloop) * time.Hour)})
			}
			// defer os.Remove(AutoTestingArgument.TestFile.URL())
			// qorJob.AddLog("./bang.py " + path.Join("public", AutoTestingArgument.TestFile.URL()))
			out1, err1 := run("./bang.py", []string{path.Join("public", AutoTestingArgument.TestFile.URL())})
			if err1 != nil {
				qorJob.AddLog(err1.Error())
				qorJob.AddLog(out1)
				err = err1
				return
			}
			qorJob.AddLog(strings.Trim(out1, "\n"))
			out2, err2 := run(strings.Trim(out1, "\n"), nil)
			if err2 != nil {
				qorJob.AddLog(err2.Error())
				err = err2
			}
			qorJob.AddLog(out2)
			return
		},
	})

	Admin.AddResource(Worker)
	return Worker
}

func run(command string, args []string) (out string, err error) {
	var buf bytes.Buffer
	var cmd *exec.Cmd
	if len(args) == 0 {
		cmd = &exec.Cmd{
			Path: command,
			Args: []string{command},
		}
	} else {
		cmd = exec.Command(command, args...)
	}
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err = cmd.Run()
	if err != nil {
		out = buf.String()
		cmd.Process.Release()
		buf.Reset()
		return
	}

	out = buf.String()

	// Clean up resource
	cmd.Process.Kill()
	buf.Reset()
	debug.FreeOSMemory()

	return
}
