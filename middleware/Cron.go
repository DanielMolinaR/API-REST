package middleware

import "github.com/robfig/cron/v3"

var c *cron.Cron

func init(){
	c = cron.New()

	c.Start()
	
	c.AddFunc("0 9 * * 1-5", sendEmail())
}
