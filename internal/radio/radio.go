package radio

import (
	"encoding/xml"
	"net/http"

	"time"

	"infoCenter/internal/check"
)

const (
	getTimeCmd       string = "fsapi/GET/netRemote.sys.clock.localTime?pin="
	getDateCmd       string = "fsapi/GET/netRemote.sys.clock.localDate?pin="
	getStationCmd    string = "fsapi/GET/netRemote.play.info.name?pin="
	getSongCmd       string = "fsapi/GET/netRemote.play.info.text?pin="
	getPowerStateCmd string = "fsapi/GET/netRemote.sys.power?pin="
	getModeCmd       string = "fsapi/GET/netRemote.sys.mode?pin="
	getValidModesCmd string = "fsapi/LIST_GET_NEXT/netRemote.sys.caps.validModes/-1?pin="
	getNavListCmd    string = "fsapi/LIST_GET_NEXT/netRemote.nav.list/-1?pin="
	getMuteStateCmd  string = "fsapi/GET/netRemote.sys.audio.mute?pin="
	setMuteStateCmd  string = "fsapi/SET/netRemote.sys.audio.mute?pin="
	getVolumeCmd     string = "fsapi/GET/netRemote.sys.audio.volume?pin="

	setVolumeCmd  string = "fsapi/SET/netRemote.sys.audio.volume?pin="
	getPresetsCmd string = "fsapi/LIST_GET_NEXT/netRemote.nav.presets/-1?pin="

	setPresetsCmd  string = "fsapi/SET/netRemote.nav.action.selectPreset?pin="
	getNavStateCmd string = "fsapi/GET/netRemote.nav.state?pin="
	setNavStateCmd string = "fsapi/SET/netRemote.nav.state?pin="

	setPowerStateCmd string = "fsapi/SET/netRemote.sys.power?pin="
	setModeCmd       string = "fsapi/SET/netRemote.sys.mode?pin="

	maxItemCmd string = "&maxItems=65536"
	valueCmd   string = "&value="
)

// FsapiResponse is the main struct
type FsapiResponse struct {
	XMLName xml.Name `xml:"fsapiResponse"`
	Status  string   `xml:"status"`
	Values  []Value  `xml:"value"`
	Items   []Item   `xml:"item"`
}

// Value is the main struct
type Value struct {
	XMLName xml.Name `xml:"value"`
	C8Array string   `xml:"c8_array"`
	U8      uint8    `xml:"u8"`
	U32     uint32   `xml:"u32"`
}

// Item is the main struct
type Item struct {
	XMLName xml.Name `xml:"item"`
	Fields  []Field  `xml:"field"`
}

// Field is the main struct
type Field struct {
	XMLName xml.Name `xml:"field"`
	C8Array string   `xml:"c8_array"`
	U8      uint8    `xml:"u8"`
}

//DeviceConfig Struct
type DeviceConfig struct {
	ip         string
	pin        string
	dateFormat string
	timeFormat string
}

//New func
func New(ip string, pin string, dateFormat string, timeFormat string) *DeviceConfig {
	d := &DeviceConfig{}
	d.ip = ip
	d.pin = pin
	d.dateFormat = dateFormat
	d.timeFormat = timeFormat
	return d
}

