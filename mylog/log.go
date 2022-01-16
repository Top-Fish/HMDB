package mylog

import (
	"bufio"
	"fmt"
	"human/global"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/360EntSecGroup-Skylar/excelize"
	excelize "github.com/xuri/excelize/v2"
)

func WriteLog(fileName string, msg string) {
	fileHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriterSize(fileHandle, len(msg))

	buf.WriteString(msg)

	err = buf.Flush()
	if err != nil {
		log.Println("flush error :", err)
	}
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

var xlsxName string

func CurTime() {
	now := time.Now()
	curTime := now.Format("2006年01月02日15时04分05秒")
	xlsxName = "./excel/" + curTime + ".xlsx"

}

func WriteExls(mz string, ga []global.Article) {

	var xlsx *excelize.File
	xlsx, err := excelize.OpenFile(xlsxName)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "cannot find the file specified") {
			fmt.Println("创建日志文件" + xlsxName)
			xlsx = excelize.NewFile()
		}
	}

	index := xlsx.NewSheet("M/Z=" + mz)
	data := map[string]string{
		"A1": "HMDB-ID",
		"B1": "Name",
		"C1": "Formula",
		"D1": "Monoisotopic Mass",
		"E1": "Adduct",
		"F1": "Adduct M/Z",
		"G1": "Delta",
		"H1": "CCS",
		"I1": "Class",
		"J1": "Subclass",
		"K1": "DirectParent",
	}
	for k, v := range data {
		xlsx.SetCellValue("M/Z="+mz, k, v)
	}
	for i, v := range ga {
		data := map[string]string{
			"A" + strconv.Itoa(i+2): v.Compound,
			"B" + strconv.Itoa(i+2): v.Name,
			"C" + strconv.Itoa(i+2): v.Formula,
			"D" + strconv.Itoa(i+2): v.Monoisotopic,
			"E" + strconv.Itoa(i+2): v.Adduct,
			"F" + strconv.Itoa(i+2): v.AdductMZ,
			"G" + strconv.Itoa(i+2): v.Delta,
			"H" + strconv.Itoa(i+2): v.CCS,
		}
		for k, v := range data {
			xlsx.SetCellValue("M/Z="+mz, k, v)
		}
	}

	//设置默认打开的表单
	xlsx.SetActiveSheet(index)

	//保存文件到指定路径
	err = xlsx.SaveAs(xlsxName)
	if err != nil {
		log.Fatal(err)
	}
}

var gaNum int = 0

func CreateExcelFile(fileName string) {
	fmt.Println("创建Excel文件:" + xlsxName)
	xlsx := excelize.NewFile()
	data := map[string]string{
		"A1": "M/Z",
		"B1": "HMDB-ID",
		"C1": "Name",
		"D1": "Formula",
		"E1": "Monoisotopic Mass",
		"F1": "Adduct",
		"G1": "Adduct M/Z",
		"H1": "Delta",
		"I1": "CCS",
		"J1": "Class",
		"K1": "Subclass",
		"L1": "DirectParent",
		"M1": "Synonyms Name",
		"N1": "IUPAC Name",
		"O1": "Samples",
	}
	for k, v := range data {
		xlsx.SetCellValue("Sheet1", k, v)

	}
	// xlsx.SetRowHeight("Sheet1", 1, 50)
	xlsx.SetRowStyle("Sheet1", 1, 1, 2)
	xlsx.SetColWidth("Sheet1", "A", "A", 10)
	xlsx.SetColWidth("Sheet1", "B", "B", 15)
	xlsx.SetColWidth("Sheet1", "C", "C", 30)
	xlsx.SetColWidth("Sheet1", "D", "I", 15)
	xlsx.SetColWidth("Sheet1", "J", "O", 30)
	//设置默认打开的表单
	index := xlsx.GetSheetIndex("Sheet1")
	xlsx.SetActiveSheet(index)

	//保存文件到指定路径
	err := xlsx.SaveAs(xlsxName)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteExls2(mz string, ga []global.Article) {

	var xlsx *excelize.File
	xlsx, err := excelize.OpenFile(xlsxName)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "cannot find the file specified") {
			CreateExcelFile(xlsxName)
			xlsx, _ = excelize.OpenFile(xlsxName)
		}
	}

	if len(ga) == 0 {
		data := map[string]string{
			"A" + strconv.Itoa(gaNum+2): mz,
			"B" + strconv.Itoa(gaNum+2): "-",
			"C" + strconv.Itoa(gaNum+2): "-",
			"D" + strconv.Itoa(gaNum+2): "-",
			"E" + strconv.Itoa(gaNum+2): "-",
			"F" + strconv.Itoa(gaNum+2): "-",
			"G" + strconv.Itoa(gaNum+2): "-",
			"H" + strconv.Itoa(gaNum+2): "-",
			"I" + strconv.Itoa(gaNum+2): "-",
			"J" + strconv.Itoa(gaNum+2): "-",
			"K" + strconv.Itoa(gaNum+2): "-",
			"L" + strconv.Itoa(gaNum+2): "-",
			"M" + strconv.Itoa(gaNum+2): "-",
			"N" + strconv.Itoa(gaNum+2): "-",
			"O" + strconv.Itoa(gaNum+2): "-",
		}
		for k, v := range data {
			xlsx.SetCellValue("Sheet1", k, v)
		}
		gaNum++
	} else {
		for _, v := range ga {
			data := map[string]string{
				"A" + strconv.Itoa(gaNum+2): mz,
				"B" + strconv.Itoa(gaNum+2): v.Compound,
				"C" + strconv.Itoa(gaNum+2): v.Name,
				"D" + strconv.Itoa(gaNum+2): v.Formula,
				"E" + strconv.Itoa(gaNum+2): v.Monoisotopic,
				"F" + strconv.Itoa(gaNum+2): v.Adduct,
				"G" + strconv.Itoa(gaNum+2): v.AdductMZ,
				"H" + strconv.Itoa(gaNum+2): v.Delta,
				"I" + strconv.Itoa(gaNum+2): v.CCS,
				"J" + strconv.Itoa(gaNum+2): v.Class,
				"K" + strconv.Itoa(gaNum+2): v.SubClass,
				"L" + strconv.Itoa(gaNum+2): v.Parent,
				"M" + strconv.Itoa(gaNum+2): v.Synonyms,
				"N" + strconv.Itoa(gaNum+2): v.IUPAC,
				"O" + strconv.Itoa(gaNum+2): v.Samples,
			}
			for k, v := range data {
				xlsx.SetCellValue("Sheet1", k, v)
			}
			gaNum++
		}
	}

	//设置默认打开的表单
	index := xlsx.GetSheetIndex("Sheet1")
	xlsx.SetActiveSheet(index)

	//保存文件到指定路径
	err = xlsx.SaveAs(xlsxName)
	if err != nil {
		log.Fatal(err)
	}
}

func Over() {
	fmt.Println("请按任意键结束...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
