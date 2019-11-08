package crontasker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	debug bool
)

func debugf(format string, args ...interface{}) {
	if debug {
		log.Printf("[crontasker]DEBUG "+format, args...)
	}
}

type Task struct {
	CronSpec string
	LastTime time.Duration
	Cmd      string
	Args     []string
}

func (t *Task) Run() {
	if t.LastTime == 0 {
		t.runOnce()
	} else {
		t.runWithDeadline()
	}
}

func (t *Task) runOnce() {
	debugf("Run %s args %v\n", t.Cmd, t.Args)
	cmd := exec.Command(t.Cmd, t.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Task failed with err, %v", err)
	}

	debugf("Task end")
}

func (t *Task) runWithDeadline() {

	deadline := time.Now().Add(t.LastTime)

	for {
		debugf("Run %s args %v, endat %v, now %v\n", t.Cmd, t.Args, deadline, time.Now())
		ctx, cancel := context.WithDeadline(context.Background(), deadline)

		cmd := exec.Command(t.Cmd, t.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if err := cmd.Start(); err != nil {
			log.Printf("Task failed to start, %v", err)
			cancel()
			time.Sleep(time.Second)
			continue
		}

		go func() {
			<-ctx.Done()
			debugf("ctx done")
			if cmd.Process != nil {
				// kill Process Group
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
					log.Printf("calling kill error %v", err)
				}
			}
			return
		}()

		if err := cmd.Wait(); err != nil {
			log.Printf("Task Process ended, err: %v", err)
		}
		cancel()

		if time.Now().Before(deadline) {
			log.Printf("Task restart")
			continue
		}

		break
	}

	debugf("Task end")
}

func parseTask(line string) (*Task, error) {
	cset := strings.Split(line, "|")
	if len(cset) != 3 {
		return nil, fmt.Errorf("line len split error, %s", line)
	}

	last, err := time.ParseDuration(cset[1])
	if err != nil {
		last = 0
	}

	cmds := strings.Split(cset[2], " ")
	task := &Task{
		CronSpec: cset[0],
		Cmd:      cmds[0],
		Args:     cmds[1:],
		LastTime: last,
	}
	return task, nil
}

func CronDaemon(conf string) error {
	var opts []cron.Option

	if debug {
		opts = append(opts, cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "crontasker: ", log.LstdFlags))))
	}

	c := cron.New(opts...)

	log.Printf("config: %s", conf)
	cf, err := os.Open(conf)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cf)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		line := strings.Trim(scanner.Text(), " \r\n")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		task, err := parseTask(line)
		if err != nil {
			log.Printf("Task failed to create, %v", err)
			continue
		}
		c.AddJob(task.CronSpec, task)
	}
	cf.Close()

	jobs := len(c.Entries())
	if jobs > 0 {
		log.Printf("found %d cron jobs", jobs)
		c.Run()
	} else {
		log.Printf("no task prased, exit")
	}
	return nil
}
