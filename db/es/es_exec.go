package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"server/common"
	"time"

	"github.com/liangdas/mqant/log"
	"github.com/olivere/elastic/v7"
)

// 向Es数据库中新增数据 index 数据表名 params 一定要是json数据结构
func (ES *EsClient) InsertToES(index string, params interface{}) error {
	_, err := ES.client.Index().Index(index).BodyJson(params).Do(context.Background())
	if err != nil {
		log.Error("index:%v,write:%v,err:%v", index, params, err)
		return err
	}
	return nil
}

// 向Es数据库中新增数据 指数退避写入
// func (ES *EsClient) InsertToESGO(index string, params interface{}) {
// 	// util.IndexTry(func() error {
// 	// 	_, err := ES.client.Index().Index(index).BodyJson(params).Do(context.Background())
// 	// 	if err != nil {
// 	// 		log.Error("index:%v,write:%v,err:%v", index, params, err)
// 	// 		return err
// 	// 	}
// 	// 	return nil
// 	// })
// 	util.Go(func() {
// 		AddBulk(&WriteData{Index: index, Data: params})
// 	})
// }

// 向Es数据库中新增数据 index 数据表名 params 一定要是json数据结构
func (ES *EsClient) InsertToESByID(index, id string, params interface{}) error {
	_, err := ES.client.Index().Index(index).Id(id).BodyJson(params).Do(context.Background())
	if err != nil {
		log.Error("index:%v,write:%v,err:%v", index, params, err)
		return err
	}
	return nil
}

// 向Es数据库中新增数据 指数退避写入
// func (ES *EsClient) InsertToESByIDGO(index, id string, params interface{}) {
// 	// util.IndexTry(func() error {
// 	// 	_, err := ES.client.Index().Index(index).Id(id).BodyJson(params).Do(context.Background())
// 	// 	if err != nil {
// 	// 		log.Error("index:%v,write:%v,err:%v", index, params, err)
// 	// 		return err
// 	// 	}
// 	// 	return nil
// 	// })
// 	util.Go(func() {
// 		AddBulk(&WriteData{ID: id, Index: index, Data: params})
// 	})
// }

func (ES *EsClient) Count(index string, q elastic.Query) int64 {
	count, err := ES.client.Count(index).Query(q).Do(context.Background())
	if err != nil {
		return 0
	}
	return count
}

// CountCard 统计数量,按提供字段去重
func (ES *EsClient) CountCard(index, field string, q elastic.Query) int64 {
	result, err := ES.C().Search().Index(index).Query(q).Aggregation("Total", elastic.NewCardinalityAggregation().Field(field)).Size(0).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return 0
	}
	ret := map[string]int64{}
	err = json.Unmarshal(result.Aggregations["Total"], &ret)
	if err != nil {
		log.Error("err:%v", err)
		return 0
	}
	return ret["value"]
}

// Update 更新文档，val必须是map|||struct||es script
func (ES *EsClient) Update(index, id string, val interface{}) (res *elastic.UpdateResponse, err error) {
	switch t := val.(type) {
	case *elastic.Script:
		res, err = ES.client.Update().Index(index).Id(id).Script(t).Do(context.Background())
		return
	default: // struct or map
		res, err = ES.client.Update().Index(index).Id(id).Doc(val).Do(context.Background())
		return
	}
}

func (ES *EsClient) DeleteByID(index string, id string) error {
	_, err := ES.client.Delete().Index(index).Id(id).Do(context.Background())
	return err
}

func (ES EsClient) DeleteByQuery(indices string, q elastic.Query) (err error) {
	_, err = ES.client.DeleteByQuery(indices).Query(q).WaitForCompletion(false).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
	}
	return
}

func (ES *EsClient) DeleteIndex(index ...string) (*elastic.IndicesDeleteResponse, error) {
	return ES.client.DeleteIndex(index...).Do(context.Background())
}

func (ES *EsClient) C() *elastic.Client {
	return ES.client
}

