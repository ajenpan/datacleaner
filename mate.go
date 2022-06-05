package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"datacleaner/filter"
	"datacleaner/reader"
	"datacleaner/utility"
)

func ReadFileMate(fp string) (*FileMate, error) {
	fp, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	matePath := matePath(fp)

	//read cache mate
	raw, err := os.ReadFile(matePath)
	if err != nil {
		mate, err := genFileMate(fp)
		if err != nil {
			return nil, err
		}
		if err = cacheMate(mate); err != nil {
			fmt.Println(err)
		}
		return mate, nil
	}

	mate := &FileMate{}
	err = json.Unmarshal(raw, mate)
	return mate, err
}

type FileMate struct {
	RandomLines   []string              `json:"random_lines"`
	FilePath      string                `json:"file_path"`
	FileSize      int64                 `json:"file_size"`
	FileHumanSize string                `json:"file_human_size"`
	CreateTime    string                `json:"create_time"`
	ParserScore   []*filter.ParserScore `json:"parser_score"`
	BestParser    string                `json:"best_parser"`

	StartAt  string `json:"start_at"`
	FinishAt string `json:"finish_at"`
	Stat     string `json:"stat"`
	Msg      string `json:"msg"`

	// SuccessRecords in     `json:"success_records"`
}

func matePath(fp string) string {
	fp, err := filepath.Abs(fp)
	if err != nil {
		return ""
	}
	hash := md5.Sum([]byte(fp))
	mateName := hex.EncodeToString(hash[:]) + ".json"
	matePath := filepath.Join("./.mate/", mateName)
	return matePath
}

func cacheMate(mate *FileMate) error {
	mp := matePath(mate.FilePath)
	if len(mp) == 0 {
		return fmt.Errorf("filepath is empty")
	}
	raw, err := json.MarshalIndent(mate, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(mp, raw, 0644)
	if err != nil {
		return err
	}
	return nil
}

func genFileMate(filepath string) (*FileMate, error) {
	mate, err := reader.ReadFileInfo(filepath)
	if err != nil {
		return nil, err
	}

	ret := &FileMate{
		RandomLines:   mate.RandomLines,
		FilePath:      mate.FilePath,
		FileSize:      mate.FileSize,
		FileHumanSize: utility.ByteCountIEC(mate.FileSize),
		CreateTime:    mate.CreateTime.Format("2006-01-02 15:04:05"),
	}

	temp := filter.GetParserScore(mate.RandomLines)
	for _, v := range temp {
		if v.Score > 0 {
			ret.ParserScore = append(ret.ParserScore, v)
		}
	}

	if len(ret.ParserScore) > 0 && ret.ParserScore[0].Score > 10 {
		ret.BestParser = ret.ParserScore[0].Parser
	}

	return ret, nil
}
