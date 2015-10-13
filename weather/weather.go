package weather

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC1123Z, v)
	if err != nil {
		return err
	}
	*c = CustomTime{parse}
	return nil
}

type Range struct {
	Centigrade string `xml:"centigrade,attr"`
	Value      string `xml:",chardata"`
}

type Temperature struct {
	Unit   string  `xml:"unit,attr"`
	Ranges []Range `xml:"range"`
}

type Period struct {
	Hour  string `xml:"hour,attr"`
	Value string `xml:",chardata"`
}

type RainFallChance struct {
	Unit    string   `xml:"unit,attr"`
	Periods []Period `xml:"period"`
}

type Info struct {
	Date           string         `xml:"date,attr"`
	Weather        string         `xml:"weather"`
	Img            string         `xml:"img"`
	WeatherDetail  string         `xml:"weather_detail"`
	Wave           string         `xml:"wave"`
	Temperature    Temperature    `xml:"temperature"`
	RainFallChance RainFallChance `xml:"rainfallchance"`
}

type Geo struct {
	Long string `xml:"long"`
	Lat  string `xml:"lat"`
}

type Area struct {
	Id    string `xml:"id,attr"`
	Geo   Geo    `xml:"geo"`
	Infos []Info `xml:"info"`
}

type Pref struct {
	Id    string `xml:"id,attr"`
	Areas []Area `xml:"area"`
}

type WeatherForecast struct {
	XMLName        xml.Name   `xml:"weatherforecast"`
	Title          string     `xml:"title"`
	Link           string     `xml:"link"`
	Description    string     `xml:"description"`
	PubDate        CustomTime `xml:"pubDate"`
	Author         string     `xml:"author"`
	ManagingEditor string     `xml:"managingEditor"`
	Pref           Pref       `xml:"pref"`
}

func GetWeatherForecast() *WeatherForecast {
	w := WeatherForecast{}

	r, err := http.Get("http://www.drk7.jp/weather/xml/13.xml")
	if err != nil {
		log.Print(err)
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}

	err = xml.Unmarshal(b, &w)
	if err != nil {
		log.Print(err)
	}

	return &w
}