// 查询一组数据 index表名，page页码，num一页数量，q查询语句，kind取值的类型(必须是指针,其中包含的ID会有特殊含义，与文档中的_id同步)，sort排序(false为逆序)
func (ES *EsClient) QueryList(index string, page, num int, q elastic.Query, kind interface{}, sort ...interface{}) (int64, error) {
	from := page * num
	if from < 0 || num == 0 {
		return 0, errors.New("invalid page or num")
	}
	oneType := reflect.TypeOf(kind).Kind()
	if oneType != reflect.Ptr {
		return 0, errors.New("invalid kind")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	searchRes := new(elastic.SearchResult)
	var err error
	if len(sort) > 0 {
		next := ES.client.Search(index).Query(q)
		for i := 0; i < len(sort)-1; i += 2 {
			field, ok := sort[i].(string)
			if !ok {
				log.Error("invalid sort:%v", sort...)
				return 0, errors.New("invalid sort")
			}
			seq, ok := sort[i+1].(bool)
			if !ok {
				log.Error("invalid sort:%v", sort...)
				return 0, errors.New("invalid sort")
			}
			next.Sort(field, seq)
		}
		searchRes, err = next.From(from).Size(num).Do(ctx)
	} else {
		searchRes, err = ES.client.Search(index).Query(q).From(from).Size(num).Do(ctx)
	}
	// ret := []interface{}{}
	if elastic.IsNotFound(err) {
		return 0, nil
	}
	if err != nil {
		log.Error("search error:%v", err)
		return 0, err
	}
	newArr := make([]reflect.Value, 0)
	val := reflect.ValueOf(kind)
	for _, v := range searchRes.Hits.Hits {
		t := reflect.TypeOf(kind)
		tee := reflect.New(t.Elem().Elem())
		ti := tee.Interface()
		err := json.Unmarshal(v.Source, ti)
		if err != nil {
			return 0, err
		}
		e := reflect.ValueOf(ti).Elem()
		if e.Kind() == reflect.Struct {
			if id := e.FieldByName("ID"); id.Kind() == reflect.String {
				id.SetString(v.Id)
			}
		} else if e.Kind() == reflect.Ptr {
			if e.Elem().Kind() == reflect.Struct {
				if id := e.Elem().FieldByName("ID"); id.Kind() == reflect.String {
					id.SetString(v.Id)
				}
			}
		}
		newArr = append(newArr, reflect.ValueOf(ti).Elem())
	}
	val.Elem().Set(reflect.Append(val.Elem(), newArr...))
	count := ES.Count(index, q)
	return count, nil
}

// 查询一组数据 index表名,after上次查询最后一个数据，num一页数量，q查询语句，kind取值的类型(必须是指针,其中包含的ID会有特殊含义，与文档中的_id同步)，sort排序(false为逆序)
func (ES *EsClient) SearchAfter(index string, after interface{}, num int, q elastic.Query, kind interface{}, seqSort bool) (int64, error) {
	oneType := reflect.TypeOf(kind).Kind()
	if oneType != reflect.Ptr {
		return 0, errors.New("invalid kind")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	searchRes := new(elastic.SearchResult)
	var err error
	next := ES.client.Search(index).Query(q).Sort("_seq_no", seqSort)
	// for i := 0; i < len(sort)-1; i += 2 {
	// 	field, ok := sort[i].(string)
	// 	if !ok {
	// 		log.Error("invalid sort:%v", sort...)
	// 		return 0, errors.New("invalid sort")
	// 	}
	// 	seq, ok := sort[i+1].(bool)
	// 	if !ok {
	// 		log.Error("invalid sort:%v", sort...)
	// 		return 0, errors.New("invalid sort")
	// 	}
	// 	next.Sort(field, seq)
	// }
	if after != nil {
		searchRes, err = next.Size(num).SearchAfter(after).SeqNoPrimaryTerm(true).Do(ctx)
	} else {
		searchRes, err = next.Size(num).SeqNoPrimaryTerm(true).Do(ctx)
	}
	// ret := []interface{}{}
	if elastic.IsNotFound(err) {
		return 0, nil
	}
	if err != nil {
		log.Error("search error:%v", err)
		return 0, err
	}
	newArr := make([]reflect.Value, 0)
	val := reflect.ValueOf(kind)
	var resAfter int64
	for i, v := range searchRes.Hits.Hits {
		t := reflect.TypeOf(kind)
		tee := reflect.New(t.Elem().Elem())
		ti := tee.Interface()
		err := json.Unmarshal(v.Source, ti)
		if err != nil {
			return 0, err
		}
		e := reflect.ValueOf(ti).Elem()
		if e.Kind() == reflect.Struct {
			if id := e.FieldByName("ID"); id.Kind() == reflect.String {
				id.SetString(v.Id)
			}
		} else if e.Kind() == reflect.Ptr {
			if e.Elem().Kind() == reflect.Struct {
				if id := e.Elem().FieldByName("ID"); id.Kind() == reflect.String {
					id.SetString(v.Id)
				}
			}
		}
		fmt.Println(*v.SeqNo)
		newArr = append(newArr, reflect.ValueOf(ti).Elem())
		if i == len(searchRes.Hits.Hits)-1 {
			resAfter = *v.SeqNo
		}
	}
	val.Elem().Set(reflect.Append(val.Elem(), newArr...))
	return resAfter, nil
}

func (ES *EsClient) UpdateByScript(index string, q elastic.Query, script string) (updated int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var res *elastic.BulkIndexByScrollResponse
	res, err = ES.client.UpdateByQuery(index).Query(q).Script(elastic.NewScriptInline(script)).Do(ctx)
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	updated = res.Updated
	return
}

// QueryOne 查询单条数据,kind取值的类型(必须是指针,其中包含的ID会有特殊含义，与文档中的_id同步)，sort排序(false为逆序)
func (ES *EsClient) QueryOne(index string, q elastic.Query, kind interface{}, sort ...interface{}) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	k := reflect.ValueOf(kind)
	oneType := k.Kind()
	if oneType != reflect.Ptr {
		return errors.New("invalid kind")
	}
	s := ES.client.Search(index).Query(q)
	ret := new(elastic.SearchResult)
	var err error
	if sort != nil {
		ret, err = s.Sort(sort[0].(string), sort[1].(bool)).Size(1).Do(context.Background())
	} else {
		ret, err = s.Size(1).Do(context.Background())
	}
	if err != nil && !elastic.IsNotFound(err) {
		// if !elastic.IsNotFound(err) {
		log.Error("err:%v", err)
		// }
		return err
	}
	if ret == nil || len(ret.Hits.Hits) == 0 {
		return nil
	}
	// randomIdx := rand.Intn(len(ret.Hits.Hits))
	// hit := ret.Hits.Hits[randomIdx]
	for _, v := range ret.Hits.Hits {
		// log.Debug("source:%v", v)
		err = json.Unmarshal(v.Source, kind)
		if err != nil {
			log.Error("err:%v,source:%v", err, string(v.Source))
			return err
		}
		if k.Elem().Kind() == reflect.Struct {
			// log.Debug("set:%v", v)
			if id := k.Elem().FieldByName("ID"); id.Kind() == reflect.String {
				id.SetString(v.Id)
			}
		}
		break
	}
	return nil
}

// QueryOne 查询单条数据,kind取值的类型(必须是指针,其中包含的ID会有特殊含义，与文档中的_id同步)，sort排序(false为逆序)
func (ES *EsClient) QueryOneRandom(index string, q elastic.Query, kind interface{}, sort ...interface{}) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	k := reflect.ValueOf(kind)
	oneType := k.Kind()
	if oneType != reflect.Ptr {
		return errors.New("invalid kind")
	}
	s := ES.client.Search(index).Query(q)
	ret := new(elastic.SearchResult)
	var err error
	if sort != nil {
		ret, err = s.Sort(sort[0].(string), sort[1].(bool)).Do(context.Background())
	} else {
		ret, err = s.Do(context.Background())
	}
	if err != nil {
		if !elastic.IsNotFound(err) {
			log.Error("err:%v", err)
		}
		return err
	}
	if len(ret.Hits.Hits) == 0 {
		return nil
	}
	randomIdx := rand.Intn(len(ret.Hits.Hits))
	hit := ret.Hits.Hits[randomIdx]
	err = json.Unmarshal(hit.Source, kind)
	if err != nil {
		log.Error("source:%v", string(hit.Source))
		return err
	}
	if k.Elem().Kind() == reflect.Struct {
		// log.Debug("set:%v", v)
		if id := k.Elem().FieldByName("ID"); id.Kind() == reflect.String {
			id.SetString(hit.Id)
		}
	}
	return nil
}

