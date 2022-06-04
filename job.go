package main

import (
	"datacleaner/utility"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type job struct {
	// fileMate *FileMate

	// c cleaner.LineParser
	// r reader.Reader
	// w writer.Writer

	Name    string   `json:"name"`
	Targets []string `json:"targets"`

	NoPareserTargets []string `json:"no_pareser_targets"`

	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	FinishAt time.Time `json:"finish_at"`

	Mates []*FileMate `json:"mates"`
}

func JobPath(name string) string {
	return filepath.Join("./.jobs", name+".json")
}

func NewJob(name string, targets []string) error {
	jobpath := JobPath(name)
	exist, err := utility.FileExists(jobpath)

	if exist {
		return fmt.Errorf("job %s already exists", name)
	}
	if err != nil {
		return err
	}

	j := &job{
		Name:     name,
		Targets:  targets,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	for _, v := range targets {
		mate, err := genFileMate(v)
		if err != nil {
			return err
		}
		j.Mates = append(j.Mates, mate)

		if len(mate.BestParser) == 0 {
			j.NoPareserTargets = append(j.NoPareserTargets, v)
		}
	}

	raw, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(jobpath, raw, 0644)
}

func DoJob(jobname string) error {
	return nil
}

// func (j *job) temp() error {

// 	line, err := j.r.Read()

// 	if err != nil {
// 		return err
// 	}

// 	res, ok := j.c.Parse(string(line))

// 	if !ok {
// 		return fmt.Errorf("parse error:%s", line)
// 	}

// 	err = j.w.Write(res)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
