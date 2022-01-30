package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
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
	fmt.Printf(fileName+"中需要搜索的个数:%d\n", lines)
	return
}

func ReadFile(fileName string) {
	chanNum := GetFileLines(fileName)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buf := bufio.NewScanner(f)
	chs := make([]chan int, chanNum)

	for i := 0; i < chanNum; i++ {
		if !buf.Scan() {
			break
		}
		line := buf.Text()
		chs[i] = make(chan int)

		go Gethtml(line, chs[i])

		time.Sleep(500 * time.Millisecond)
	}
	for _, ch := range chs {
		<-ch
	}
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func WriteFile(fileName string, msg string) {

	if FileExist(fileName) {
		fmt.Println(fileName + "已经下载")
		return
	}

	fileHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	defer fileHandle.Close()
	buf := bufio.NewWriterSize(fileHandle, len(msg))

	buf.WriteString(msg)

	err = buf.Flush()
	if err != nil {
		log.Println("flush error :", err)
	}
	fmt.Println(fileName + "下载成功")
}

func Gethtml(search string, ch chan int) {
	defer func() {
		ch <- 1
	}()
	if len(search) == 0 {
		return
	}
	aa := strings.ReplaceAll(search, " ", "%20")
	// aa := strings.ReplaceAll(aa, " ", "%20")
	url := "https://hmdb.ca/unearth/q?utf8=%E2%9C%93&query=" + aa + "&searcher=metabolites&button="
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("下载失败")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	WriteFile("./download/"+search+".html", string(body))
}

func Gethtml2(url string) []byte {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("下载失败")
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return body
}

func Usage() {
	fmt.Printf(
		"------------------------------------------------\n" +
			"  1. list.txt中添加需要下载的M/Z值(每行一个)\n" +
			"  2. 双击运行getHmdb.exe(一般需要多次运行才能下载完毕)\n" +
			"  3. 最后双击运行parseHmdb.exe\n" +
			"------------------------------------------------\n")
}

func begin() string {
	fmt.Println("请选择从HMDB数据库下载数据:")
	fmt.Println("\t\t: 下载搜索列表信息(list.txt)")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("请按任意键开始:")
	b := make([]byte, 1)
	os.Stdin.Read(b)
	return list
}
func over() {
	fmt.Println("请按任意键结束...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

const (
	list string = "list.txt"
)

func main() {
	usl := begin()
	ReadFile(usl)
	ParseAndSave()
	over()
}

/////////////////////////////////////////////////
func ParseAndSave() {
	mzs := ReadFileSlice(list)
	for _, v := range mzs {
		body := ReadHtml("./download/" + v + ".html")
		// body := ReadHtml("./download/2-Methylhexanoic acid.html")
		if len(body) == 0 {
			continue
		}
		ParseHtml(v, []byte(body))
		break
	}
}
func ReadFileSlice(fileName string) []string {

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
func ReadHtml(fileName string) string {
	if !FileExist(fileName) {
		return ""
	}
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return string(f)
}

func ParseHtml(search string, html []byte) {
BEGIN:
	if html == nil {
		return
	}
	html1 := strings.Replace(string(html), "\n", "", -1)
	body := strings.Replace(html1, "  ", "", -1)

	reItems := regexp.MustCompile(`<div class="unearth-search-hit">(.*?)</div>`)
	reid := regexp.MustCompile(`">HMDB(.*?)</a>`)

	items := reItems.FindAllString(body, -1)

	fmt.Printf(">>>>>>>>>:搜索条件=%s 搜索结果共计%d个条目\n", search, len(items))

	res := make([]string, 0, len(items))
	for index, item := range items {
		tds := reid.FindAllStringSubmatch(item, -1)
		for _, v := range tds {
			res = append(res, "HMDB"+v[1])
		}
		fmt.Printf("----------------第%d个-------------------\n", index)

	}
	fmt.Printf("%+v", res)

	reNextPage := regexp.MustCompile(`<li class="next_page">(.*?)</li>`)
	reNextPageUrl := regexp.MustCompile(`href="(.*?)">Next`)
	next := reNextPage.FindAllString(body, -1)
	// fmt.Printf("%+v", next)
	if len(next) > 0 {
		url := reNextPageUrl.FindAllStringSubmatch(next[1], -1)
		if len(url) > 0 {
			nextUrl := fmt.Sprintf("https://hmdb.ca%s", url[0][1])
			html = Gethtml2(nextUrl)
			goto BEGIN
		}
	}

	// mylog.WriteExls2(search, infos)
}
