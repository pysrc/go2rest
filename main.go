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
	defer db.Close()
	if err != nil {
		return
	}
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
		per, err := strconv.ParseInt(r.FormValue("per"), 10, 64)
		if err != nil || per <= 0 {
			per = 30
		}
		page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1
		}
		sqlNode := simsql.Query("and", strings.Split(params["schema"], "-"), params["table"], nil, per, page)
		queryDb(w, r, sqlNode)
	})
	router.Route("GET", "/api/v1/:table/:field/:value/:schema", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
		pa := map[string]interface{}{
			params["field"]: params["value"],
		}
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)
		execDb(w, r, simsql.Update("and", params["table"], data, pa))
	})
	router.Route("POST", "/api/v1/:table", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)
		execDb(w, r, simsql.Insert(params["table"], data))
	})
	router.Route("DELETE", "/api/v1/:table/:field/:value", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		pa := map[string]interface{}{
			params["field"]: params["value"],
		}
		execDb(w, r, simsql.Delete("or", params["table"], pa))
	})
	router.Run("127.0.0.1:8080")
}
