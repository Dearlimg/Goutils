package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

type Snowflake struct {
	node *snowflake.Node
}

// Init 初始化
func Init(startTime time.Time, machineID int64) (*Snowflake, error) {
	//初始化开始时间
	snowflake.Epoch = startTime.UnixNano() / 1000000
	//指定机器的 ID
	node, err := snowflake.NewNode(machineID)
	if err != nil {
		return nil, err
	}
	return &Snowflake{node: node}, nil
}

// GetID 生成ID
func (sn *Snowflake) GetID() int64 {
	return sn.node.Generate().Int64() //生成并返回唯一的 snowflake ID，并转化为 Int64 类型的值（也可以是其他类型）
}
