package main

import _ "github.com/farseer810/file-manager/dao"

func main() {
	server := GetServer()
	server.Run()
}
