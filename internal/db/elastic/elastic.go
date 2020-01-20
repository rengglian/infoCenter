package db

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"infoCenter/internal/check"
	"infoCenter/internal/structs"

	"github.com/elastic/go-elasticsearch/v7"
)

//Connection external
type Connection struct {
	*elasticsearch.Client
}

//New external
func New(cert string, key string, server string, user string, pw string) *Connection {
	cer, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println(err)
		return nil
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			server,
		},
		Username: user,
		Password: pw,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cer},
			},
		},
	}
	es, _ := elasticsearch.NewClient(cfg)
	return &Connection{es}
}

//ListEnviroLocations func
func (es *Connection) ListEnviroLocations() []string {

	var r map[string]interface{}
	queryStr := `
	{
		"aggs" : { 
			"location":{ 
				"terms" : { 
							"field" : "Location.keyword", 
							"size":10000 
						}
					}
				}
	}
	`
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("healthstate-enviro"),
		es.Search.WithBody(strings.NewReader(queryStr)),
		es.Search.WithTrackTotalHits(false),
		es.Search.WithSize(0),
	)

	check.Error("Search Error", err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&r)
	check.Error("JSON Decode Error", err)

	s := []string{}
	for _, value := range r["aggregations"].(map[string]interface{})["location"].(map[string]interface{})["buckets"].([]interface{}) {
		s = append(s, value.(map[string]interface{})["key"].(string))
	}
	sort.Strings(s)
	return s
}

//GetLast24hr func
func (es *Connection) GetLast24hr(location string) structs.LocationEntry {

	var r map[string]interface{}
	var locationEntry structs.LocationEntry
	queryStr := `
	{
		"query": {
			  "bool": {
			  "must": [
				{
			"range" : {
				"@timestamp": {
					"gte" : "now-24h",
					"lt" :  "now"
				  }
			  }
			  },
					{
					  "match_phrase": {
					  "Location.keyword": {
					  "query": "%s"
					  }
				  }
				  }
				  
			  ]
			  }
		  },
		  "sort": [
		{
		"@timestamp": {
		"order": "desc"
		}
		}
		]
	  }
	`
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("healthstate-enviro"),
		es.Search.WithBody(strings.NewReader(fmt.Sprintf(queryStr, location))),
		es.Search.WithTrackTotalHits(false),
		es.Search.WithSize(288),
	)

	check.Error("Search Error", err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&r)
	check.Error("JSON Decode Error", err)

	locationEntry.Location = location
	locationEntry.Voltage = 0.0
	locationEntry.Altitude = 0.0
	locationEntry.Temperature = 0.0
	locationEntry.MinTemperature = 100.0
	locationEntry.MaxTemperature = -100.0
	locationEntry.Pressure = 0.0
	locationEntry.Humidity = 0.0
	locationEntry.LastTimeStr = "unknown"
	locationEntry.MinTimeStr = "unknown"
	locationEntry.MaxTimeStr = "unknown"
	locationEntry.Icon = "bed.svg"

	lenResponse := float64(len(r["hits"].(map[string]interface{})["hits"].([]interface{})))
	if lenResponse > 0 {
		tmpTime, err := time.Parse(time.RFC3339, r["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"].(string))
		var tmpMinTime time.Time
		var tmpMaxTime time.Time
		locationEntry.Temperature = r["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["_source"].(map[string]interface{})["Temperature"].(float64)
		locationEntry.Humidity = r["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["_source"].(map[string]interface{})["Humidity"].(float64)
		check.Error("Time Parse Error", err)
		for _, value := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			tmpTemperature := value.(map[string]interface{})["_source"].(map[string]interface{})["Temperature"].(float64)
			if tmpTemperature > locationEntry.MaxTemperature {
				locationEntry.MaxTemperature = tmpTemperature
				tmpMaxTime, err = time.Parse(time.RFC3339, value.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"].(string))
				check.Error("Time Parse Error", err)
			}
			if tmpTemperature < locationEntry.MinTemperature {
				locationEntry.MinTemperature = tmpTemperature
				tmpMinTime, err = time.Parse(time.RFC3339, value.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"].(string))
				check.Error("Time Parse Error", err)
			}

		}
		switch locationEntry.Location {
		case "Office":
			locationEntry.Icon = "desk.svg"
			locationEntry.ID = "songCell"
		case "Living Room":
			locationEntry.Icon = "television.svg"
			locationEntry.ID = "stationCell"
		default:
			locationEntry.Icon = "bed.svg"
			locationEntry.ID = "blueCell"
		}

		locationEntry.LastTimeStr = tmpTime.Add(time.Hour * 1).Format("Mon Jan _2 15:04:05")
		locationEntry.MinTimeStr = tmpMinTime.Add(time.Hour * 1).Format("Mon Jan _2 15:04:05")
		locationEntry.MaxTimeStr = tmpMaxTime.Add(time.Hour * 1).Format("Mon Jan _2 15:04:05")
	}
	return locationEntry
}

//GetLastWeek func
func (es *Connection) GetLastWeek(location string) structs.TempDetails {

	var r map[string]interface{}
	queryStr := `
	{
		"_source": [
		  "@timestamp",
		  "Temperature",
		  "Humidity"
		],
		"sort": [
		  {
			"@timestamp": {
			  "order": "asc"
			}
		  }
		],
		"query": {
		  "bool": {
			"must": [
			  {
				"exists": {
				  "field": "Temperature"
				}
			  },
			  {
				"wildcard": {
				  "Location.keyword": {
					"value": "%s"
				  }
				}
			  },
			  {
				"range": {
				  "@timestamp": {
					"gt": "now-1w",
					"lt": "now"
				  }
				}
			  }
			]
		  }
		}
	  }
	`
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("healthstate-enviro"),
		es.Search.WithBody(strings.NewReader(fmt.Sprintf(queryStr, location))),
		es.Search.WithTrackTotalHits(false),
		es.Search.WithSize(0),
		es.Search.WithSize(2100),
	)

	check.Error("Search Error", err)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&r)
	check.Error("JSON Decode Error", err)

	items := []structs.TempEntry{}
	entry := structs.TempDetails{Points: items}

	lenResponse := float64(len(r["hits"].(map[string]interface{})["hits"].([]interface{})))
	if lenResponse > 0 {
		for _, value := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {

			var item structs.TempEntry
			item.TimeStr = value.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"].(string)
			item.Temperature = value.(map[string]interface{})["_source"].(map[string]interface{})["Temperature"].(float64)
			item.Humidity = value.(map[string]interface{})["_source"].(map[string]interface{})["Humidity"].(float64)
			entry.AddItem(item)
		}
	}

	return entry
}
