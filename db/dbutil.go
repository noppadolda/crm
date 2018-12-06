package database

import (
	config "ats_eng_api/config"
	"ats_eng_api/encryption"
	"database/sql"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	_ "gopkg.in/goracle.v2"
)

var db_key = "@t$eng@p!DBp@sS!"

func SelectData(sqlStmnt string) (string, error) {
	size := 0
	dsn, _ := encryption.DecryptBase64([]byte(config.GetConfig("database")), []byte(db_key))

	db, err := sql.Open("goracle", string(dsn))
	if err != nil {
		log.Errorln("dbutil:SelectData:sql.Open : " + err.Error())
		//fmt.Println("SelectData :: " + err.Error())
		return "", err
	}
	rows, err := db.Query(sqlStmnt)
	if err != nil {
		log.Println("dbutil:SelectData:rows : " + err.Error())
		//fmt.Println("rows :: " + err.Error())
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		log.Println("dbutil:SelectData:column : " + err.Error())
		//fmt.Println("column :: " + err.Error())
		return "", err
	}
	count := len(columns)
	//tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	/* for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[strings.ToLower(col)] = v
		}
		tableData = append(tableData, entry)
		size++
	} */
	jArr := "\"result\":[]"
	for rows.Next() {
		jOut := ""
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		//entry := make(map[string]interface{})
		cols := len(columns)
		for cols-1 >= 0 {
			var v interface{}
			val := values[cols-1]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			jOut, _ = sjson.Set(jOut, strings.ToLower(columns[cols-1]), v)
			cols--
		}
		//fmt.Println(jOut)
		m, _ := gjson.Parse(jOut).Value().(map[string]interface{})
		if size == 0 {
			jArr, _ = sjson.Set(jArr, "result.0", m)
		} else {
			jArr, _ = sjson.Set(jArr, "result.-1", m)
		}

		size++
	}
	//jArr = gjson.Get(jArr, "data").String()

	defer db.Close()
	//fmt.Println(gjson.Get(jArr, "result.#"))
	/* if gjson.Get(jArr, "result.#").Int() == 1 {
		jArr = gjson.Get(jArr, "result.0").String()
	} else {
		jArr = "{" + gjson.Get(jArr, "result").String() + "}"
	} */
	//fmt.Println(jArr)
	/* jsonData, err := json.Marshal(tableData)
	if size == 1 {
		jsonData, err = json.Marshal(tableData[0])
	} */

	return jArr, nil
}

func ExecutetData(sqlStmnt string) (string, error) {
	dsn, _ := encryption.DecryptBase64([]byte(config.GetConfig("database")), []byte(db_key))

	db, err := sql.Open("goracle", string(dsn))
	if err != nil {
		log.Errorln("dbutil:ExecutetData:sql.Open : " + err.Error())
		//fmt.Println("SelectData :: " + err.Error())
		return "", err
	}
	result, err := db.Exec(sqlStmnt)
	jreSult := ""
	if err != nil {
		log.Println("dbutil:ExecutetData:Exec : " + err.Error())
		jreSult, _ = sjson.Set(jreSult, "rows_success", 0)
		jreSult, _ = sjson.Set(jreSult, "message", err.Error())
		jreSult, _ = sjson.Set(jreSult, "success", false)
		//fmt.Println("rows :: " + err.Error())
	} else {
		//fmt.Println(result.RowsAffected)
		effect, _ := result.RowsAffected()
		jreSult, _ = sjson.Set(jreSult, "rows_success", effect)
		jreSult, _ = sjson.Set(jreSult, "message", "Success.")
		jreSult, _ = sjson.Set(jreSult, "success", true)
	}
	defer db.Close()

	return jreSult, err
}

func ExecuteTransaction(sqlStmnt []string) (string, error) {

	dsn, _ := encryption.DecryptBase64([]byte(config.GetConfig("database")), []byte(db_key))
	isError := false
	var successRow int64 = 0
	jreSult := ""
	db, err := sql.Open("goracle", string(dsn))
	if err != nil {
		log.Errorln("dbutil:ExecuteTransaction:sql.Open : " + err.Error())
		//fmt.Println("SelectData :: " + err.Error())
		return "", err
	}
	tx, _ := db.Begin()
	for _, sql := range sqlStmnt {
		result, err := tx.Exec(sql)
		if err != nil {
			log.Errorln("dbutil:ExecuteTransaction:sql.Exec : " + err.Error())
			//fmt.Println("SelectData :: " + err.Error())
			isError = true
			break
		}
		aff, _ := result.RowsAffected()
		successRow += aff
	}
	if !isError {
		tx.Commit()
		jreSult, _ = sjson.Set(jreSult, "rows_success", successRow)
		jreSult, _ = sjson.Set(jreSult, "message", "Success.")
		jreSult, _ = sjson.Set(jreSult, "success", true)
	} else {
		tx.Rollback()
		jreSult, _ = sjson.Set(jreSult, "rows_success", 0)
		jreSult, _ = sjson.Set(jreSult, "message", err.Error())
		jreSult, _ = sjson.Set(jreSult, "success", false)
		log.Errorln("dbutil:ExecuteTransaction:sql.Exec : Rollback")
	}
	defer db.Close()

	return jreSult, nil
}