// Upsert 没找到插入，找到则执行更新
func (ES *EsClient) Upsert(index, id string, val interface{}, uVal ...interface{}) (ret *elastic.UpdateResponse, err error) {
	ctx := context.Background()
	var s *elastic.UpdateService
	switch t := val.(type) {
	case *elastic.Script:
		s = ES.client.Update().Index(index).Id(id).Script(t)
	default: // struct or map
		s = ES.client.Update().Index(index).Id(id).Doc(val)
	}
	if uVal == nil {
		ret, err = s.Do(ctx)
		return
	}
	// 是否强制执行更新脚本
	if len(uVal) > 1 {
		ret, err = s.Upsert(uVal[0]).ScriptedUpsert(uVal[1].(bool)).Do(ctx)
	} else {
		ret, err = s.Upsert(uVal[0]).Do(ctx)
	}
	return
}

// IncrBy 根据字段叠加
func (ES *EsClient) IncrBy(index, id, field string, val interface{}, uVal interface{}) (ret *elastic.UpdateResponse, err error) {
	ctx := context.Background()
	s := ES.client.Update().Index(index).Id(id).Script(elastic.NewScript(fmt.Sprintf("if(ctx._source.%v==null){ctx._source.%v=%v}else{ctx._source.%v+=%v}", field, field, val, field, val)))
	if uVal == nil {
		ret, err = s.Do(ctx)
		return
	}
	ret, err = s.Upsert(uVal).ScriptedUpsert(true).Do(ctx)
	return
}

