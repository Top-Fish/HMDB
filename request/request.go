package request

import (
	"fmt"
	"human/global"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetHmdbMain() {

	url := "https://hmdb.ca/spectra/ms/search"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "applitdsion/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error")
	}

	if res.StatusCode != 200 {
		fmt.Println("GetHmdbMain:获取Cookie错误")
		return
	}
	fmt.Println(res.Header)

	global.CurCookie = res.Header.Get("Set-Cookie")
	global.CurETag = res.Header.Get("ETag")

}

func GetHtml(params string) {

	url := "https://hmdb.ca/spectra/ms/search"
	method := "POST"

	payload := strings.NewReader("query_masses=175.01&ms_search_ion_mode=positive&adduct_type%5B%5D=Unknown&tolerance=0.05&tolerance_units=Da&ccs_predictors=&ccs_tolerance=&commit=Search")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println(res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	fmt.Println(res.StatusCode)
	// hmdb.ParseHtml(body)
}

func GetHmdbDetail(hmdbId string) []byte {
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// mylog.WriteLog("hmdb"+hmdbId+".html", string(body))

	return body
}
