package main

import (
	"fmt"
	"time"

	"infoCenter/internal/config"
	db "infoCenter/internal/db/elastic"
	"infoCenter/internal/structs"
)

//ReqInfo struct
type ReqInfo struct {
	CreatedStr  string
	ElapasedStr string
}

var cfg = config.New("./config/config.yaml")

var elastic = db.New(cfg.GetElasticCert(), cfg.GetElasticKey(), cfg.GetElasticServer(), cfg.GetElasticUser(), cfg.GetElasticPW())

func testListEnviroLocations() {
	var reqinfo ReqInfo
	start := time.Now()
	slist := elastic.ListEnviroLocations()
	reqinfo.ElapasedStr = time.Since(start).String()
	reqinfo.CreatedStr = time.Now().Format("Mon Jan _2 15:04:05 2006")

	fmt.Println(reqinfo)
	fmt.Println(slist)
}

func testGetLast24hr() {
	var reqinfo ReqInfo
	items := []structs.LocationEntry{}
	location := structs.LocationList{Items: items}
	start := time.Now()
	slist := []string{"Office", "Living Room"}
	for _, device := range slist {
		location.AddItem(elastic.GetLast24hr(device))
	}
	reqinfo.ElapasedStr = time.Since(start).String()
	reqinfo.CreatedStr = time.Now().Format("Mon Jan _2 15:04:05 2006")
	fmt.Println(reqinfo)
	fmt.Println(location)
}

func testGetLastWeek() {
	var reqinfo ReqInfo
	start := time.Now()
	location := "Office"
	entry := elastic.GetLastWeek(location)
	reqinfo.ElapasedStr = time.Since(start).String()
	reqinfo.CreatedStr = time.Now().Format("Mon Jan _2 15:04:05 2006")
	fmt.Println(reqinfo)
	fmt.Println(len(entry.Points))
}

func main() {

	fmt.Println("testListEnviroLocations")
	testListEnviroLocations()
	fmt.Println("testGetLast24hr")
	testGetLast24hr()
	fmt.Println("testGetLastWeek")
	testGetLastWeek()
	fmt.Println("stopTest")
}