func (ES *EsClient) Refresh(index string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := ES.client.Refresh(index).Do(ctx)
	if err != nil && !elastic.IsNotFound(err) {
		log.Error("refresh err:%v", err)
	}
}

// GroupBy 聚合查询数量
// size 不设置默认为10
func (ES *EsClient) GroupBy(index, field string, q elastic.Query, size int) (*common.GroupBuckets, error) {
	query := ES.C().Search().Index(index).Query(q)
	if size > 0 {
		query.Aggregation("total", elastic.NewTermsAggregation().Field(field).Size(size))
	} else {
		query.Aggregation("total", elastic.NewTermsAggregation().Field(field))
	}
	result, err := query.Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	res := new(common.GroupBuckets)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return res, nil
}

// GroupBySub 聚合查询数量
// size 不设置默认为10
func (ES *EsClient) GroupBySub(index, field, sub string, q elastic.Query, size int) (*common.Group2CardBuckets, error) {
	query := ES.C().Search().Index(index).Query(q).Aggregation("total", elastic.NewTermsAggregation().Field(field).SubAggregation("sub1", elastic.NewTermsAggregation().Field(sub)))
	if size > 0 {
		query.Size(size)
	}
	result, err := query.Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	res := new(common.Group2CardBuckets)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return res, nil
}

