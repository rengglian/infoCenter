package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

var templates = template.Must(template.ParseFiles("./templates/infoCenter.html"))
var validPath = regexp.MustCompile("^/(infoCenter|shutdown)/([a-zA-Z0-9]+)$")

var server = &http.Server{Addr: ":8080", Handler: nil}

// Page is the main struct
type Page struct {
	Title       string
	Info        Information
	RequestInfo ReqInfo
}

// Information is the main struct
type Information struct {
	DateTime    string
	StationName string
	SongName    string
	Status      string
}

// ReqInfo is the main struct
type ReqInfo struct {
	ElapasedStr string
}

// FsapiResponse is the main struct
type FsapiResponse struct {
	XMLName xml.Name `xml:"fsapiResponse"`
	Status  string   `xml:"status"`
	Value   []Value  `xml:"value"`
}

// Value is the main struct
type Value struct {
	XMLName xml.Name `xml:"value"`
	C8Array string   `xml:"c8_array"`
	U8      uint8    `xml:"u8"`
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func loadInformation() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	ipAddr := "http://192.168.0.143/"
	var info Information
	if getPowerState(ipAddr) == 1 {
		info.DateTime = getDateTimeString(ipAddr)
		info.StationName = getStationString(ipAddr)
		info.SongName = getSongString(ipAddr)
		info.Status = "Power On"
	} else {
		info.DateTime = "None"
		info.StationName = "None"
		info.SongName = "None"
		info.Status = "Power Off"
	}

	reqinfo.ElapasedStr = time.Since(start).String()
	return &Page{Title: "InfoCenter", Info: info, RequestInfo: reqinfo}, nil
}

func shutdownHandler(w http.ResponseWriter, r *http.Request, title string) {
	err := exec.Command("sudo", "systemctl", "poweroff").Run()
	checkError("Exec Error", err)
}

func loadInformationHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadInformation()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "infoCenter", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTimeString(ipAddr string) string {

	var str strings.Builder
	str.WriteString(ipAddr)
	str.WriteString("fsapi/GET/netRemote.sys.clock.localTime?pin=1234")
	resp, err := http.Get(str.String())
	checkError("HTTP New Request Error", err)
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	checkError("XML Decode Error", err)

	return fsapiResponse.Value[0].C8Array
}

func getDateString(ipAddr string) string {

	var str strings.Builder
	str.WriteString(ipAddr)
	str.WriteString("fsapi/GET/netRemote.sys.clock.localDate?pin=1234")
	resp, err := http.Get(str.String())
	checkError("HTTP New Request Error", err)
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	checkError("XML Decode Error", err)

	return fsapiResponse.Value[0].C8Array
}

func getDateTimeString(ipAddr string) string {
	var str strings.Builder
	str.WriteString(getDateString(ipAddr))
	str.WriteString(getTimeString(ipAddr))
	dateTime, err := time.Parse("20060102150405", str.String())
	checkError("Pase Time Failed", err)
	return dateTime.Format("Mon Jan _2 2006 - 15:04")
}

func getStationString(ipAddr string) string {

	var str strings.Builder
	str.WriteString(ipAddr)
	str.WriteString("fsapi/GET/netRemote.play.info.name?pin=1234")
	resp, err := http.Get(str.String())
	checkError("HTTP New Request Error", err)
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	checkError("XML Decode Error", err)

	return fsapiResponse.Value[0].C8Array
}

func getSongString(ipAddr string) string {

	var str strings.Builder
	str.WriteString(ipAddr)
	str.WriteString("fsapi/GET/netRemote.play.info.text?pin=1234")
	resp, err := http.Get(str.String())
	checkError("HTTP New Request Error", err)
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	checkError("XML Decode Error", err)

	return fsapiResponse.Value[0].C8Array
}

func getPowerState(ipAddr string) uint8 {

	var str strings.Builder
	str.WriteString(ipAddr)
	str.WriteString("fsapi/GET/netRemote.sys.power?pin=1234")
	resp, err := http.Get(str.String())
	checkError("HTTP New Request Error", err)
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	checkError("XML Decode Error", err)

	return fsapiResponse.Value[0].U8
}

func main() {

	http.HandleFunc("/infoCenter/", makeHandler(loadInformationHandler))
	http.HandleFunc("/shutdown/", makeHandler(shutdownHandler))
	http.Handle("/addons/", http.StripPrefix("/addons/", http.FileServer(http.Dir("addons"))))

	fmt.Println("Reporting is Running")

	go func() {
		err := server.ListenAndServe()
		if err.Error() == "http: Server closed" {

		} else {
			checkError("Listen And Serve Error", err)
		}

	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	checkError("Shutdown", err)
	fmt.Println("Good Bye")
}
