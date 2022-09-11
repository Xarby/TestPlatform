package Util

import "strings"

func SplitFormat(instring string) string {
	return 	strings.Replace(strings.Replace(instring,"\n","",-1),"\t","",-1)
}
