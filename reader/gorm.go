package reader

import (
	"fmt"
	"io"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"datacleaner/object"
)

type GormReader struct {
	DB    *gorm.DB
	stop  chan bool
	wg    *sync.WaitGroup
	cache chan object.Object

	kerr   error
	rwLock sync.RWMutex
}

func NewSqlserverReader(dsn string) (*GormReader, error) {
	if dsn == "" {
		dsn = "sqlserver://sa:123@127.0.0.1?database=master"
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

	batch := 1000

	w := &GormReader{
		DB:    dbc,
		wg:    &sync.WaitGroup{},
		stop:  make(chan bool, 1),
		cache: make(chan object.Object, batch),
	}
	offset := 0
	results := make([]map[string]interface{}, 0, batch)

	w.wg.Add(1)
	go func() {
		c := uint64(0)
		defer func() {
			close(w.cache)
			fmt.Println("reader closed: ", c)
			w.wg.Done()
		}()

		for {
			res := w.DB.Debug().Raw("SELECT Name, CtfId, Gender, Address, Mobile, EMail FROM [dbo].[cdsgus] where id >= ? and id < ? and CtfTp='ID'", offset, offset+batch).Scan(&results)
			if res.Error != nil {
				w.rwLock.Lock()
				w.kerr = res.Error
				w.rwLock.Unlock()
				return
			}

			if res.RowsAffected == 0 {
				return
			}

			offset += batch
			for _, v := range results {
				select {
				case <-w.stop:
					w.rwLock.Lock()
					w.kerr = fmt.Errorf("reader closed")
					w.rwLock.Unlock()
					return
				case w.cache <- v:
					c++
				}
			}
			results = results[:0]
		}
	}()

	return w, nil
}

func (w *GormReader) Read() (object.Object, error) {
	w.rwLock.RLock()
	defer w.rwLock.RUnlock()
	if w.kerr != nil {
		return nil, w.kerr
	}

	obj, ok := <-w.cache
	if ok {
		return obj, nil
	} else {
		return nil, io.EOF
	}
}

func (w *GormReader) Close() error {
	w.stop <- true
	close(w.stop)
	w.wg.Wait()
	return nil
}
