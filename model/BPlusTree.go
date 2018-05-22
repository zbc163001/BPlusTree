package model

import (
	"fmt"

)

const M = 4
const INT_MAX = int(^uint(0) >> 1)
const INT_MIN = ^INT_MAX
const LIMIT_M_2  = (M  + 1)/2
type Position *BPlusFullNode


type BPlusLeafNode struct {
	Next  *BPlusFullNode
	datas []int
}
//叶子节点应该为Children为空，但leafNode中datas不为空 Next一般不为空
type BPlusFullNode struct {
	KeyNum   int
	Key      []int
	isLeaf   bool
    Children []*BPlusFullNode
    leafNode *BPlusLeafNode
}

type BPlusTree struct {
	keyMax int
	root   *BPlusFullNode
	ptr  *BPlusFullNode
}


func MallocNewNode() *BPlusFullNode {
    NewLeaf := MallocNewLeaf()
	NewNode := BPlusFullNode{
		KeyNum:   0,
		Key:      make([]int, M + 1),
		isLeaf:   true,
		Children: make([]*BPlusFullNode, M + 1),
		leafNode: NewLeaf,
	}
	for i,_ := range NewNode.Key{
		NewNode.Key[i] = INT_MIN
	}

	return &NewNode
}

func MallocNewLeaf() *BPlusLeafNode{

	NewLeaf := BPlusLeafNode{
		Next:  nil,
		datas: make([]int, M + 1),
	}
	for i,_ := range NewLeaf.datas {
		NewLeaf.datas[i] = i
	}
	return &NewLeaf
}

func(tree *BPlusTree)  Initialize() Position{

	var T Position
/* 根结点 */
	T = MallocNewNode()
	tree.ptr = T
	tree.root = T
	return T
}

func FindMostLeft(P Position) Position{
	var Tmp Position
	Tmp = P
	if Tmp.isLeaf == true || Tmp == nil{
		return Tmp
	}else if Tmp.Children[0].isLeaf == true{
		return  Tmp.Children[0]
	}else{
		for (Tmp != nil && Tmp.Children[0].isLeaf != true ) {
			Tmp = Tmp.Children[0]
		}
	}
	return  Tmp.Children[0]
}

func  FindMostRight(P Position ) Position{
	var Tmp Position
	Tmp = P

	if Tmp.isLeaf == true || Tmp == nil{
		return Tmp
	}else if Tmp.Children[Tmp.KeyNum - 1].isLeaf == true{
		return  Tmp.Children[Tmp.KeyNum - 1]
	}else{
		for (Tmp != nil && Tmp.Children[Tmp.KeyNum - 1].isLeaf != true ) {
			Tmp = Tmp.Children[Tmp.KeyNum - 1]
		}
	}

	return Tmp.Children[Tmp.KeyNum - 1]
}

/* 寻找一个兄弟节点，其存储的关键字未满，若左右都满返回nil */
func FindSibling(Parent Position,i int ) Position{
	var Sibling Position
	var upperLimit int
	upperLimit = M
	Sibling = nil
	if i == 0{
		if Parent.Children[1].KeyNum < upperLimit{

			Sibling = Parent.Children[1]
		}
	} else if (Parent.Children[i - 1].KeyNum < upperLimit){
		Sibling = Parent.Children[i - 1]
	}else if (i + 1 < Parent.KeyNum && Parent.Children[i + 1].KeyNum < upperLimit){
		Sibling = Parent.Children[i + 1]
	}
	return Sibling
}

/* 查找兄弟节点，其关键字数大于M/2 ;没有返回nil j用来标识是左兄还是右兄*/

func  FindSiblingKeyNum_M_2( Parent Position,i int, j *int) Position{
	var lowerLimit int
	var Sibling Position
	Sibling = nil

	lowerLimit = LIMIT_M_2

	if (i == 0){
		if (Parent.Children[1].KeyNum > lowerLimit){
			Sibling = Parent.Children[1]
			*j = 1
		}
	}else{
		if (Parent.Children[i - 1].KeyNum > lowerLimit){
			Sibling = Parent.Children[i - 1]
			*j = i - 1
		} else if (i + 1 < Parent.KeyNum && Parent.Children[i + 1].KeyNum > lowerLimit){
			Sibling = Parent.Children[i + 1]
			*j = i + 1
		}

	}
	return Sibling
}

/* 当要对X插入data的时候，i是X在Parent的位置，insertIndex是data要插入的位置，j可由查找得到
   当要对Parent插入X节点的时候，posAtParent是要插入的位置，Key和j的值没有用
 */
