package main

// นำเข้า package fmt มาใช้งาน
import (
	"crm/crmstruct" // referpath
	"crm/tvs_orderinfo"
	"fmt"
)

type Crm_request struct {
	OrderType          string
	tvs_orderinfoobj   tvs_orderinfo.Tvs_order
	crm_datainboundobj crmcustomer_struct.Crminbounddata
}

func main() {
	var Crm_requestobj Crm_request
	Crm_requestobj.OrderType = "!!!"
	Crm_requestobj.crm_datainboundobj.Crm_inboundfnname = "CREATECUSTOMER"
	Crm_requestobj.tvs_orderinfoobj.OrderType = "CUSTOMERORDER"
	fmt.Println(Crm_requestobj)
 
}