func (ES *EsClient) GroupBySubs(index string, q elastic.Query, size int, ret interface{}, field ...string) error {
	agg := elastic.NewTermsAggregation().Field(field[0])
	if len(field) > 1 {
		for i := 1; i < len(field); i++ {
			agg.SubAggregation(field[i], elastic.NewTermsAggregation().Field(field[i]))
		}
	}
	query := ES.C().Search().Index(index).Query(q).Aggregation("total", agg)
	if size > 0 {
		query.Size(size)
	}
	result, err := query.Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	for _, v := range result.Aggregations {
		fmt.Println(string(v))
	}
	if err := json.Unmarshal(result.Aggregations["total"], &ret); err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// GroupByCard 聚合查询数量,包含子查询card
// size 不设置默认为10
func (ES *EsClient) GroupByCard(index, field, sub string, q elastic.Query, size int) (*common.GroupCardBuckets, error) {
	query := ES.C().Search().Index(index).Query(q)
	if size > 0 {
		query.Aggregation("total", elastic.NewTermsAggregation().Field(field).SubAggregation("sub", elastic.NewCardinalityAggregation().Field(sub)).Size(size))
	} else {
		query.Aggregation("total", elastic.NewTermsAggregation().Field(field).SubAggregation("sub", elastic.NewCardinalityAggregation().Field(sub)))
	}
	result, err := query.Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	res := new(common.GroupCardBuckets)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return res, nil
}

// GroupBy2Card 聚合查询数量,包含双重子查询card
// size 不设置默认为10
func (ES *EsClient) GroupBy2Card(index, field, sub1, sub2 string, q elastic.Query, size int) (*common.Group2CardBuckets, error) {
	query := ES.C().Search().Index(index).Query(q).Aggregation("total", elastic.NewTermsAggregation().Field(field).SubAggregation("sub1", elastic.NewTermsAggregation().Field(sub1).SubAggregation("sub2", elastic.NewCardinalityAggregation().Field(sub2))))
	if size > 0 {
		query.Size(size)
	}
	result, err := query.Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	res := new(common.Group2CardBuckets)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return res, nil
}

// GroupBy 聚合查询数量 field聚合的字段 sumField求和的字段
func (ES *EsClient) GroupSumBy(index, field string, q elastic.Query, ret interface{}, order string, or bool, size int, sumField ...string) error {
	agg := elastic.NewTermsAggregation().Field(field)
	for _, v := range sumField {
		agg.SubAggregation(v, elastic.NewSumAggregation().Field(v))
	}
	if order != "" {
		agg.OrderByAggregation(order, or)
	}
	agg.SubAggregation("Top", elastic.NewTopHitsAggregation().From(0).Size(1))
	if size == 0 {
		size = 5000
	}
	agg.Size(size)
	result, err := ES.C().Search().Index(index).Query(q).Aggregation("total", agg).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	if err := json.Unmarshal(result.Aggregations["total"], &ret); err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// Group2SumBy 聚合查询数量 field聚合的字段 sub2求和的字段
func (ES *EsClient) Group2SumBy(index, field, sub1, sub2 string, q elastic.Query, order string, or bool, size int) (*common.Group2SumBuckets, error) {
	agg := elastic.NewTermsAggregation().Field(field).SubAggregation("sub1", elastic.NewTermsAggregation().Field(sub1).SubAggregation("sub2", elastic.NewSumAggregation().Field(sub2)))
	if order != "" {
		agg.OrderByAggregation(order, or)
	}
	// agg.SubAggregation("Top", elastic.NewTopHitsAggregation().From(0).Size(1))
	if size == 0 {
		size = 5000
	}
	agg.Size(size)
	ret := new(common.Group2SumBuckets)
	result, err := ES.C().Search().Index(index).Query(q).Aggregation("total", agg).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := json.Unmarshal(result.Aggregations["total"], &ret); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return ret, nil
}

// SumBy 聚合求和数量
func (ES *EsClient) SumBy(index, field string, q elastic.Query) (float64, error) {
	result, err := ES.C().Search().Index(index).Query(q).Aggregation("total", elastic.NewSumAggregation().Field(field)).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	// fmt.Println(string(result.Aggregations["total"]))
	res := new(common.SumByResult)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return res.Value, err
	}
	return res.Value, nil
}

// SumByInt64 聚合求和数量
func (ES *EsClient) SumByInt64(index, field string, q elastic.Query) int64 {
	result, err := ES.C().Search().Index(index).Query(q).Aggregation("total", elastic.NewSumAggregation().Field(field)).Size(0).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return 0
	}
	// fmt.Println(string(result.Aggregations["total"]))
	res := new(common.SumByResult)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return int64(res.Value)
	}
	return int64(res.Value)
}

// AVG 聚合求平均值
func (ES *EsClient) AVG(index, field string, q elastic.Query) (float64, error) {
	result, err := ES.C().Search().Index(index).Query(q).Aggregation("total", elastic.NewAvgAggregation().Field(field)).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return 0, err
	}
	// fmt.Println(string(result.Aggregations["total"]))
	res := new(common.SumByResult)
	if err := json.Unmarshal(result.Aggregations["total"], &res); err != nil {
		log.Error("err:%v", err)
		return res.Value, err
	}
	return res.Value, nil
}

// Exist 判断是否存在
func (ES *EsClient) Exist(index, id string) bool {
	ret, err := ES.C().Exists().Index(index).Id(id).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return false
	}
	return ret
}

// Exist 判断是否存在
func (ES *EsClient) IndexExist(index string) bool {
	ret, err := ES.C().IndexExists(index).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return false
	}
	return ret
}

// Exist 判断是否存在
func (ES *EsClient) CreateExist(index string) error {
	_, err := ES.C().CreateIndex(index).Do(context.Background())
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// Exist 判断是否存在
func (ES *EsClient) AddAlias(aliasName, indexName string) error {
	alias := elastic.NewAliasAddAction(aliasName).Index(indexName).IsWriteIndex(true)
	aliasResult, err := ES.client.Alias().Action(alias).Do(context.Background())
	if err != nil || !aliasResult.Acknowledged {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

func (ES *EsClient) Rollover(aliasName, indexName string, conditon map[string]interface{}) error {
	req := ES.client.RolloverIndex(aliasName).NewIndex(indexName)
	if conditon != nil {
		for k, v := range conditon {
			req.AddCondition(k, v)
		}
	}
	aliasResult, err := req.Do(context.Background())
	if err != nil || !aliasResult.Acknowledged {
		log.Error("err:%v", err)
		return err
	}
	return nil
}
