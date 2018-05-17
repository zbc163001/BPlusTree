package main

import ("fmt")
type father struct{
	name string
}

type node1 struct{
	father
	name1 string
}

type node2 struct {
	father
	name2 string
}

func (n father) print1() {
	fmt.Printf("father" +
		n.name)
}
func (n node1) print() {
fmt.Printf("node1" +
	"")
}

func (n node2) print() {

}
func main() {
	var i int
	i = 1
	if (i == 1) {
		fmt.Printf("Hello, world.\n")
	}else if (i == 2)  {
		fmt.Printf("help")
	}else {
		fmt.Print("default")
	}
	node := node1{
		name1: "name2",
		father: father{"father"},
		  //也可以去掉逗号让大括号跟在"后面
		}
		node.print()
}