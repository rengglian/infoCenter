package main

import (
	"fmt"

	"infoCenter/internal/radio"
)

var dir200 = radio.New("http://192.168.1.204/", "1234", "Mon Jan _2 2006", "15:04")

func main() {

	fmt.Println(dir200.GetPowerState())
	fmt.Println(dir200.GetDateString())
	fmt.Println(dir200.GetTimeString())
	fmt.Println(dir200.GetStationString())
	fmt.Println(dir200.GetSongString())
	fmt.Println(dir200.GetValideModes())
	fmt.Println(dir200.GetNavList())
	fmt.Println(dir200.GetMode())
	fmt.Println(dir200.SetMode(0))
	fmt.Println(dir200.GetPresetList())
	fmt.Println(dir200.SetPowerState(true))

}
