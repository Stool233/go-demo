package cron

import (
	"demo/http_demo"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"time"
)

func StartCron() {
	s := gocron.NewScheduler()
	s.Every(1).Minutes().Do(http_demo.GetDictionaries)

	sc := s.Start()

	go wait(s, sc)
	<-sc

}

func wait(s *gocron.Scheduler, sc chan bool) {
	time.Sleep(8 * time.Minute)
	fmt.Println("All task removed")
	close(sc) // close the channel
}
