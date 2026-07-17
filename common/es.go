package common

const (
	ESIndexBalance = "balance" // 流水记录别名
)

// GroupBuckets group聚合查询对象
type GroupBuckets struct {
	Buckets []struct {
		Key       interface{}
		Doc_count int
	}
}

// GroupCardBuckets group聚合查询对象
type GroupCardBuckets struct {
	Buckets []struct {
		Key       interface{}
		Doc_count int
		Sub       struct {
			Value int
		}
	}
}

// Group2CardBuckets group聚合查询对象
type Group2CardBuckets struct {
	Buckets []struct {
		Key       interface{}
		Doc_count int
		Sub1      struct {
			Buckets []struct {
				Key       interface{}
				Doc_count int
				Sub2      struct {
					Value int
				}
			}
		}
	}
}

// GroupSumBuckets group聚合查询对象
type GroupSumBuckets struct {
	Buckets []struct {
		Key       interface{}
		Doc_count int
		Value     struct {
			Value float64
		}
	}
}

// Group2SumBuckets group聚合查询对象
type Group2SumBuckets struct {
	Buckets []struct {
		Key       interface{}
		Doc_count int
		Sub1      struct {
			Buckets []struct {
				Key       interface{}
				Doc_count int
				Sub2      struct {
					Value float64
				}
			}
		}
	}
}

// SumByResult sum聚合查询对象
type SumByResult struct {
	Value float64
}

// ESPanic 崩溃记录
type ESPanic struct {
	Module string
	Func   string
	Error  string
	Time   int64
}