func(tree *BPlusTree) InsertElement (isData bool,Parent Position, X Position, Key int, posAtParent int , insertIndex int, data int)  Position {

	var k int
	if (isData){
/* 插入data*/
		k = X.KeyNum - 1
		for (k >= insertIndex){
			X.Key[k + 1] = X.Key[k]
			X.leafNode.datas[k + 1] = X.leafNode.datas[k]
			k--
		}

		X.Key[insertIndex] = Key
		X.leafNode.datas[insertIndex] = data

		if (Parent != nil) {
			Parent.Key[posAtParent] = X.Key[0] //可能min_key 已发生改变
		}

		X.KeyNum++

	}else{
/* 插入节点 */
/* 对树叶节点进行连接 */
		if (X.isLeaf == true){
			if (posAtParent > 0){
				Parent.Children[posAtParent- 1].leafNode.Next = X
			}
			X.leafNode.Next = Parent.Children[posAtParent]
			//更新叶子指针
			if X.Key[0] <= tree.ptr.Key[0]{
				tree.ptr = X
			}
		}

		k = Parent.KeyNum - 1
		for (k >= posAtParent){   //插入节点时key也要对应的插入
			Parent.Children[k + 1] = Parent.Children[k]
			Parent.Key[k + 1] = Parent.Key[k]
			k--
		}
		Parent.Key[posAtParent] = X.Key[0]
		Parent.Children[posAtParent] = X
		Parent.KeyNum++
	}
	return X
}
/*
两个参数X posAtParent 有些重复 posAtParent可以通过X的最小关键字查找得到
*/
func(tree *BPlusTree) RemoveElement(isData bool,Parent Position ,X Position , posAtParent int, deleteIndex int ) Position {

	var  k,keyNum int

	if (isData){
		keyNum = X.KeyNum
/* 删除key */
		k = deleteIndex + 1
		for (k < keyNum){
			X.Key[k - 1] = X.Key[k]
			X.leafNode.datas[k - 1] = X.leafNode.datas[k - 1]
			k++
		}

		X.Key[keyNum - 1] = INT_MIN
		X.leafNode.datas[keyNum - 1] = INT_MIN
		Parent.Key[posAtParent] = X.Key[0]
		X.KeyNum--
	}else{
/* 删除节点 */
/* 修改树叶节点的链接 */
		if (X.isLeaf == true && posAtParent > 0){
			Parent.Children[posAtParent - 1].leafNode.Next = Parent.Children[posAtParent + 1]
		}

		keyNum = Parent.KeyNum
		k = posAtParent + 1
		for (k < keyNum){
			Parent.Children[k - 1] = Parent.Children[k]
			Parent.Key[k - 1] = Parent.Key[k]
			k++
		}

        if X.Key[0] == tree.ptr.Key[0]{ // refresh ptr
        	tree.ptr = Parent.Children[0]
		}
		Parent.Children[Parent.KeyNum - 1] = nil
		Parent.Key[Parent.KeyNum - 1] = INT_MIN

		Parent.KeyNum--

	}
	return X
}

/* Src和Dst是两个相邻的节点，i是Src在Parent中的位置；
 将Src的元素移动到Dst中 ,n是移动元素的个数*/
