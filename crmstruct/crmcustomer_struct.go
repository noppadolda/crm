package crmcustomer_struct

import (
	"fmt"
)

// นำเข้า package fmt มาใช้งาน

type Crminbounddata struct {
	Customerno              int
	Crm_inboundfnname       string
	Suhistorno              int
	Geometry     struct {
		Area      int
		Perimeter int
	}
}

func init() {
	fmt.Println("crm struct package initialized")

}
func Area(len, wid float64) float64 {
	area := len * wid
	return area
}
