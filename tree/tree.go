package tree

type TreeNode struct {
	Value  string
	Leafs  []*TreeNode
	Parent *TreeNode
	Deep   int
}

//func NewTreeNode(cust *model.Customer) *TreeNode {
//node := new(TreeNode)
//node.Value = cust
//return node
//}

func NewTreeNode(custid string) *TreeNode {
	node := new(TreeNode)
	node.Value = custid
	return node
}

type treeAction func(tree *TreeNode)

func (node *TreeNode) TravelTree(walk treeAction) {
	if node == nil {
		return
	}
	walk(node)
	for _, child := range node.Leafs {
		if child == nil {
			return
		} else {
			child.TravelTree(walk)
		}
	}
}