func(tree *BPlusTree) MoveElement(src Position , dst Position , parent Position , posAtParent int,eNum int )  Position {
	var TmpKey,data int
	var Child Position
	var j int
	var srcInFront bool

	srcInFront = true

	if (src.Key[0] < dst.Key[0]) {
		srcInFront = false
	}
	j = 0
/* 节点Src在Dst前面 */
	if (srcInFront){
		if (src.isLeaf == false){
			for (j < eNum) {
				Child = src.Children[src.KeyNum - 1]
				tree.RemoveElement(false, src, Child, src.KeyNum - 1, INT_MIN) //每删除一个节点keyNum也自动减少1 队尾删
				tree.InsertElement(false, dst, Child, INT_MIN, 0, INT_MIN,INT_MIN) //队头加
				j++
			}
		}else{
			for (j < eNum) {
				TmpKey = src.Key[src.KeyNum -1]
				data = src.leafNode.datas[src.KeyNum - 1]
				tree.RemoveElement(true, parent, src, posAtParent, src.KeyNum - 1)
				tree.InsertElement(true, parent, dst, TmpKey, posAtParent + 1, 0,data)
				j++
			}

		}

		parent.Key[posAtParent+ 1] = dst.Key[0]
/* 将树叶节点重新连接 */
		if (src.KeyNum > 0) {
			FindMostRight(src).leafNode.Next = FindMostLeft(dst) //似乎不需要重连，src的最右本身就是dst最左的上一元素
		}else {
			if src.isLeaf == true {
				parent.Children[posAtParent - 1 ].leafNode.Next = dst
			}
			//TODO 看看在其他地方是否有类似的判断
			tree.RemoveElement(false ,parent.parent，parent ,parentIndex,INT_MIN )
		}
	}else{
		if (src.isLeaf == false){
			for (j < eNum) {
				Child = src.Children[0]
				tree.RemoveElement(false, src, Child, 0, INT_MIN)
				tree.InsertElement(false, dst, Child, INT_MIN, dst.KeyNum, INT_MIN,INT_MIN)
				j++
			}

		}else{
			for (j < eNum) {
				TmpKey = src.Key[0]
				data = src.leafNode.datas[0]
				tree.RemoveElement(true, parent, src, posAtParent, 0)
				tree.InsertElement(true, parent, dst, TmpKey, posAtParent - 1, dst.KeyNum,data)
				j++
			}

		}

		parent.Key[posAtParent] = src.Key[0]
		if (src.KeyNum > 0) {
			FindMostRight(dst).leafNode.Next = FindMostLeft(src)
		}else {
			if src.isLeaf == true {
				dst.leafNode.Next = src.leafNode.Next
			}
			tree.RemoveElement(false ,parent.parent，parent ,parentIndex,INT_MIN )
		}
	}

	return parent
}
//i为节点X的位置
func(tree *BPlusTree)  SplitNode(Parent Position, X Position, i int) Position{
	var  j,k, keyNum int
	var NewNode Position

	NewNode = MallocNewNode()

	k = 0
	j = X.KeyNum / 2
	keyNum = X.KeyNum
	for (j < keyNum){
		if (X.Children[0] != nil){  //Internal node
			NewNode.Children[k] = X.Children[j]
			X.Children[j] = nil
		}
		NewNode.Key[k] = X.Key[j]
		X.Key[j] = INT_MIN
		NewNode.KeyNum++
		X.KeyNum--
		j++
		k++
	}

	if (Parent != nil) {

		tree.InsertElement(false, Parent, NewNode, INT_MIN, i+1, INT_MIN, INT_MIN)
	}else{
/* 如果是X是根，那么创建新的根并返回 */
		Parent = MallocNewNode()
		tree.InsertElement(false , Parent, X, INT_MIN, 0, INT_MIN, INT_MIN)
		tree.InsertElement(false, Parent, NewNode,INT_MIN, 1, INT_MIN,INT_MIN)

		return Parent
	}

	return X   //为什么返回一个X一个Parent
}

/* 合并节点,X少于M/2关键字，S有大于或等于M/2个关键字*/
func(tree *	BPlusTree) MergeNode( Parent Position,  X Position, S Position, i int) Position{
	var Limit int

/* S的关键字数目大于M/2 */
	if (S.KeyNum > LIMIT_M_2){
/* 从S中移动一个元素到X中 */
		tree.MoveElement(S, X, Parent, i,1)
	}else{
/* 将X全部元素移动到S中，并把X删除 */
		Limit = X.KeyNum
		tree.MoveElement(X,S, Parent, i,Limit) //最多时S恰好MAX
		tree.RemoveElement(false, Parent, X, i, INT_MIN)
	}
	return Parent
}

//自上往下递归，但是第一个节点怎么指定？  TODO  从根开始向下遍历？ 切片的初始化记得先搞好 倒数第二层节点
func(tree *BPlusTree)  RecursiveInsert( beInsertedElement Position, Key int, i int , Parent Position,data int) Position{
	var  InsertIndex,upperLimit int
	var  Sibling Position

/* 查找分支 */
	InsertIndex = 0
	for  (InsertIndex < beInsertedElement.KeyNum && Key >= beInsertedElement.Key[InsertIndex]){
/* 重复值不插入 */
		if (Key == beInsertedElement.Key[InsertIndex]){
			return beInsertedElement
		}
		InsertIndex++
	}
	//似乎不用
	//if (InsertIndex != 0 && beInsertedElement.Children[0] != nil) {
	//	InsertIndex--
	//}

/* 树叶 */
	if (beInsertedElement.isLeaf == true) {
		beInsertedElement = tree.InsertElement(true, Parent, beInsertedElement, Key, i, InsertIndex,data) //返回叶子节点
		/* 内部节点 */
	}else {
		//
		beInsertedElement.Children[InsertIndex] = tree.RecursiveInsert(beInsertedElement.Children[InsertIndex], Key, InsertIndex, beInsertedElement,data)
	}
/* 调整节点 */

	upperLimit = M

	if (beInsertedElement.KeyNum > upperLimit){
/* 根 */
		if (Parent == nil){
/* 分裂节点 */
			beInsertedElement = tree.SplitNode(Parent, beInsertedElement, i)
		} else{
			Sibling = FindSibling(Parent, i)
			if (Sibling != nil){
/* 将T的一个元素（Key或者Child）移动的Sibing中 */
				tree.MoveElement(beInsertedElement, Sibling, Parent, i, 1)
			}else{
/* 分裂节点 */
				beInsertedElement = tree.SplitNode(Parent, beInsertedElement, i)
			}
		}

	}

	if (Parent != nil) {
		Parent.Key[i] = beInsertedElement.Key[0]
	}


return beInsertedElement
}

