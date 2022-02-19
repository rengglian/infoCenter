package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"time"

	"infoCenter/internal/check"
	"infoCenter/internal/config"
	db "infoCenter/internal/db/elastic"
	"infoCenter/internal/radio"
	"infoCenter/internal/structs"
)

var templates = template.Must(template.ParseFiles("./templates/home.html",
	"./templates/radio.html",
	"./templates/radioControl.html",
	"./templates/about.html",
	"./templates/lab.html",
	"./templates/tempGraph.html",
	"./templates/badi.html"))

var validPath = regexp.MustCompile("^/(home|radio|about|lab|tempGraph|radioControl|badi)/([a-zA-Z0-9_ ]+)$")

var cfg = config.New("./config/config.yaml")

var server = &http.Server{Addr: cfg.GetServerPort(), Handler: nil}

var elastic = db.New(cfg.GetElasticCert(), cfg.GetElasticKey(), cfg.GetElasticServer(), cfg.GetElasticUser(), cfg.GetElasticPW())

var dir200 = radio.New(cfg.GetRadioIP(), cfg.GetRadioPin(), "Mon Jan _2 2006", "15:04")

// Page is the main struct
type Page struct {
	Title        string
	Info         Information
	RequestInfo  ReqInfo
	LocationInfo structs.LocationList
	TempInfo     structs.TempDetails
}

// Information is the main struct
type Information struct {
	DateStr     string
	TimeStr     string
	Time        string
	StationName string
	SongName    string
	Status      string
	Presets     []string
	MuteState   string
	ModeWeb     string
	ModeMp3     string
	ModeDab     string
	ModeUkw     string
	ModeAux     string
}

// ReqInfo is the main struct
type ReqInfo struct {
	ElapasedStr string
}

func loadHomeInformation() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	var info Information
	info.DateStr = start.Format("Mon Jan _2 2006")
	info.TimeStr = start.Format("15:04")

	items := []structs.LocationEntry{}
	locations := structs.LocationList{Items: items}
	slist := elastic.ListEnviroLocations()
	for _, location := range slist {
		locations.AddItem(elastic.GetLast24hr(location))
	}
	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: "InfoCenter", Info: info, LocationInfo: locations, RequestInfo: reqinfo}, nil
}
func loadRadioInformation() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()

	var info Information
	info.DateStr = start.Format("Mon Jan _2 2006")
	info.TimeStr = start.Format("15:04")
	if dir200.GetPowerState() == 1 {
		info.StationName = dir200.GetStationString()
		info.SongName = dir200.GetSongString()
		info.Status = "shutdownOn"
	} else {
		info.StationName = "NONE"
		info.SongName = "NONE"
		info.Status = "shutdownOff"
	}

	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: "InfoCenter", Info: info, RequestInfo: reqinfo}, nil
}

func loadRadioControl() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	var info Information
	info.Presets = nil
	if dir200.GetPowerState() == 1 {
		switch dir200.GetMode() {
		case 0:
			info.ModeWeb = "activeButton"
			info.ModeMp3 = ""
			info.ModeDab = ""
			info.ModeUkw = ""
			info.ModeAux = ""
		case 2:
			info.ModeWeb = ""
			info.ModeMp3 = "activeButton"
			info.ModeDab = ""
			info.ModeUkw = ""
			info.ModeAux = ""
		case 3:
			info.ModeWeb = ""
			info.ModeMp3 = ""
			info.ModeDab = "activeButton"
			info.ModeUkw = ""
			info.ModeAux = ""
		case 4:
			info.ModeWeb = ""
			info.ModeMp3 = ""
			info.ModeDab = ""
			info.ModeUkw = "activeButton"
			info.ModeAux = ""
		case 5:
			info.ModeWeb = ""
			info.ModeMp3 = ""
			info.ModeDab = ""
			info.ModeUkw = ""
			info.ModeAux = "activeButton"
		}
		info.Status = "shutdownOn"
		info.Presets = dir200.GetPresetList()
		if dir200.GetMuteState() == 1 {
			info.MuteState = "mute.png"
		} else {
			info.MuteState = "noMute.svg"
		}

	} else {
		info.ModeWeb = ""
		info.ModeMp3 = ""
		info.ModeDab = ""
		info.ModeUkw = ""
		info.ModeAux = ""
		info.Status = "shutdownOff"
		info.MuteState = "mute.png"
	}

	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: "InfoCenter", Info: info, RequestInfo: reqinfo}, nil
}

