package util

import "fmt"

const DEBUGGING bool = true

func DebugPrint(str string) {
	if DEBUGGING {
		fmt.Println(str);
	}
}