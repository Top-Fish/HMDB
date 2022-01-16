package hmdb

import (
	"fmt"
	"regexp"
	"strings"

	"human/global"
	"human/mylog"
)

func GetConditions() (cds []global.Conds) {
	cds = make([]global.Conds, 0)
	cds = append(cds, global.Conds{
		QueryMass:     "175.01",
		IonMode:       "positive",
		AdductType:    "Unknown",
		Tolerance:     "0.05",
		ToleranceUnit: "Da",
		CcsPredictor:  "",
		CcsTolerance:  "",
	})
	return
}

///////////////////////////////////////////////////////////////////
func ParseHtml(mz string, html []byte) {
	html1 := strings.Replace(string(html), "\n", "", -1)
	body := strings.Replace(html1, "  ", "", -1)

	reItems := regexp.MustCompile(`<tr>(.*?)</tr>`)
	reArticle := regexp.MustCompile(`<td>(.*?)</td>`)
	reAdductMz := regexp.MustCompile(`<td>(.*?)<br`)
	reHmdbId := regexp.MustCompile(`">(.*?)</a>`)

	items := reItems.FindAllString(body, -1)

	fmt.Printf(">>>>>>>>>:M/Z=%s 搜索结果共计%d个条目\n", mz, len(items))

	infos := make([]global.Article, 0, len(items))

	for index, item := range items {
		tds := reArticle.FindAllStringSubmatch(item, -1)
		id := reHmdbId.FindAllStringSubmatch(tds[global.Compound][1], -1)

		amz := reAdductMz.FindAllStringSubmatch(tds[global.AdductMZ][0], -1)
		ss := strings.Replace(tds[global.Formula][1], "<sub>", "", -1)
		formula := strings.Replace(ss, "</sub>", "", -1)
		aaa := global.Article{
			Compound:     id[0][1],
			Name:         tds[global.Name][1],
			Formula:      formula,
			Monoisotopic: tds[global.Monoisotopic][1],
			Adduct:       tds[global.Adduct][1],
			AdductMZ:     amz[0][1],
			Delta:        tds[global.Delta][1],
			CCS:          tds[global.CCS][1],
		}
		details := ParseHmdbDetail(aaa.Compound)
		aaa.Class = details.Class
		aaa.SubClass = details.SubClass
		aaa.Parent = details.Parent
		aaa.Synonyms = details.Synonyms
		aaa.IUPAC = details.IUPAC
		aaa.Samples = details.Sample

		infos = append(infos, aaa)
		fmt.Printf("----------------第%d个-------------------\n", index)

	}

	mylog.WriteExls2(mz, infos)
}

////
//https://hmdb.ca/metabolites/HMDB0304953
//https://hmdb.ca/metabolites/HMDB0060015

type Details struct {
	Class    string
	SubClass string
	Parent   string
	Synonyms string
	IUPAC    string
	Sample   string
}

func ParseHmdbDetail(hmdbId string) (ret *Details) {

	ret = new(Details)

	body := mylog.ReadHtml("./html/metabolites/" + hmdbId + ".html")
	if len(body) == 0 {
		return
	}
	body = strings.Replace(string(body), "\n", "", -1)
	html := strings.Replace(body, "  ", "", -1)

	reTrTag := regexp.MustCompile(`<tr>(.*?)</tr>`)
	reClass := regexp.MustCompile(`">(.*?)<span class`)
	reSubClass := regexp.MustCompile(`">(.*?)<span class`)
	reParent := regexp.MustCompile(`">(.*?)<span class`)

	// reSynonyms := regexp.MustCompile(`<td>(.*?)</td>`)
	reIUPAC := regexp.MustCompile(`<td>(.*?)</td>`)
	reSample := regexp.MustCompile(`<div class="word-break-all">(.*?)</div>`)

	items := reTrTag.FindAllString(html, -1)
	for _, v := range items {
		if strings.Contains(v, "<th>Class</th>") {
			aa := reClass.FindAllStringSubmatch(v, -1)
			if len(aa) > 0 {
				ret.Class = aa[0][1]
			}
		} else if strings.Contains(v, "<th>Sub Class</th>") {
			aa := reSubClass.FindAllStringSubmatch(v, -1)
			if len(aa) > 0 {
				ret.SubClass = aa[0][1]
			}
		} else if strings.Contains(v, "<th>Direct Parent</th>") {
			aa := reParent.FindAllStringSubmatch(v, -1)
			if len(aa) > 0 {
				ret.Parent = aa[0][1]
			}
		} else if strings.Contains(v, "<th>Synonyms</th>") {
			if strings.Contains(v, "Not Available") {
				ret.Synonyms = "Not Available"
			} else {
				ret.Synonyms = "Find it by yourself!!!"
			}

		} else if strings.Contains(v, "<th>IUPAC Name</th>") {
			aa := reIUPAC.FindAllStringSubmatch(v, -1)
			if len(aa) > 0 {
				ret.IUPAC = aa[0][1]
			}
		} else if strings.Contains(v, "<th>SMILES</th>") {
			aa := reSample.FindAllStringSubmatch(v, -1)
			if len(aa) > 0 {
				ret.Sample = aa[0][1]
			}
		}

	}
	return
}
