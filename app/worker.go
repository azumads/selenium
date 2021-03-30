package app

import (
	"bytes"
	"os/exec"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/qor/admin"
	"github.com/qor/worker"
)

var Loops = [][]string{
	{"0", "Never"},
	{"1", "1H"},
	{"2", "2H"},
	{"3", "3H"},
	{"6", "6H"},
	{"12", "12H"},
	{"24", "24H"},
	{"48", "48H"},
}

type RunTestArgument struct {
	Project   Project
	ProjectID uint
	Loop      string
}

const ONCE_WRITE_COUNT = 1000

func AddWorker() *worker.Worker {
	Worker := worker.New()

	autoTestingRes := Admin.NewResource(&RunTestArgument{})
	autoTestingRes.Meta(&admin.Meta{
		Name:   "Loop",
		Config: &admin.SelectOneConfig{Collection: Loops},
	})

	Worker.RegisterJob(&worker.Job{
		Name:     "Run Test",
		Resource: autoTestingRes,
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) (err error) {
			RunTestArgument := arg.(*RunTestArgument)
			testcases := []TestCase{}
			DB.Where("project_id = ?", RunTestArgument.Project.ID).Find(&testcases)
			loop := RunTestArgument.Loop
			intloop, _ := strconv.ParseInt(loop, 10, 0)
			if intloop != 0 {
				DB.Create(&ScheduledTest{
					Project:   RunTestArgument.Project,
					ProjectID: RunTestArgument.Project.ID,
					JobId:     qorJob.GetJobID(),
					LoopHour:  loop,
					NextRun:   time.Now().Add(time.Duration(intloop) * time.Hour)})
			}

			qorJob.AddLog("start to run project " + RunTestArgument.Project.Name + "'s test cases")
			for i, tc := range testcases {
				qorJob.AddLog("----------------------------------------------------------------------")
				qorJob.AddLog("start to run test case " + strconv.Itoa(i+1) + ": " + tc.Name)
				runArgs := []string{path.Join("public", tc.TestFile.URL())}
				if tc.CsvFile.URL() != "" {
					runArgs = append(runArgs, path.Join("public", tc.CsvFile.URL()))
				}
				script := "./bang.py"
				if IsProd() {
					script = "./bang_linux.py"
				}
				out1, err1 := run(script, runArgs)
				if err1 != nil {
					qorJob.AddLog(err1.Error())
					qorJob.AddLog(out1)
					err = err1
					if RunTestArgument.Project.NotifyEmail != "" {
						SendNotifyErrorEmail(RunTestArgument.Project.NotifyEmail, RunTestArgument.Project.Name, tc.Name, qorJob.GetJobID())
					}

					return
				}

				// qorJob.AddLog(strings.Trim(out1, "\n"))
				out2, err2 := run(strings.Trim(out1, "\n"), nil)
				if err2 != nil {
					qorJob.AddLog(err2.Error())
					err = err2
					qorJob.AddLog(out2)
					if RunTestArgument.Project.NotifyEmail != "" {
						SendNotifyErrorEmail(RunTestArgument.Project.NotifyEmail, RunTestArgument.Project.Name, tc.Name, qorJob.GetJobID())
					}
					return
				}
				qorJob.AddLog(strings.Trim(out2, `...
----------------------------------------------------------------------`))

			}

			return
		},
	})

	Admin.AddResource(Worker, &admin.Config{Name: "Run test"})
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