func loadAboutInformation() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	var info Information
	info.DateStr = start.Format("Mon Jan _2 2006")
	info.TimeStr = start.Format("15:04")

	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: "InfoCenter", Info: info, RequestInfo: reqinfo}, nil
}

func loadTempInformation(location string) (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	var info Information
	entry := elastic.GetLastWeek(location)

	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: location, Info: info, TempInfo: entry, RequestInfo: reqinfo}, nil
}

func loadLabInformation() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	var info Information
	info.DateStr = start.Format("Mon Jan _2 2006")
	info.TimeStr = start.Format("15:04")

	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: "InfoCenter", Info: info, RequestInfo: reqinfo}, nil
}

func loadTestInformation() (*Page, error) {
	var reqinfo ReqInfo
	start := time.Now()
	var info Information
	info.DateStr = start.Format("Mon Jan _2 2006")
	info.TimeStr = start.Format("15:04")

	reqinfo.ElapasedStr = time.Since(start).Truncate(time.Millisecond).String()
	return &Page{Title: "InfoCenter", Info: info, RequestInfo: reqinfo}, nil
}

func loadHomeHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadHomeInformation()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "home", p)
}

func loadRadioHandler(w http.ResponseWriter, r *http.Request, title string) {
	if title == "shutdown" {
		if dir200.GetPowerState() == 1 {
			dir200.SetPowerState(false)
		} else {
			dir200.SetPowerState(true)
		}
		http.Redirect(w, r, "/radio/index", http.StatusFound)
	}

	p, err := loadRadioInformation()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "radio", p)
}

func loadRadioControlHandler(w http.ResponseWriter, r *http.Request, title string) {

	switch title {
	case "web":
		dir200.SetMode(0)
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "mp3":
		dir200.SetMode(2)
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "dab":
		dir200.SetMode(3)
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "ukw":
		dir200.SetMode(4)
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "aux":
		dir200.SetMode(5)
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "mute":
		if dir200.GetMuteState() == 1 {
			dir200.SetMuteState(false)
		} else {
			dir200.SetMuteState(true)
		}
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "increaseVolume":
		dir200.IncreaseVolume()
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "decreaseVolume":
		dir200.DecreaseVolume()
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	case "preset":
		values, _ := r.URL.Query()["value"]
		nbr, _ := strconv.Atoi(values[0])
		dir200.SetPresetSate(nbr)
		http.Redirect(w, r, "/radioControl/index", http.StatusFound)
	}

	p, err := loadRadioControl()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "radioControl", p)
}

func loadAboutHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadAboutInformation()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "about", p)
}

func loadLabHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadAboutInformation()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "lab", p)
}

func loadTempGraphHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadTempInformation(title)
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "tempGraph", p)
}

func loadBadiHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadAboutInformation()
	if err != nil {
		title = "Information"
		http.Redirect(w, r, "/TODO/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "badi", p)
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

func main() {

	http.HandleFunc("/home/", makeHandler(loadHomeHandler))
	http.HandleFunc("/radio/", makeHandler(loadRadioHandler))
	http.HandleFunc("/about/", makeHandler(loadAboutHandler))
	http.HandleFunc("/lab/", makeHandler(loadLabHandler))
	http.HandleFunc("/tempGraph/", makeHandler(loadTempGraphHandler))
	http.HandleFunc("/radioControl/", makeHandler(loadRadioControlHandler))
	http.HandleFunc("/badi/", makeHandler(loadBadiHandler))
	http.Handle("/addons/", http.StripPrefix("/addons/", http.FileServer(http.Dir("addons"))))

	fmt.Println("Reporting is Running")
	fmt.Println(cfg.GetServerPort())
	fmt.Println(cfg.GetRadioIP())
	fmt.Println(cfg.GetRadioPin())

	go func() {
		err := server.ListenAndServe()
		if err.Error() == "http: Server closed" {

		} else {
			check.Error("Listen And Serve Error ", err)
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
	check.Error("Shutdown", err)
	fmt.Println("Good Bye")
}
