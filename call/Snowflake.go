package call

import (
	"github.com/bwmarrin/snowflake"
)

var (
	snowNode *snowflake.Node
)

func NewSnowflake(id int64) {
	id = id % 1023
	var err error
	snowNode, err = snowflake.NewNode(id)
	if err != nil {
		panic(err)
	}
}

func SnowNode() *snowflake.Node {
	return snowNode
}
