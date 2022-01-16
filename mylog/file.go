package mylog

import (
	"bufio"
	"fmt"
	"os"
)

func GetFileLines(fileName string) (lines int) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	lines = 0
	buf := bufio.NewScanner(f)

	for buf.Scan() {
		lines++
	}
	fmt.Printf("list.txt中需要下载M/Z个数:%d\n", lines)
	return
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReadFile(fileName string) []string {

	chanNum := GetFileLines(fileName)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	buf := bufio.NewScanner(f)

	names := make([]string, 0, chanNum)

	for i := 0; i < chanNum; i++ {
		if !buf.Scan() {
			break
		}
		line := buf.Text()
		names = append(names, line)
	}

	return names
}
