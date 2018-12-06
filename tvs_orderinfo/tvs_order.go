package tvs_orderinfo

import (
	"fmt"
	"time"
)

// นำเข้า package fmt มาใช้งาน

type Tvs_order struct {
	TrackingNo string
	OrderType  string
	Level      string
	OrderDate  time.Time
	TVSNo      int
	MobileNo   string
	SerialNo   string
	Reference1 string
	Reference2 string
	Reference3 string
	Reference4 string
	Reference5 string
	Start    time.Time
	End time.Time
	Duration int 
	ProcessConfig[] Process
	Result_Code  string 
	Result_Desc  string 

}
type Process struct {
	 name string
	 start time.Time
	 end   time.Time 
}
func init() {
	fmt.Println("tvs_orderinfo initialized")

}
