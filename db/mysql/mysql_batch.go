package mysql

import (
	"reflect"
	"server/util"
	"time"

	"github.com/liangdas/mqant/log"
	"gorm.io/gorm"
)

var (
	WriteMap  = make(map[reflect.Type][]interface{})
	WriteChan = make(chan *WriteData, 10000)
	QueueFlag = false
)

type WriteData struct {
	ID   string
	Type reflect.Type
	Data interface{}
}

func InitBatchQueue(mysql *gorm.DB) {
	if QueueFlag {
		return
	}
	QueueFlag = true
	util.Go(func() {
		t := time.NewTimer(time.Second)
		for {
			select {
			case <-t.C:
				t.Reset(time.Second)
				for t, v := range WriteMap {
					if len(v) == 0 {
						delete(WriteMap, t)
						continue
					}
					st := reflect.SliceOf(t)
					s := reflect.MakeSlice(st, 0, 0)
					for _, one := range v {
						s = reflect.Append(s, reflect.ValueOf(one))
					}
					tableName := reflect.ValueOf(v[0]).MethodByName("TableName").Call(nil)[0].String()
					err := mysql.Table(tableName).Create(s.Interface()).Error
					if err != nil {
						log.Error("err:%v", err)
					}
					delete(WriteMap, t)
				}
			case one := <-WriteChan:
				m, ok := WriteMap[one.Type]
				if ok {
					m = append(m, one.Data)
				} else {
					m = []interface{}{one.Data}
				}
				WriteMap[one.Type] = m
			}
			// fmt.Println(total)
		}
	})
}

func AddBatch(w *WriteData) {
	WriteChan <- w
}

// 向数据库中新增数据 批量写入
func (u *MysqlClient) InsertToMysqlGo(data interface{}) {
	// util.IndexTry(func() error {
	// 	_, err := ES.client.Index().Index(index).BodyJson(params).Do(context.Background())
	// 	if err != nil {
	// 		log.Error("index:%v,write:%v,err:%v", index, params, err)
	// 		return err
	// 	}
	// 	return nil
	// })
	util.Go(func() {
		ref := reflect.TypeOf(data)
		AddBatch(&WriteData{Type: ref, Data: data})
	})
}
