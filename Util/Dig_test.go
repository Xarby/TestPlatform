package Util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDegDomain(t *testing.T) {

	//split1 := []string{"5.6.2.1","9.9.9.9"}
	//sort.Strings(split1)
	//fmt.Println(split1)
	a:= DigDoamin("10.1.104.55","www.baidu.com","A")
	b:= DigDoamin("10.1.121.116","www.baidu.com","A")
	//fmt.Println(a)
	//sort.Strings(b)
	//fmt.Println(len(b))
	//fmt.Println(b)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println( reflect.DeepEqual(a, b) )
	//fmt.Println(len(a))
	//for _, v := range a {
	//	fmt.Println(v)
	//}
	//fmt.Println(a == []string{"5.6.2.1","9.9.9.9"})
	//fmt.Println(len(a))
	//fmt.Println(len("9.9.9.9"))
	//
	//if a != "9.9.9.9" {
	//	t.Error("不通过")
	//}
}