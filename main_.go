
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pysrc/rest"
	"github.com/pysrc/simsql"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// 允许操作的表与字段
var db_allow = map[string][]string{

	"demo_city": []string{

		"demo_city_id", // NullAble=NO | 城市ID

		"demo_city_name", // NullAble=YES | 城市名

	},

	"demo_user": []string{

		"demo_id", // NullAble=NO | 用户ID

		"demo_name", // NullAble=YES | 用户名

		"demo_date", // NullAble=YES | 注册日期

		"demo_city_name", // NullAble=YES | 城市ID

	},

}

// 判断该表是否允许
func table_allow(table string) bool {
	_, ok := db_allow[table]
	return ok
}

// 判断该字段是否允许
func field_allow(table, field string) bool {
	if table_allow(table) {
		for _, v := range db_allow[table] {
			if v == field {
				return true
			}
		}
	}
	return false
}

func SendJson(args interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(args)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/my_test?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var queryDb = func(w http.ResponseWriter, r *http.Request, sqlNode simsql.SqlNode) {
		fmt.Println(sqlNode)
		var res []interface{}
		rs, _ := db.Query(sqlNode.Query, sqlNode.Args...)
		keys, _ := rs.Columns()
		leng := len(keys)
		for rs.Next() {
			pvals := make([]interface{}, leng)
			vals := make([]string, leng)
			for i := 0; i < leng; i++ {
				pvals[i] = &vals[i]
			}
			rs.Scan(pvals...)
			im := make(map[string]string)
			for i := 0; i < leng; i++ {
				im[keys[i]] = vals[i]
			}
			res = append(res, im)
		}
		SendJson(res, w)
	}
	var execDb = func(w http.ResponseWriter, r *http.Request, sqlNode simsql.SqlNode) {
		fmt.Println(sqlNode.Query, sqlNode.Args)
		tx, err := db.Begin()
		if err != nil {
			return
		}
		res, err := tx.Exec(sqlNode.Query, sqlNode.Args...)
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
		rs, err := res.RowsAffected()
		if err != nil {
			return
		}
		id, _ := res.LastInsertId()
		SendJson(map[string]int64{"rows": rs, "lastId": id}, w)
	}
	var router rest.Router
	router.Route("GET", "/api/v1/:table/:schema", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		if !table_allow(params["table"]) {
			http.NotFound(w, r)
			return
		}
		schemas := strings.Split(params["schema"], "-")
		for _, v := range schemas {
			if !field_allow(params["table"], v) {
				http.NotFound(w, r)
				return
			}
		}
		per, err := strconv.ParseInt(r.FormValue("per"), 10, 64)
		if err != nil || per <= 0 {
			per = 30
		}
		page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1
		}
		sqlNode := simsql.Query("and", schemas, params["table"], nil, per, page)
		queryDb(w, r, sqlNode)
	})
	router.Route("GET", "/api/v1/:table/:field/:value/:schema", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		if !table_allow(params["table"]) {
			http.NotFound(w, r)
			return
		}
		schemas := strings.Split(params["schema"], "-")
		for _, v := range schemas {
			if !field_allow(params["table"], v) {
				http.NotFound(w, r)
				return
			}
		}
		if !field_allow(params["table"], params["field"]) {
			http.NotFound(w, r)
			return
		}
		pa := map[string]interface{}{
			params["field"]: params["value"],
		}
		per, err := strconv.ParseInt(r.FormValue("per"), 10, 64)
		if err != nil || per <= 0 {
			per = 30
		}
		page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1
		}
		queryDb(w, r, simsql.Query("and", strings.Split(params["schema"], "-"), params["table"], pa, per, page))
	})
	router.Route("PUT", "/api/v1/:table/:field/:value", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		if !table_allow(params["table"]) {
			http.NotFound(w, r)
			return
		}
		if !field_allow(params["table"], params["field"]) {
			http.NotFound(w, r)
			return
		}
		pa := map[string]interface{}{
			params["field"]: params["value"],
		}
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)
		execDb(w, r, simsql.Update("and", params["table"], data, pa))
	})
	router.Route("POST", "/api/v1/:table", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		if !table_allow(params["table"]) {
			http.NotFound(w, r)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)
		execDb(w, r, simsql.Insert(params["table"], data))
	})
	router.Route("DELETE", "/api/v1/:table/:field/:value", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		if !table_allow(params["table"]) {
			http.NotFound(w, r)
			return
		}
		if !field_allow(params["table"], params["field"]) {
			http.NotFound(w, r)
			return
		}
		pa := map[string]interface{}{
			params["field"]: params["value"],
		}
		execDb(w, r, simsql.Delete("or", params["table"], pa))
	})
	router.Run("127.0.0.1:8080")
}

