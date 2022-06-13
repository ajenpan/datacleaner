package reader

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type GormReader struct {
	DB *gorm.DB

	sy    *sync.WaitGroup
	cache chan map[string]interface{}
}

//TODO:
func NewSqlserverReader(dsn string) (*GormReader, error) {
	fmt.Println(os.Getenv("GODEBUG"))

	if dsn == "" {
		dsn = "sqlserver://sa:123@127.0.0.1?database=tt" //master
	}
	dbc, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		DisableNestedTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}
	ret := &GormReader{
		DB:    dbc,
		sy:    &sync.WaitGroup{},
		cache: make(chan map[string]interface{}, 1000),
	}

	return ret, nil
}

func (w *GormReader) Read() ([]byte, error) {
	// w.cache <- line
	return nil, nil
}

func (w *GormReader) Close() error {
	close(w.cache)

	w.sy.Wait()
	return nil
}
