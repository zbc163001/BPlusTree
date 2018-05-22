package main

import ("fmt"
	//"strconv"
)
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
	"\n")
}

func (n node2) print() {

}
func main() {
	//var i int
	//i = 1
	//if (i == 1) {
	//	fmt.Printf("Hello, world.\n")
	//}else if (i == 2)  {
	//	fmt.Printf("help")
	//}else {
	//	fmt.Print("default")
	//}
	//node := node1{
	//	name1: "name2",
	//	father: father{"father"},
	//	  //也可以去掉逗号让大括号跟在"后面
	//	}
	//	node.print()
		 var s string
		 s  = "abc"
		// fmt.Printf(s)
		var bytes []byte
		bytes = []byte(s)
		var s2 string
		s2 = string(bytes)
 	//	 a,err := strconv.Atoi(s)
	//fmt.Printf(strconv.Itoa(a))
	//	 if err == nil {
	//	 	fmt.Printf(strconv.Itoa(a))
	//	 }
		 args := make(map[string]int)
		 args["hello"] = 1
		 args["world"] = 2
		 fmt.Print(s2)
		 var test *int64
		 *test = 1
		 //testAdd(&test)
		 fmt.Print(*test)
}
func testAdd(test *int64)  {
	*test = *test + 1
}