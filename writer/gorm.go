package writer

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type GormWriter struct {
	DB *gorm.DB

	sy    *sync.WaitGroup
	cache chan map[string]interface{}
}

func NewGormWriter(dsn string) (*GormWriter, error) {
	if dsn == "" {
		dsn = "root:123456@tcp(localhost:3306)/passwd?charset=utf8mb4&parseTime=True&loc=Local"
	}
	dbc, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableNestedTransaction: true, //关闭嵌套事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Error),
	})

	if err != nil {
		return nil, err
	}
	ret := &GormWriter{
		DB:    dbc,
		sy:    &sync.WaitGroup{},
		cache: make(chan map[string]interface{}, 1000),
	}
	ret.sy.Add(1)

	go func() {
		defer ret.sy.Done()
		defer fmt.Println("writer close")

		batch := make([]map[string]interface{}, 0, 500)
		for line := range ret.cache {
			batch = append(batch, line)

			if len(batch) >= 500 {
				if err := ret.DB.Table("passwd").Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(batch, 100).Error; err != nil {
					fmt.Println(err)
				}
				batch = batch[:0]
			}
		}

		if len(batch) > 0 {
			if err := ret.DB.Table("passwd").Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(batch, 100).Error; err != nil {
				fmt.Println(err)
			}
		}

	}()

	return ret, nil
}

func (w *GormWriter) Write(line map[string]interface{}) error {
	w.cache <- line
	return nil
}

func (w *GormWriter) Close() error {
	close(w.cache)

	w.sy.Wait()
	return nil
}