/* 插入 */
func(tree *BPlusTree) Insert( T Position, Key int,data int) Position{
	return tree.RecursiveInsert(T, Key, 0, nil,data)
}

func(tree *BPlusTree) RecursiveRemove( T Position, Key int, i int, Parent Position) Position{

	var  j int
	var Sibling,Tmp Position //TODO 查看tmp具体用处
	var NeedAdjust bool

	Sibling = nil

/* 查找分支 */
	j = 0
	for (j < T.KeyNum && Key >= T.Key[j]){
		if (Key == T.Key[j]) {
			break
		}
		j++
	}

	if (T.Children[0] == nil){
/* 没找到 */
		if (Key != T.Key[j] || j == T.KeyNum) {
			return T
		}
	}else {
		if (j == T.KeyNum || Key < T.Key[j]) {
			j--
		}
	}


/* 树叶 */
	if (T.Children[0] == nil){
		T = tree.RemoveElement(1, Parent, T, i, j)
	}else{
		T.Children[j] = tree.RecursiveRemove(T.Children[j], Key, j, T)
	}

	NeedAdjust = false
/* 树的根或者是一片树叶，或者其儿子数在2到M之间 */
	if (Parent == nil && T.Children[0] != nil && T.KeyNum < 2){
		NeedAdjust = true
	} else if (Parent != nil && T.Children[0] != nil && T.KeyNum < LIMIT_M_2){
		/* 除根外，所有非树叶节点的儿子数在[M/2]到M之间。(符号[]表示向上取整) */
		NeedAdjust = true
	} else if (Parent != nil && T.Children[0] == nil && T.KeyNum < LIMIT_M_2){
		/* （非根）树叶中关键字的个数也在[M/2]和M之间 */
		NeedAdjust = true
	}

/* 调整节点 */
	if (NeedAdjust){
/* 根 */
		if (Parent == nil){
			if(T.Children[0] != nil && T.KeyNum < 2){
				Tmp = T
				T = T.Children[0]

				return T
			}

		}else{
/* 查找兄弟节点，其关键字数目大于M/2 */
			Sibling = FindSiblingKeyNum_M_2(Parent, i,&j)
			if (Sibling != nil){
				tree.MoveElement(Sibling, T, Parent, j, 1)
			}else{
				if (i == 0){
					Sibling = Parent.Children[1]
				} else{
					Sibling = Parent.Children[i - 1]
				}

				Parent = tree.MergeNode(Parent, T, Sibling, i)
				T = Parent.Children[i]
			}
		}

	}


return T
}

/* 删除 */
func(tree *BPlusTree) Remove( T Position, Key int) Position{
	return tree.RecursiveRemove(T, Key, 0, nil)
}


/* 销毁 */
func Destroy(T Position) Position{
	var i,j int
	if (T != nil){
		i = 0
		for (i < T.KeyNum + 1){
			Destroy(T.Children[i])
			i++
		}

		fmt.Printf("Destroy:(")
		j = 0
		for (j < T.KeyNum){/*  T->Key[i] != Unavailable*/
			fmt.Printf("%d:",T.Key[j])
			j++
		}
		fmt.Printf(") ")
	}

	return T
}

func RecursiveTravel( T Position, Level int){
	var i int
	if (T != nil){
		fmt.Printf("  ")
		fmt.Printf("[Level:%d]-->",Level)
		fmt.Printf("(")
		i = 0
		for (i < T.KeyNum){/*  T->Key[i] != Unavailable*/
			fmt.Printf("%d:",T.Key[i])
			i++
		}
		fmt.Printf(")")

		Level++

		i = 0
		for (i <= T.KeyNum) {
			RecursiveTravel(T.Children[i], Level)
			i++
		}
	}
}

/* 遍历 */
func Travel( T Position){
	RecursiveTravel(T, 0)
	fmt.Printf("\n")
}

/* 遍历树叶节点的数据 */
func TravelData( T Position){
	var Tmp Position
	var i int
	if (T == nil){
		return
	}
	fmt.Printf("All Data:")
	Tmp = T
	for (Tmp.Children[0] != nil){
		Tmp = Tmp.Children[0]
	}
/* 第一片树叶 */
	for(Tmp != nil){
		i = 0
		for (i < Tmp.KeyNum){
			fmt.Printf(" %d",Tmp.Key[i++])
		}
		Tmp = Tmp.leafNode.Next
	}
}

//TODO
//查找算法，根据KEY找到对应的区块
func main() {
	var p,q,j   int   = 1,2,3
	var test bool
	if (p==q || q== j){
		test = true
	}
	fmt.Println(test)

}
