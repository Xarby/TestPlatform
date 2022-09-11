package Util

import "testing"

func TestCheckIP(t *testing.T) {
	if res := CheckIP("1.2.3.4") ;res == false{
		t.Fatal("预期结果正常，返回却是异常")
	}
}