//GetTimeString function
func (d *DeviceConfig) GetTimeString() string {

	str := d.ip + getTimeCmd + d.pin
	resp, err := http.Get(str)
	if err != nil {
		t := time.Now()
		return t.Format(d.timeFormat)
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	dateTime, err := time.Parse("150405", fsapiResponse.Values[0].C8Array)
	check.Error("Time Parse Failed: ", err)
	return dateTime.Format(d.timeFormat)
}

//GetDateString function
func (d *DeviceConfig) GetDateString() string {

	str := d.ip + getDateCmd + d.pin
	resp, err := http.Get(str)
	if err != nil {
		t := time.Now()
		return t.Format(d.dateFormat)
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	dateTime, err := time.Parse("20060102", fsapiResponse.Values[0].C8Array)
	check.Error("Time Parse Failed: ", err)
	return dateTime.Format(d.dateFormat)
}

//GetStationString function
func (d *DeviceConfig) GetStationString() string {

	str := d.ip + getStationCmd + d.pin
	resp, err := http.Get(str)
	if err != nil {
		return "No Connection"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].C8Array
}

//GetSongString function
func (d *DeviceConfig) GetSongString() string {

	str := d.ip + getSongCmd + d.pin
	resp, err := http.Get(str)
	if err != nil {
		return "No Connection"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].C8Array
}

//GetPowerState function
func (d *DeviceConfig) GetPowerState() uint8 {

	str := d.ip + getPowerStateCmd + d.pin
	resp, err := http.Get(str)
	if err != nil {
		return 0
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].U8

}

//GetMuteState function
func (d *DeviceConfig) GetMuteState() uint8 {

	str := d.ip + getMuteStateCmd + d.pin
	resp, err := http.Get(str)
	if err != nil {
		return 0
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].U8

}

//SetPowerState function
func (d *DeviceConfig) SetPowerState(value bool) string {

	bitSetVar := "0"
	if value {
		bitSetVar = "1"
	}

	str := d.ip + setPowerStateCmd + d.pin + valueCmd + bitSetVar
	resp, err := http.Get(str)
	if err != nil {
		return "Not Connected"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	return fsapiResponse.Status

}

//SetMuteState function
func (d *DeviceConfig) SetMuteState(value bool) string {

	bitSetVar := "0"
	if value {
		bitSetVar = "1"
	}

	str := d.ip + setMuteStateCmd + d.pin + valueCmd + bitSetVar
	resp, err := http.Get(str)
	if err != nil {
		return "Not Connected"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	return fsapiResponse.Status

}

//GetValideModes fucntion
func (d *DeviceConfig) GetValideModes() []string {
	var respStr []string
	str := d.ip + getValidModesCmd + d.pin + maxItemCmd
	resp, err := http.Get(str)
	if err != nil {
		respStr = append(respStr, "No Connection")
		return respStr
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	for _, item := range fsapiResponse.Items {

		var tmpStr string
		selectable := "false\t"
		if item.Fields[1].U8 == 1 {
			selectable = "true\t"
		}

		streamable := "false\t"
		if item.Fields[3].U8 == 1 {
			streamable = "true\t"
		}
		tmpStr += "ID:\t"
		tmpStr += item.Fields[0].C8Array
		tmpStr += "\tSelectable:\t"
		tmpStr += selectable
		tmpStr += "Label:\t"
		tmpStr += item.Fields[2].C8Array
		tmpStr += "\tStreamable:\t"
		tmpStr += streamable
		respStr = append(respStr, tmpStr)
	}

	return respStr
}

//GetNavList fucntion
func (d *DeviceConfig) GetNavList() []string {
	var respStr []string
	str := d.ip + getNavListCmd + d.pin + maxItemCmd
	resp, err := http.Get(str)
	if err != nil {
		respStr = append(respStr, "No Connection")
		return respStr
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	for _, item := range fsapiResponse.Items {
		var tmpStr string
		selectable := "false\t"
		if item.Fields[1].U8 == 1 {
			selectable = "true\t"
		}

		streamable := "false\t"
		if item.Fields[3].U8 == 1 {
			streamable = "true\t"
		}
		tmpStr += "ID:\t"
		tmpStr += item.Fields[0].C8Array
		tmpStr += "\tSelectable:\t"
		tmpStr += selectable
		tmpStr += "Label:\t"
		tmpStr += item.Fields[2].C8Array
		tmpStr += "\tStreamable:\t"
		tmpStr += streamable
		respStr = append(respStr, tmpStr)
	}

	return respStr
}

//GetPresetList fucntion
func (d *DeviceConfig) GetPresetList() []string {

	var respStr []string

	var fsapiResponse FsapiResponse
	fsapiResponse.Status = d.SetNavState(1)

	if fsapiResponse.Status == "FS_OK" {

		str := d.ip + getPresetsCmd + d.pin + maxItemCmd
		resp, err := http.Get(str)
		if err != nil {
			respStr = append(respStr, "No Connection")
			return respStr
		}

		err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
		check.Error("XML Decode Error", err)

		for _, item := range fsapiResponse.Items {
			respStr = append(respStr, item.Fields[0].C8Array)
		}
		d.SetNavState(0)
	}
	return respStr
}

//GetNavState function
func (d *DeviceConfig) GetNavState() uint8 {

	str := d.ip + getNavStateCmd + d.pin

	resp, err := http.Get(str)
	if err != nil {
		return 0
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].U8

}

//SetNavState function
func (d *DeviceConfig) SetNavState(value int) string {
	var valueStr string
	switch value {
	case 0:
		valueStr = "0"
	case 1:
		valueStr = "1"
	default:
		valueStr = "0"
	}

	str := d.ip + setNavStateCmd + d.pin + valueCmd + valueStr
	resp, err := http.Get(str)
	if err != nil {
		return "Not Connected"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	return fsapiResponse.Status

}

//SetPresetSate function
func (d *DeviceConfig) SetPresetSate(value int) string {
	var valueStr string
	var fsapiResponse FsapiResponse
	fsapiResponse.Status = d.SetNavState(1)
	if fsapiResponse.Status == "FS_OK" {

		switch value {
		case 0:
			valueStr = "0"
		case 1:
			valueStr = "1"
		case 2:
			valueStr = "2"
		case 3:
			valueStr = "3"
		case 4:
			valueStr = "4"
		case 5:
			valueStr = "5"
		case 6:
			valueStr = "6"
		case 7:
			valueStr = "7"
		case 8:
			valueStr = "8"
		case 9:
			valueStr = "9"
		default:
			valueStr = "0"
		}

		str := d.ip + setPresetsCmd + d.pin + valueCmd + valueStr
		resp, err := http.Get(str)
		if err != nil {
			return "Not Connected"
		}

		err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
		check.Error("XML Decode Error", err)

		d.SetNavState(0)
	}

	return fsapiResponse.Status

}

//GetMode function
func (d *DeviceConfig) GetMode() uint32 {

	str := d.ip + getModeCmd + d.pin

	resp, err := http.Get(str)
	if err != nil {
		return 0
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].U32

}

//GetVolume function
func (d *DeviceConfig) GetVolume() uint8 {

	str := d.ip + getVolumeCmd + d.pin

	resp, err := http.Get(str)
	if err != nil {
		return 0
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)
	return fsapiResponse.Values[0].U8

}

//SetMode function
func (d *DeviceConfig) SetMode(value int) string {

	var valueStr string
	switch value {
	case 0:
		valueStr = "0"
	case 2:
		valueStr = "2"
	case 3:
		valueStr = "3"
	case 4:
		valueStr = "4"
	case 5:
		valueStr = "5"
	default:
		valueStr = "0"
	}

	str := d.ip + setModeCmd + d.pin + valueCmd + valueStr

	resp, err := http.Get(str)
	if err != nil {
		return "Not Connected"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	return fsapiResponse.Status

}

//SetVolume function
func (d *DeviceConfig) SetVolume(value uint8) string {

	var valueStr string
	switch value {
	case 0:
		valueStr = "0"
	case 1:
		valueStr = "1"
	case 2:
		valueStr = "2"
	case 3:
		valueStr = "3"
	case 4:
		valueStr = "4"
	case 5:
		valueStr = "5"
	case 6:
		valueStr = "6"
	case 7:
		valueStr = "7"
	case 8:
		valueStr = "8"
	case 9:
		valueStr = "9"
	default:
		valueStr = "3"
	}

	str := d.ip + setVolumeCmd + d.pin + valueCmd + valueStr

	resp, err := http.Get(str)
	if err != nil {
		return "Not Connected"
	}
	var fsapiResponse FsapiResponse
	err = xml.NewDecoder(resp.Body).Decode(&fsapiResponse)
	check.Error("XML Decode Error", err)

	return fsapiResponse.Status

}

//IncreaseVolume function
func (d *DeviceConfig) IncreaseVolume() {
	actVol := d.GetVolume()
	actVol++
	d.SetVolume(actVol)
}

//DecreaseVolume function
func (d *DeviceConfig) DecreaseVolume() {
	actVol := d.GetVolume()
	actVol--
	d.SetVolume(actVol)
}
