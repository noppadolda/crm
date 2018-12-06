package handler

import (
	"ats_eng_api/encryption"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/sjson"

	"github.com/tidwall/gjson"

	dbutil "ats_eng_api/database"

	"github.com/gorilla/mux"
)

type FunctionMap struct {
	Function  interface{}
	Authorize bool
}

var strAppUserExecute = "app_user_execute"

/* var m = map[string]interface{}{
	"campaign/getCampaign":    getCampaign,
	"campaign/setCampaign":    setCampaign,
	"campaign/deleteCampaign": deleteCampaign,
	"campaign/checkPricePlan": checkPricePlan,
	"getGeneric":              getGeneric,
	"getToken":                GetToken,
	"Login":                   Login,
	"test":                    test,
	"Encrypt":                 Encrypt,
	"Decrypt":                 Decrypt,
} */
var m = map[string]FunctionMap{
	"campaign/getCampaign":    {Function: getCampaign, Authorize: false},
	"campaign/setCampaign":    {Function: setCampaign, Authorize: false},
	"campaign/deleteCampaign": {Function: deleteCampaign, Authorize: false},
	"campaign/checkPricePlan": {Function: checkPricePlan, Authorize: false},
	"campaign/getPricePlan":   {Function: getPricePlan, Authorize: false},
	"campaign/setPricePlan":   {Function: setPricePlan, Authorize: false},
	"getGeneric":              {Function: getGeneric, Authorize: false},
	"getToken":                {Function: GetToken, Authorize: false},
	"CheckToken":              {Function: CheckTokenAndRenew, Authorize: false},
	"Login":                   {Function: Login, Authorize: false},
	"test":                    {Function: test, Authorize: false},
	"Encrypt":                 {Function: Encrypt, Authorize: false},
	"Decrypt":                 {Function: Decrypt, Authorize: false},
}

var respCode = map[int]string{
	200: "Success",
	201: "No data found.",
	211: "Login fail.",
	221: "Token not authorize.",
	222: "Token time out.",
	223: "Renew Token.",
	500: "Server Error.",
	501: "Method not found.",
}

var resp string = "{" +
	"\"status_code\" : 500," +
	"\"status_description\" : \"Error\"," +
	"\"data\" : null" +
	"}"

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	vars := mux.Vars(r)
	var jsonData string
	authPass := true
	//var err error
	//fmt.Println("HandleRequest")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var funct = vars["function"]
	if vars["category"] != "" {
		funct = vars["category"] + "/" + funct
	}
	//fmt.Println(funct)
	if m[funct].Function != nil {
		/* if len(query) > 0 {
			query := r.URL.Query()
			jsonData, _ = callDynamically(funct, r, query)
		} else {
			jsonData, _ = callDynamically(funct, r)
		} */
		//CheckToken(r.Header.Get("Authorization"))
		//fmt.Println(GetTokenUser(r.Header.Get("Authorization")))
		if m[funct].Authorize {
			if !CheckToken(r.Header.Get("Authorization")) {
				resp, _ = sjson.Set(resp, "status_code", 221)
				resp, _ = sjson.Set(resp, "status_description", respCode[221])
				resp, _ = sjson.Set(resp, "data", nil)
				jsonData = resp
				authPass = false
			}
		}
		if authPass {
			jsonData, _ = callDynamically(funct, r)
		}
	} else {
		resp, _ = sjson.Set(resp, "status_code", 501)
		resp, _ = sjson.Set(resp, "status_description", respCode[501])
		jsonData = resp
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, jsonData)
}

func callDynamically(name string, r *http.Request, args ...interface{}) (string, error) {
	argsLength := len(args)
	var jsonData string
	var err error
	//fmt.Println(r.Header.Get("Authorization"))
	//fmt.Println(argsLength)
	if argsLength == 1 {
		switch v := args[0].(type) {
		case url.Values:
			query := args[0].(url.Values)
			jsonData, err = m[name].Function.(func(url.Values) (string, error))(query)
		case string:
			fmt.Printf("%q is %v bytes long\n", v, len(v))
		default:
			jsonData, err = m[name].Function.(func() (string, error))()
		}

	} else {
		query := r.URL.Query()
		if len(query) > 0 {
			jsonData, err = m[name].Function.(func(url.Values) (string, error))(query)
		} else {
			b, err := ioutil.ReadAll(r.Body)

			defer r.Body.Close()
			if err != nil {
				resp, _ = sjson.Set(resp, "status_code", 500)
				resp, _ = sjson.Set(resp, "status_description", respCode[500])
				return resp, err
			}
			bSize := len(b)
			user := GetTokenUser(r.Header.Get("Authorization"))
			sb, _ := sjson.Set(string(b), strAppUserExecute, user)
			sb, _ = sjson.Set(string(b), "token", r.Header.Get("Authorization"))
			if bSize == 0 {
				jsonData, err = m[name].Function.(func() (string, error))()
			} else {
				jsonData, err = m[name].Function.(func(string) (string, error))(string(sb))
			}
		}

	}
	/*if argsLength == 3 {
		for _, key := range args[2].MapKeys() {
			strct := args[2].MapIndex(key)
			fmt.Println(key.Interface(), strct.Interface())
		}
		jsonData, err := m[name].(func(map[string][]string) (string, error))(args[2])
	} else {
		jsonData, err := m[name].(func() (string, error))()
	} */
	return jsonData, err
}

func getGeneric(query url.Values) (string, error) {
	code := query["code"]
	sql := "select * from ats_generic_config where config_code = '" + code[0] + "' "
	jsonData, err := dbutil.SelectData(sql)
	fmt.Println(jsonData)
	if err != nil {
		resp, _ = sjson.Set(resp, "status_code", 500)
		resp, _ = sjson.Set(resp, "status_description", respCode[500])
		return resp, err
	}
	resp, _ = sjson.Set(resp, "status_code", 200)
	resp, _ = sjson.Set(resp, "status_description", respCode[200])
	m, _ := gjson.Parse(gjson.Get(jsonData, "result.0.config_value").String()).Value().(map[string]interface{})
	jsonData, _ = sjson.Set(resp, "data", m)
	return jsonData, nil
}

func test() (string, error) {

	resp, _ = sjson.Set(resp, "status_code", 200)
	resp, _ = sjson.Set(resp, "status_description", respCode[200])
	return resp, nil
}

func Encrypt(data string) (string, error) {
	rt, _ := encryption.EncryptBase64([]byte(gjson.Get(data, "data").String()), []byte(gjson.Get(data, "key").String()))
	//rt, _ := encryption.Encrypt([]byte(gjson.Get(data, "data").String()), []byte(gjson.Get(data, "key").String()))
	resp, _ = sjson.Set(resp, "status_code", 200)
	resp, _ = sjson.Set(resp, "status_description", respCode[200])
	resp, _ = sjson.Set(resp, "data", string(rt))
	return resp, nil
}

func Decrypt(data string) (string, error) {
	rt, _ := encryption.DecryptBase64([]byte(gjson.Get(data, "data").String()), []byte(gjson.Get(data, "key").String()))
	//rt, _ := encryption.Decrypt([]byte(gjson.Get(data, "data").String()), []byte(gjson.Get(data, "key").String()))
	resp, _ = sjson.Set(resp, "status_code", 200)
	resp, _ = sjson.Set(resp, "status_description", respCode[200])
	resp, _ = sjson.Set(resp, "data", string(rt))
	return resp, nil
}
