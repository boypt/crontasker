package crontasker

import (
	"testing"
	"time"
)

func TestTask_Run(t *testing.T) {
	type fields struct {
		EndDu time.Duration
		Cmd   string
		Args  []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"1", fields{time.Second * -1, "/bin/sleep", []string{"25"}},
		},
		{
			"2", fields{time.Second * 3, "/bin/sleep", []string{"2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zt := &Task{
				LastTime: tt.fields.EndDu,
				Cmd:      tt.fields.Cmd,
				Args:     tt.fields.Args,
			}
			zt.Run()
		})
	}
}

func TestTask_runOnce(t *testing.T) {
	type fields struct {
		CronSpec string
		LastTime time.Duration
		Cmd      string
		Args     []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{"1", fields{"* * * * *", 0, "ls", []string{"/"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zt := &Task{
				CronSpec: tt.fields.CronSpec,
				LastTime: tt.fields.LastTime,
				Cmd:      tt.fields.Cmd,
				Args:     tt.fields.Args,
			}
			zt.runOnce()
		})
	}
}
