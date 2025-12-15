package chi


/*
一个用go实现的跟http无关的Go Radix tree
*/
type RTree struct{
    root *rnode
}
type rnode struct{
    prefix string 
	children []*rnode
	leaf bool
	value any
}
// New 创建一个空树
func NewRtree()*RTree{
	return &RTree{}
}
// Insert 插入 key 对应的value 
// 如果key已存在，会覆盖旧值
func (t *RTree) Insert(key string, value any) {
	if key == ""{
		return 
	}
   insert(t.root,key,value)
}
// Get 查找 key 对应的 value
func (t *RTree) Get(key string)(any,bool){
	if key==""{
		return nil,false
	}
	return get(t.root,key)
}

func insert(n *rnode,key string,value any){
    search:=key 

	// 如果当前节点没有孩子，直接挂一个新节点
	if len(n.children)==0{
		child:=&rnode{
            prefix: search,
			leaf: true,
			value: value,
		}
		n.children=append(n.children,child)
		return 
	}

	// 尝试在现有children 中找到有公共前缀的
	for i,child:=range n.children{
        l:=longestCommonPrefix(search,child.prefix)
		if l==0{
			// child 不匹配
			continue
		}
		// 情况 1: child.prefix 完全是 search 的前缀
		// search: "cart"
		// child.prefix: "car"
		// -> 继续往child 里插入 "t"
		if l==len(child.prefix){
			rest:=search[l:]
			if rest==""{
               // key 恰好在这个节点结束
			   child.leaf=true
			   child.value=value
			   return 
			}
			insert(child,rest,value)
			return 
		}
		// 情况2：search 和 child.prefix 只有部分重叠
		// search: "carb"
		// child.prefix: "cart"
		// l=3，公共前缀 "car"
		// 需要把 child 拆分成中间节点 mid 和原 child:
		// mid.prefix = "car"
		// mid.children = [原 child（prefix="t"...）]
		// 然后再在 mid 下边挂一条新的边表示剩余的 search。
		common:=search[:l]

		mid:=&rnode{
			prefix: common,
			// children: []*rnode{},
		}

		// 原来的child 剩余部分挂到mid下
		child.prefix=child.prefix[l:]
		mid.children=append(mid.children,child)

		// 用mid替换原来的 child
		n.children[i]=mid

		rest:=search[l:]
		if rest==""{
			// 新 key 到mid 就结束了
			mid.leaf=true
			mid.value=value
			return 
		}
		// 新的key还有剩余，在mid下面再挂一个新的child
		newChild:=&rnode{
			prefix: rest,
			leaf: true,
			value: value,
		}
        mid.children=append(mid.children, newChild)
		return
	}
	// 情况3：没有任何child 有公共前缀，直接新建一个
	child:=&rnode{
		prefix: search,
		leaf: true,
		value: value,
	}
	n.children=append(n.children, child)
}

func get(n *rnode,key string)(any,bool){
	search:=key
	// 在children 中找一个 prefix 能匹配 search 开头的
	for _,child:=range n.children{
		if len(search)<len(child.prefix){
			// search 比 prefix 还短，不可能完全匹配
			continue
		}
		if search[:len(child.prefix)]!=child.prefix{
			continue
		}
		rest:=search[len(child.prefix):]
		if rest==""{
			// key正好在这个节点结束
			if child.leaf{
				return child.value,true
			}
			return nil,false
		}
		// key 还有剩余，递归往下周吧
		return get(child,rest)
	}
	return nil,false
}
// longestCommonPrefix 返回a和b的最长公共前缀长度
func longestCommonPrefix(a,b string)int{
	n:=len(a)
	if len(b)<n{
		n=len(b)
	}
	i:=0
	for i<n&&a[i]==b[i]{
		i++
	}
	return i 
}