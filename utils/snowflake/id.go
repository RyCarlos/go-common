// Package snowflake -----------------------------
// @file      : id.go
// @author    : Carlos
// @contact   : 534994749@qq.com
// @time      : 2025/6/18 16:20
// -------------------------------------------
package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"sync"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	node *snowflake.Node
	once sync.Once
)

// InitNode Init 初始化雪花Id节点
func InitNode(machineId int64) error {
	var err error
	once.Do(func() {
		node, err = snowflake.NewNode(machineId)
	})
	return err
}

func BuildId() int64 {
	if node == nil {
		panic("snowflake node not initialized")
	}
	return node.Generate().Int64()
}
