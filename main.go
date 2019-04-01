package main

import (
	"docker"
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello")
	t := time.Duration(10 * time.Second)
	tt := time.NewTicker(t)
	for {
		select {
		case <-tt.C:
			thread()
		}
	}
}

func thread() {
	v := docker.GetDiskStat(10)
	for _, aa := range v {
		fmt.Printf("name is:%v  value is:%v\n", aa.Name, aa.Value)
	}
}
