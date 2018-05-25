package main

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


func MallocNewNode(isLeaf bool) *BPlusFullNode {
	var NewNode *BPlusFullNode
	if isLeaf == true {
		NewLeaf := MallocNewLeaf()
		NewNode = &BPlusFullNode{
			KeyNum:   0,
			Key:      make([]int, M + 1),  //申请M + 1是因为插入时可能暂时出现节点key大于M 的情况,待后期再分裂处理
			isLeaf:   isLeaf,
			Children: nil,
			leafNode: NewLeaf,
		}
	}else{
		NewNode = &BPlusFullNode{
			KeyNum:  0,
			Key:     make([]int, M + 1),
			isLeaf:   isLeaf,
			Children: make([]*BPlusFullNode, M + 1),
			leafNode: nil,
		}
	}
	for i,_ := range NewNode.Key{
		NewNode.Key[i] = INT_MIN
	}

	return NewNode
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

func(tree *BPlusTree) Initialize() {

/* 根结点 */
	T := MallocNewNode(true)
	tree.ptr = T
	tree.root = T
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
			X.leafNode.datas[k - 1] = X.leafNode.datas[k]
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

/* Src和Dst是两个相邻的节点，posAtParent是Src在Parent中的位置；
 将Src的元素移动到Dst中 ,eNum是移动元素的个数*/
func(tree *BPlusTree) MoveElement(src Position , dst Position , parent Position , posAtParent int,eNum int )  Position {
	var TmpKey,data int
	var Child Position
	var j int
	var srcInFront bool

	srcInFront = false

	if (src.Key[0] < dst.Key[0]) {
		srcInFront = true
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
			//  此种情况肯定是merge merge中有实现先移动再删除操作
			//tree.RemoveElement(false ,parent.parent，parent ,parentIndex,INT_MIN )
		}
	}else{
		if (src.isLeaf == false){
			for (j < eNum) {
				Child = src.Children[0]
				tree.RemoveElement(false, src, Child, 0, INT_MIN)  //从src的队头删
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
			//tree.RemoveElement(false ,parent.parent，parent ,parentIndex,INT_MIN )
		}
	}

	return parent
}
//i为节点X的位置
func(tree *BPlusTree)  SplitNode(Parent Position, beSplitedNode Position, i int) Position{
	var  j,k, keyNum int
	var NewNode Position

	if beSplitedNode.isLeaf == true {
		NewNode = MallocNewNode(true)
	}else{
		NewNode = MallocNewNode(false)
	}

	k = 0
	j = beSplitedNode.KeyNum / 2
	keyNum = beSplitedNode.KeyNum
	for (j < keyNum){
		if (beSplitedNode.isLeaf == false){ //Internal node
			NewNode.Children[k] = beSplitedNode.Children[j]
			beSplitedNode.Children[j] = nil
		}else {
			NewNode.leafNode.datas[k] = beSplitedNode.leafNode.datas[j]
			beSplitedNode.leafNode.datas[j] = INT_MIN
		}
		NewNode.Key[k] = beSplitedNode.Key[j]
		beSplitedNode.Key[j] = INT_MIN
		NewNode.KeyNum++
		beSplitedNode.KeyNum--
		j++
		k++
	}

	if (Parent != nil) {
		tree.InsertElement(false, Parent, NewNode, INT_MIN, i+1, INT_MIN, INT_MIN)
		// parent > limit 时的递归split recurvie中实现
	}else{
/* 如果是X是根，那么创建新的根并返回 */
		Parent = MallocNewNode(false)
		tree.InsertElement(false , Parent, beSplitedNode, INT_MIN, 0, INT_MIN, INT_MIN)
		tree.InsertElement(false, Parent, NewNode,INT_MIN, 1, INT_MIN,INT_MIN)
		tree.root = Parent
		return Parent
	}

	return beSplitedNode
	// 为什么返回一个X一个Parent?
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
		tree.MoveElement(X,S, Parent, i,Limit) //最多时S恰好MAX MoveElement已考虑了parent.key的索引更新
		tree.RemoveElement(false, Parent, X, i, INT_MIN)
	}
	return Parent
}

func(tree *BPlusTree)  RecursiveInsert( beInsertedElement Position, Key int, posAtParent int , Parent Position,data int) (Position, bool){
	var  InsertIndex,upperLimit int
	var  Sibling Position
	var  result bool
    result = true
/* 查找分支 */
	InsertIndex = 0
	for  (InsertIndex < beInsertedElement.KeyNum && Key >= beInsertedElement.Key[InsertIndex]){
/* 重复值不插入 */
		if (Key == beInsertedElement.Key[InsertIndex]){
			return beInsertedElement ,false
		}
		InsertIndex++
	}
    //key必须大于被插入节点的最小元素，才能插入到此节点，故需回退一步
	if (InsertIndex != 0 && beInsertedElement.isLeaf == false) {
		InsertIndex--
	}

/* 树叶 */
	if (beInsertedElement.isLeaf == true) {
		beInsertedElement = tree.InsertElement(true, Parent, beInsertedElement, Key, posAtParent, InsertIndex,data) //返回叶子节点
		/* 内部节点 */
	}else {
		beInsertedElement.Children[InsertIndex],result = tree.RecursiveInsert(beInsertedElement.Children[InsertIndex], Key, InsertIndex, beInsertedElement,data)
		//更新parent发生在split时
	}
/* 调整节点 */

	upperLimit = M
	if (beInsertedElement.KeyNum > upperLimit){
/* 根 */
		if (Parent == nil){
/* 分裂节点 */
			beInsertedElement = tree.SplitNode(Parent, beInsertedElement, posAtParent)
		} else{
			Sibling = FindSibling(Parent, posAtParent)
			if (Sibling != nil){
/* 将T的一个元素（Key或者Child）移动的Sibing中 */
				tree.MoveElement(beInsertedElement, Sibling, Parent, posAtParent, 1)
			}else{
/* 分裂节点 */
				beInsertedElement = tree.SplitNode(Parent, beInsertedElement, posAtParent)
			}
		}

	}
	if (Parent != nil) {
		Parent.Key[posAtParent] = beInsertedElement.Key[0]
	}

return beInsertedElement, result
}

/* 插入 */
func(tree *BPlusTree) Insert(  Key int,data int) (Position,bool){
	return tree.RecursiveInsert(tree.root, Key, 0, nil, data) //从根节点开始插入
}

func(tree *BPlusTree) RecursiveRemove( beRemovedElement Position, Key int, posAtParent int, Parent Position) (Position, bool){

	var  deleteIndex int
	var Sibling Position
	var NeedAdjust bool
	var result bool
	Sibling = nil

	/* 查找分支   TODO查找函数可以在参考这里的代码 或者实现一个递归遍历*/
	deleteIndex = 0
	for (deleteIndex < beRemovedElement.KeyNum && Key >= beRemovedElement.Key[deleteIndex]){
		if (Key == beRemovedElement.Key[deleteIndex]) {
			break
		}
		deleteIndex++
	}

	if (beRemovedElement.isLeaf == true){
/* 没找到 */
		if (Key != beRemovedElement.Key[deleteIndex] || deleteIndex == beRemovedElement.KeyNum) {
			return beRemovedElement, false
		}
	}else {
		if (deleteIndex == beRemovedElement.KeyNum || Key < beRemovedElement.Key[deleteIndex]) {
			deleteIndex-- //准备到下层节点查找
		}
	}

/* 树叶 */
	if (beRemovedElement.isLeaf == true){
		beRemovedElement = tree.RemoveElement(true, Parent, beRemovedElement, posAtParent, deleteIndex)
	}else{
		beRemovedElement.Children[deleteIndex],result = tree.RecursiveRemove(beRemovedElement.Children[deleteIndex], Key, deleteIndex, beRemovedElement)
	}

	NeedAdjust = false
	//有子节点的root节点，当keyNum小于2时
	if (Parent == nil && beRemovedElement.isLeaf == false && beRemovedElement.KeyNum < 2){
		NeedAdjust = true
	} else if (Parent != nil && beRemovedElement.isLeaf == false && beRemovedElement.KeyNum < LIMIT_M_2){
		/* 除根外，所有中间节点的儿子数不在[M/2]到M之间时。(符号[]表示向上取整) */
		NeedAdjust = true
	} else if (Parent != nil && beRemovedElement.isLeaf == true && beRemovedElement.KeyNum < LIMIT_M_2){
		/* （非根）树叶中关键字的个数不在[M/2]到M之间时 */
		NeedAdjust = true
	}

/* 调整节点 */
	if (NeedAdjust){
/* 根 */
		if (Parent == nil){
			if(beRemovedElement.isLeaf == false && beRemovedElement.KeyNum < 2){
				//树根的更新操作 树高度减一
				beRemovedElement = beRemovedElement.Children[0]
				tree.root = beRemovedElement.Children[0]
				return beRemovedElement,true
			}

		}else{
/* 查找兄弟节点，其关键字数目大于M/2 */
			Sibling = FindSiblingKeyNum_M_2(Parent, posAtParent,&deleteIndex)
			if (Sibling != nil){
				tree.MoveElement(Sibling, beRemovedElement, Parent, deleteIndex, 1)
			}else{
				if (posAtParent == 0){
					Sibling = Parent.Children[1]
				} else{
					Sibling = Parent.Children[posAtParent- 1]
				}

				Parent = tree.MergeNode(Parent, beRemovedElement, Sibling, posAtParent)
				//Merge中已考虑空节点的删除
				beRemovedElement = Parent.Children[posAtParent]
			}
		}

	}

	return beRemovedElement ,result
}

/* 删除 */
func(tree *BPlusTree) Remove(Key int) (Position,bool){
	return tree.RecursiveRemove(tree.root, Key, 0, nil)
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

func(tree *BPlusTree) FindData(key int) (int, bool) {
	var currentNode *BPlusFullNode
	var index int
	currentNode = tree.root
	for index < currentNode.KeyNum {
		index = 0
		for key >= currentNode.Key[index] && index < currentNode.KeyNum{
			index ++
		}
		if  index == 0 {
			return  INT_MIN, false
		}else{
			index--
			if currentNode.isLeaf == false {
				currentNode = currentNode.Children[index]
			}else{
				if key == currentNode.Key[index] {
					return currentNode.leafNode.datas[index], true
				}else{
					return  INT_MIN, false
				}
			}
		}

	}
	return  INT_MIN, false
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
			fmt.Printf(" %d",Tmp.Key[i])
			i++
		}
		Tmp = Tmp.leafNode.Next
	}
}

func main (){
	var tree BPlusTree
	(&tree).Initialize()
	var i int
	i = 1
	for i< 9 {
		_ ,result:= tree.Insert(i,i * 10)
		fmt.Print(i)
		if result == false {
			print("数据已存在")
		}
		i++
	}

	tree.Remove(7)
	tree.Remove(6)
	tree.Remove(5)
	resultDate,success:=tree.FindData(5)
	if success == true {
		fmt.Print(resultDate)
		fmt.Printf("\n")
	}

	//遍历结点元素
	i = 0
	for i < tree.root.Children[1].KeyNum{
		fmt.Println(tree.root.Children[1].leafNode.datas[i])
		i++
	}
}


