package structs

//LocationList struct
type LocationList struct {
	Items []LocationEntry
}

//LocationEntry struct
type LocationEntry struct {
	Location       string
	Voltage        float64
	Altitude       float64
	Temperature    float64
	MinTemperature float64
	MinTimeStr     string
	MaxTemperature float64
	MaxTimeStr     string
	Pressure       float64
	Humidity       float64
	LastTimeStr    string
	Icon           string
	ID             string
}

//TempDetails struct
type TempDetails struct {
	Points []TempEntry
}

//TempEntry Struct
type TempEntry struct {
	TimeStr     string
	Temperature float64
	Humidity    float64
}

//AddItem exported funtion
func (entry *LocationList) AddItem(item LocationEntry) []LocationEntry {
	entry.Items = append(entry.Items, item)
	return entry.Items
}

//AddItem exported funtion
func (entry *TempDetails) AddItem(item TempEntry) []TempEntry {
	entry.Points = append(entry.Points, item)
	return entry.Points
}
