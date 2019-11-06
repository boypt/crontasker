package crontasker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type Task struct {
	CronSpec string
	LastTime time.Duration
	Cmd      string
	Args     []string
}

func (t *Task) Run() {

	deadline := time.Now().Add(t.LastTime)

	for {
		log.Printf("Run %s args %v, endat %v, now %v\n", t.Cmd, t.Args, deadline, time.Now())
		ctx, cancel := context.WithDeadline(context.Background(), deadline)

		cmd := exec.Command(t.Cmd, t.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Printf("Task failed to start, %v", err)
			cancel()
			time.Sleep(time.Second)
			continue
		}

		go func() {
			<-ctx.Done()
			log.Printf("ctx canceled, calling killing")
			if cmd.Process != nil {
				cmd.Process.Kill()
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

	log.Printf("Task end")
}

func RunRtTask(line string) (*Task, error) {
	cset := strings.Split(line, "|")
	if len(cset) != 3 {
		return nil, fmt.Errorf("line len split error, %s", line)
	}

	last, err := time.ParseDuration(cset[1])
	if err != nil {
		return nil, fmt.Errorf("line duration split error, %s", cset[1])
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

func CronDaemon(conf string) {
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	log.Printf("config: %s", conf)
	cf, err := os.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cf)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		line := strings.Trim(scanner.Text(), " \r\n")
		if line == "" {
			continue
		}
		task, err := RunRtTask(line)
		if err != nil {
			log.Printf("Task failed to create, %v", err)
			continue
		}
		c.AddJob(task.CronSpec, task)
	}
	cf.Close()

	for _, en := range c.Entries() {
		log.Println(en)
	}

	c.Run()
}
