package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ReqParam struct {
	SearchMode     string   `json:"ms_search_ion_mode"`
	AdductType     []string `json:"adduct_type"`
	Tolerance      string   `json:"tolerance"`
	ToleranceUnits string   `json:"tolerance_units"`
	CcsPredictors  string   `json:"ccs_predictors"`
	CcsTolerance   string   `json:"ccs_tolerance"`
}

func (rp ReqParam) Trans2String() string {
	str := make([]byte, 0, 1024)

	{
		ss := fmt.Sprintf("&ms_search_ion_mode=%s", rp.SearchMode)
		str = append(str, []byte(ss)...)
	}

	{
		for _, v := range rp.AdductType {
			ss := "&adduct_type%5B%5D=" + v
			str = append(str, []byte(ss)...)
		}

	}
	{
		ss := fmt.Sprintf("&tolerance=%s", rp.Tolerance)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&tolerance_units=%s", rp.ToleranceUnits)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&ccs_predictors=%s", rp.CcsPredictors)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&ccs_tolerance=%s", rp.CcsTolerance)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&commit=%s", "Search")
		str = append(str, []byte(ss)...)
	}
	return strings.ReplaceAll(string(str), "+", "%2B")
}

func ReadReqParamConfig(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	r := io.Reader(f)

	param := &ReqParam{}
	if err = json.NewDecoder(r).Decode(param); err != nil {
		panic(err)
	}

	return param.Trans2String()
}

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
	fmt.Printf(fileName+"中需要下载M/Z个数:%d\n", lines)
	return
}

func ReadFile(fileName string) {
	fileName = "conf/" + fileName
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
		switch fileName {
		case "conf/list.txt":
			go Gethtml(line, chs[i])
		case "conf/detailist.txt":
			go GetHmdbDetail(line, chs[i])
		default:
			fmt.Println("文件名称有误")
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	for _, ch := range chs {
		<-ch
	}
	fmt.Println("\t\tGame Over!!!")
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

func Gethtml(mz string, ch chan int) {
	defer func() {
		ch <- 1
	}()
	if len(mz) == 0 {
		return
	}
	url := "https://hmdb.ca/spectra/ms/search"
	method := "POST"

	// payload := strings.NewReader("query_masses=" + mz + "&ms_search_ion_mode=negative&adduct_type%5B%5D=M-H&adduct_type%5B%5D=M%2BCl&tolerance=5&tolerance_units=ppm&ccs_predictors=&ccs_tolerance=&commit=Search")
	payload := strings.NewReader("query_masses=" + mz + ReqConf)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
	// fmt.Println(string(body))
	WriteFile("./html/search/hmdb-"+mz+".html", string(body))
}

func Usage() {
	fmt.Printf(
		"------------------------------------------------\n" +
			"  1. list.txt中添加需要下载的M/Z值(每行一个)\n" +
			"  2. 双击运行getHmdb.exe(一般需要多次运行才能下载完毕)\n" +
			"  3. 最后双击运行parseHmdb.exe\n" +
			"------------------------------------------------\n")
}

func GetHmdbDetail(hmdbId string, ch chan int) []byte {
	defer func() {
		ch <- 1
	}()
	if len(hmdbId) != 11 {
		return nil
	}
	url := "https://hmdb.ca/metabolites/" + hmdbId
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	// req.Header.Add("Cookie", "_hmdb_session=cUU1OGZWNWpxSXBxazJuUVo0dDFjdCt1OEhXTkVrS1JEM014V0pQeEFEMGhmbDJQcWJ5d2FBU1lBT1BhM0RScnBsSnU3N016bW8vTWxaaU1ZNnV3M3JnYjgvY1lkenkvVmpzeW42WjJqMjQvdVNCWVBqdXF1VFl3NEp2U3JLbVJLTEE3MlVQZlNUMjdQUW9ZczUwMENRPT0tLUdGVDVxZnd3cGJFUUw4NTQ0NUh5RWc9PQ%3D%3D--627cd870610f64251d2769a8a996f26071595c2e")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println(hmdbId + "下载失败")
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	WriteFile("./html/metabolites/"+hmdbId+".html", string(body))
	return nil
}
func begin() string {
	fmt.Println("请选择从HMDB数据库下载数据:")
	fmt.Println("\t\t1: 下载搜索列表信息(list.txt)")
	fmt.Println("\t\t2: 下载详细信息(detailist.txt)")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("您的选择是:")
	b := make([]byte, 1)
	os.Stdin.Read(b)

	switch b[0] {
	case '1':
		return list
	case '2':
		return detail
	default:
		fmt.Println("选择错误, 程序结束!!!")
		os.Exit(0)
	}
	return ""
}
func over() {
	fmt.Println("请按任意键结束...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

const (
	list     string = "list.txt"
	detail   string = "detailist.txt"
	confFile string = "conf/param.json"
)

var ReqConf string

func main() {
	ReqConf = ReadReqParamConfig(confFile)
	usl := begin()
	ReadFile(usl)
	over()
}
