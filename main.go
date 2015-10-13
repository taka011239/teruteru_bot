package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/robfig/cron"
	"github.com/taka011239/teruteru_bot/utils"
	"github.com/taka011239/teruteru_bot/weather"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	AREA_ID_TOKYO           = "東京地方"
	THRESHOLD_PRECIPITATION = "50"
	USER                    = "@manji0112"
)

func getArea(areas []weather.Area) *weather.Area {
	for _, area := range areas {
		if area.Id == AREA_ID_TOKYO {
			return &area
		}
	}
	return nil
}

func getTodayInfo(infos []weather.Info) *weather.Info {
	today := time.Now().Format("2006/01/02")
	for _, info := range infos {
		if info.Date == today {
			return &info
		}
	}
	return nil
}

func isRainy(chance *weather.RainFallChance) bool {
	for i, period := range chance.Periods {
		// 00-06時は無視
		if i == 0 {
			continue
		}

		if period.Value >= THRESHOLD_PRECIPITATION {
			return true
		}
	}
	return false
}

func needUmbrella(f *weather.WeatherForecast) bool {
	if area := getArea(f.Pref.Areas); area != nil {
		if info := getTodayInfo(area.Infos); info != nil {
			return isRainy(&info.RainFallChance)
		}
	}
	log.Printf("Can not get the information")
	return false
}

var c *utils.Config

func init() {
	log.Print("**** Load Config ****")
	c = utils.NewConfig()
	c.Load()
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGTERM,
		syscall.SIGINT)

	api := anaconda.NewTwitterApi(c.AccessToken, c.AccessTokenSecret)

	c := cron.New()
	c.AddFunc("0 0 7 * * *", func() {
		log.Print("**** START ****")

		f := weather.GetWeatherForecast()
		if needUmbrella(f) {
			m := USER + " 傘持って行ってね♡"
			tweet, _ := api.PostTweet(m, nil)
			log.Print(tweet.Text)
		}

		log.Print("**** END ****")
	})
	c.Start()

	select {
	case s := <-sigCh:
		log.Printf("Recv: %v\n", s)
		os.Exit(0)
	}
}
