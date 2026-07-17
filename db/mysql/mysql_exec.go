package mysql

import (
	"fmt"
	"reflect"

	"github.com/liangdas/mqant/log"
	"gorm.io/gorm"
)

// Upsert 存在更新不存在插入
func (u *MysqlClient) Upsert(condition string, values interface{}, tx ...*gorm.DB) (int64, error) {
	if tx == nil {
		return Upsert(condition, values, u.client)
	}
	return Upsert(condition, values, tx[0])
}

// Upsert 存在更新不存在插入
func Upsert(condition string, values interface{}, d *gorm.DB) (int64, error) {
	var count int64
	if err := d.Model(values).Where(condition).Count(&count).Error; err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	if count > 0 {
		tx := d.Model(values).Where(condition).Updates(values)
		if tx.Error != nil {
			log.Error("err:%v", tx.Error)
		}
		return tx.RowsAffected, tx.Error
	}
	tx := d.Model(values).Create(values)
	if tx.Error != nil {
		log.Error("err:%v", tx.Error)
	}
	return tx.RowsAffected, tx.Error
}

func (u *MysqlClient) UpsertMap(condition string, model, update interface{}) error {
	var count int64
	if err := u.client.Model(model).Where(condition).Count(&count).Error; err != nil {
		if err != nil {
			log.Error("err:%v", err)
		}
		return err
	}
	if count > 0 {
		if err := u.client.Model(model).Where(condition).Updates(update).Error; err != nil {
			log.Error("err:%v", err)
			return err
		}
	} else {
		if err := u.client.Create(model).Error; err != nil {
			log.Error("err:%v", err)
			return err
		}
	}
	return nil
}

func (u *MysqlClient) UpdateW(model, update interface{}, condition string) error {
	if err := u.client.Model(model).Where(condition).Updates(update).Error; err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (u *MysqlClient) Update(model, update interface{}) error {
	if err := u.client.Model(model).Where(model).Updates(update).Error; err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (u *MysqlClient) Del(model interface{}, condi ...interface{}) error {
	if err := u.client.Delete(model, condi...).Error; err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (u *MysqlClient) UpdateRes(model, update interface{}, tx ...*gorm.DB) (int64, error) {
	d := u.client
	if tx != nil {
		d = tx[0]
	}
	res := d.Model(model).Where(model).Updates(update)
	return res.RowsAffected, res.Error
}

func (u *MysqlClient) UpdateWRes(model, update interface{}, condition string, tx ...*gorm.DB) (int64, error) {
	d := u.client
	if tx != nil {
		d = tx[0]
	}
	res := d.Model(model).Where(condition).Updates(update)
	err := res.Error
	if err != nil {
		log.Error("err:%v", err)
	}
	return res.RowsAffected, err
}

// SelectField 查询选择的字段，扫描进dst
func (u *MysqlClient) SelectField(dst interface{}, fields ...string) error {
	return u.client.Model(dst).Where(dst).Select(fields).Scan(dst).Error
}

func (u *MysqlClient) Get(data interface{}) error {
	err := u.client.Model(data).Where(data).First(data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("err:%v", err)
	}
	return err
}

func (u *MysqlClient) GetLast(data interface{}) error {
	err := u.client.Model(data).Where(data).Last(data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("err:%v", err)
	}
	return err
}

func (u *MysqlClient) Create(data interface{}) error {
	if err := u.client.Create(data).Error; err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// QueryListM 查询一组数据，以model为检索条件,model必须实现TableName方法
func (u *MysqlClient) QueryListM(page, num int, model, ret interface{}, order ...interface{}) error {
	m := reflect.ValueOf(model)
	f := m.MethodByName("TableName")
	fret := f.Call(nil)
	tn := fret[0].String()
	tx := u.client.Table(tn).Where(model)
	if order != nil {
		tx.Order(order[0])
	}
	return tx.Offset(page * num).Limit(num).Scan(ret).Error
}

// QueryListW 查询一组数据，以where为检索条件
func (u *MysqlClient) QueryListW(page, num int, order string, model, ret, query interface{}, args ...interface{}) (int64, error) {
	m := reflect.ValueOf(model)
	f := m.MethodByName("TableName")
	fret := f.Call(nil)
	tn := fret[0].String()
	tx := u.client.Table(tn).Where(query, args...)
	if order != "" {
		tx.Order(order)
	}
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	if err := tx.Offset(page * num).Limit(num).Scan(ret).Error; err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	return count, nil
}

// QueryList 查询一组数据
func (u *MysqlClient) QueryList(page, num int, query, order string, model, ret interface{}) (int64, error) {
	tx := u.client.Model(model).Where(query)
	if order != "" {
		tx.Order(order)
	}
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	if err := tx.Offset(page * num).Limit(num).Scan(ret).Error; err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	return count, nil
}

// QueryList 查询一组数据
func (u *MysqlClient) QueryAll(query, order string, model, ret interface{}) (int64, error) {
	tx := u.client.Model(model).Where(query)
	if order != "" {
		tx.Order(order)
	}
	// if err := tx.Count(&count).Error; err != nil {
	// 	log.Error("err:%v", err)
	// 	return 0, err
	// }
	if err := tx.Scan(ret).Error; err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	var count int64
	ref := reflect.ValueOf(ret).Elem()
	if ref.Kind() == reflect.Slice {
		count = int64(ref.Len())
	}
	return count, nil
}

// DistinctCount 去重查询个数
func (u *MysqlClient) DistinctCount(model interface{}, con, name string) (count int64) {
	if err := u.client.Model(model).Where(con).Distinct(name).Count(&count).Error; err != nil {
		log.Error("err:%v", err)
	}
	return
}

// Count 根据条件查询个数
func (u *MysqlClient) Count(model interface{}, condition string) (count int64) {
	if err := u.client.Model(model).Where(condition).Count(&count).Error; err != nil {
		log.Error("err:%v", err)
	}
	return
}

// Count 根据条件查询个数
func (u *MysqlClient) CountTable(tableName, condition string) (count int64) {
	if err := u.client.Table(tableName).Where(condition).Count(&count).Error; err != nil {
		log.Error("err:%v", err)
	}
	return
}

// Count 根据条件查询个数
func (u *MysqlClient) Exist(model interface{}) bool {
	var count int64
	if err := u.client.Model(model).Where(model).Count(&count).Error; err != nil {
		log.Error("err:%v", err)
		return false
	}
	return count > 0
}

// Sum 求和
func (u *MysqlClient) Sum(model interface{}, condition, name string) (count int64) {
	err := u.client.Model(model).Where(condition).Select(fmt.Sprintf("IFNULL(SUM(%v),0)", name)).Scan(&count).Error
	if err != nil {
		log.Error("err:%v", err)
	}
	return
}

func (u *MysqlClient) SumTable(tableName, condition, name string) (count int64) {
	err := u.client.Table(tableName).Where(condition).Select(fmt.Sprintf("IFNULL(SUM(%v),0)", name)).Scan(&count).Error
	if err != nil {
		log.Error("err:%v", err)
	}
	return
}

// sql原生语句查询
func (u *MysqlClient) QueryBySql(sqlStr string, res interface{}) error {
	err := u.client.Raw(sqlStr).Scan(res).Error
	if err != nil {
		log.Error("err:%v", err)
		return err
	} else {
		return nil
	}
}

// sql原生语句查询
func (u *MysqlClient) QueryCountBySql(sqlStr string, count interface{}) error {
	err := u.client.Raw(sqlStr).Scan(count).Error
	if err != nil {
		log.Error("err:%v", err)
		return err
	} else {
		return nil
	}
}
