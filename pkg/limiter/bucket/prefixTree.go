package bucket

/*
前缀树
*/

type PrefixTree struct {
	suffix map[string]*PrefixTree //存储子节点的映射表
	result interface{}            //存储结果数据
}

// NewPrefixTree 创建新的 PrefixTree 实例
func NewPrefixTree() *PrefixTree {
	return &PrefixTree{suffix: make(map[string]*PrefixTree)}
}

// Put 在 PrefixTree 中插入数据
func (t *PrefixTree) Put(prefix []string, v interface{}) {
	root := t
	for _, s := range prefix {
		if root.suffix[s] == nil { // 如果当前节点的子节点中没有 s 对应的节点
			root.suffix[s] = NewPrefixTree() // 在当前节点的子节点中创建一个新的 PrefixTree 节点
		}
		root = root.suffix[s] // 将当前节点指向 s 对应的子节点
	}
	root.result = v // 在最终节点存储结果数据 v
}

func (t *PrefixTree) Get(prefix []string) interface{} {
	root := t
	for _, s := range prefix {
		if root.suffix[s] != nil { // 如果当前节点的子节点中存在 s 对应的节点
			root = root.suffix[s] // 将当前节点指向 s 对应的子节点
		} else {
			break // 如果节点不存在，跳出循环
		}
	}
	return root.result
}
