package main

import (
	"fmt"
	"fstabmanager/fstab"
)

func main() {

	fmt.Println("FsTab Manager ...")
	ft, err1 := fstab.NewFsTabDb("/tmp/fstab.test")
	if err1 != nil {
		fmt.Println("Failed")
	}
	var err error

	err = ft.AddMount("/data", "/mnt/data", "none", "bind", "0", "0")
	if err != nil {
		fmt.Println(err)
	}
	err = ft.AddMount("/data", "/mnt/data", "none", "bind,noauto", "0", "0")
	if err != nil {
		fmt.Println(err)
	}

	err = ft.AddMount("/data1", "/mnt/data", "none", "bind,noauto1", "0", "0")
	if err != nil {
		fmt.Println(err)
	}
	err = ft.AddMount("/data", "/mnt/data1", "none", "bind,1noauto", "0", "0")
	if err != nil {
		fmt.Println(err)
	}
	ft.Save()

	for i, r := range ft.Records {
		fmt.Println(i, r.String())
	}
}
