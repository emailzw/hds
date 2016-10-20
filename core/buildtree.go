package core

import (
	"hds/tree"
	"hds/utils"
	"sort"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func BuildTree(rootid string, deep int) *tree.TreeNode {
	redisCon0 := utils.GetRedisConnection(0)
	defer redisCon0.Close()
	node := tree.NewTreeNode(rootid)
	node.Deep = deep
	r, e := redis.Strings(redisCon0.Do("SMEMBERS", rootid))
	utils.CheckError(e)
	sort.Strings(r)
	for _, v := range r {
		child := BuildTree(v, deep+1)
		child.Parent = node
		node.Leafs = append(node.Leafs, child)
	}
	return node
}
