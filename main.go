package main

import (
	"human/hmdb"
	"human/mylog"
)

func main() {
	mylog.CurTime()

	mzs := mylog.ReadFile("conf/list.txt")
	for _, v := range mzs {

		body := mylog.ReadHtml("html/search/hmdb-" + v + ".html")
		if len(body) == 0 {
			continue
		}
		hmdb.ParseHtml(v, []byte(body))
	}
	mylog.Over()
	// hmdb.ParseHmdbDetail("HMDB0031127")
}
