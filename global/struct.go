package global

import "fmt"

type Conds struct {
	QueryMass     string
	IonMode       string
	AdductType    string
	Tolerance     string
	ToleranceUnit string
	CcsPredictor  string
	CcsTolerance  string
}

type Article struct {
	Compound     string
	Name         string
	Formula      string
	Monoisotopic string
	Adduct       string
	AdductMZ     string
	Delta        string
	CCS          string
	Class        string
	SubClass     string
	Parent       string
	Synonyms     string
	IUPAC        string
	Samples      string
}

func (e Article) Write2Exls() {

}

const (
	Compound int = iota
	Name
	Formula
	Monoisotopic
	Adduct
	AdductMZ
	Delta
	CCS
)

func (cds Conds) Trans2String() string {
	str := make([]byte, 0, 1024)
	// {
	// 	ss := fmt.Sprintf("authenticity_token=%s", "n7/akkNtsV+pJ325CAiGq4k1RmXJeWQ7w3A2vlPgDugyFDrvFMHfiABM/4qVn2ORShenObikGggnb3LHjhsmIw==")
	// 	str = append(str, []byte(ss)...)
	// }

	{
		ss := fmt.Sprintf("query_masses=%s", cds.QueryMass)
		str = append(str, []byte(ss)...)
	}

	{
		ss := fmt.Sprintf("&ms_search_ion_mode=%s", cds.IonMode)
		str = append(str, []byte(ss)...)
	}

	{
		ss := fmt.Sprintf("&adduct_type%%5B%%5D=%s", cds.AdductType)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&tolerance=%s", cds.Tolerance)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&tolerance_units=%s", cds.ToleranceUnit)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&ccs_predictors=%s", cds.CcsPredictor)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&ccs_tolerance=%s", cds.CcsTolerance)
		str = append(str, []byte(ss)...)
	}
	{
		ss := fmt.Sprintf("&commit=%s", "Search")
		str = append(str, []byte(ss)...)
	}
	return string(str)
}
