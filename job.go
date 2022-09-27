package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"datacleaner/filter"
	"datacleaner/object"
	"datacleaner/reader"
	"datacleaner/utils"
	"datacleaner/writer"
)

type job struct {
	Name           string   `json:"name"`
	Targets        []string `json:"targets"`
	TotalSize      int64    `json:"total_size"`
	TotalHumanSize string   `json:"total_human_size"`

	NoPareserTargets []string `json:"no_pareser_targets"`

	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
	FinishAt string `json:"finish_at"`

	Mates []*FileMate `json:"mates"`
}

func JobPath(name string) string {
	return filepath.Join("./.jobs", name+".json")
}

func NewJob(name string, targets []string) (*job, error) {
	jobpath := JobPath(name)
	exist, err := utils.FileExists(jobpath)

	if exist {
		raw, err := os.ReadFile(jobpath)
		if err != nil {
			return nil, err
		}
		j := &job{}
		err = json.Unmarshal(raw, &j)
		return j, err
	}
	if err != nil {
		return nil, err
	}

	j := &job{
		Name:     name,
		Targets:  targets,
		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdateAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	for _, v := range targets {
		mate, err := genFileMate(v)
		if err != nil {
			return nil, err
		}
		j.Mates = append(j.Mates, mate)

		if len(mate.BestParser) == 0 {
			j.NoPareserTargets = append(j.NoPareserTargets, v)
		}

		j.TotalSize += mate.FileSize
	}

	j.TotalHumanSize = utils.ByteCountIEC(j.TotalSize)

	if err := j.Store(); err != nil {
		fmt.Println("store job error:", err)
	}

	return j, nil
}

func DoJob(jobname string) error {
	return nil
}

func (j *job) Store() error {
	j.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
	jobpath := JobPath(j.Name)
	raw, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(jobpath, raw, 0644)
	if err != nil {
		fmt.Println("store job error:", err)
	}
	return err
}

func (j *job) Run() error {
	defer func() {
		j.Store()
	}()

	w, err := writer.NewEsWriter("accountpasswd", nil)
	// w, err := writer.NewGormWriter("")
	if err != nil {
		return err
	}
	defer w.Close()

	doMete := func(v *FileMate) error {
		r, err := reader.NewByLine(v.FilePath)
		if err != nil {
			return err
		}
		defer r.Close()

		c := filter.AllParser[v.BestParser]
		if c == nil {
			return fmt.Errorf("parser %s not found", v.BestParser)
		}

		readCount := int64(1)
		parseCount := int64(1)
		writeCount := int64(1)

		startAt := time.Now()
		var lErr error

		reader := make(chan object.Object, 1000)
		writer := make(chan object.Object, 1000)

		sy := &sync.WaitGroup{}
		sy.Add(3)

		go func() {
			defer sy.Done()
			defer close(reader)

			for {
				line, err := r.Read()
				if err != nil {
					lErr = err
					break
				}
				if line == nil {
					break
				}

				readCount++
				reader <- line
			}
		}()

		go func() {
			defer sy.Done()
			defer close(writer)

			for line := range reader {
				if len(line) == 0 {
					continue
				}
				res, ok := c.Do(line)
				if ok && res != nil {
					parseCount++
					writer <- res
				}
			}
		}()

		go func() {
			defer sy.Done()

			for res := range writer {

				if werr := w.Write(res); werr != nil {
					fmt.Println("write err: ", err)
				} else {
					writeCount++
				}
			}
		}()

		// for {
		// 	startC++
		// 	line, err := r.Read()
		// 	if err != nil {
		// 		lErr = err
		// 		break
		// 	}
		// 	if line == nil {
		// 		lErr = nil
		// 		break
		// 	}
		// 	res, ok := c.Parse(string(line))
		// 	if !ok {
		// 		// fmt.Println("parse error,line: ", string(line))
		// 		continue
		// 	}
		// 	werr := w.Write(res)
		// 	if werr != nil {
		// 		fmt.Println("write err: ", err)
		// 	} else {
		// 		finishC++
		// 	}
		// }

		sy.Wait()

		cost := time.Since(startAt)
		ns := float64(writeCount) / (cost.Seconds())
		fmt.Printf("%s finish, r:%d, c:%d, w:%d, cost:%v, %v/s, size:%v\n", v.FilePath, readCount, parseCount, writeCount, cost.Seconds(), ns, v.FileHumanSize)
		return lErr
	}

	sy := &sync.WaitGroup{}

	for i, v := range j.Mates {
		v.StartAt = time.Now().Format("2006-01-02 15:04:05")

		fmt.Printf("%s [%d/%d]: %s %s\n", v.StartAt, i+1, len(j.Mates), v.FilePath, v.FileHumanSize)

		if len(v.Stat) != 0 {
			continue
		}

		if len(v.BestParser) == 0 {
			continue
		}

		sy.Add(1)
		go func(i int, v *FileMate) {
			defer sy.Done()

			err := doMete(v)
			if err == io.EOF {
				err = nil
			}

			if err != nil {
				v.Stat = "error"
				v.Msg = err.Error()
			} else {
				v.Msg = "ok"
				v.Stat = "done"
			}
			v.FinishAt = time.Now().Format("2006-01-02 15:04:05")
		}(i, v)
		// j.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
	}

	sy.Wait()
	j.FinishAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}
