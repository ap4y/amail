package main

import (
	"log"

	"ap4y.me/cloud-mail/notmuch"
)

func main() {
	c, err := notmuch.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	threads, err := c.Search("folder:INBOX", 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(threads)
	log.Println("=========")

	thread, err := c.Show("thread:" + threads[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", thread[0][0])

	count, err := c.Count("folder:INBOX")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(count)

	reply, err := c.Reply("id:20210521190054.eim6yn774i7ltc44@jetsam.local", "all")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", reply)

	if err := c.Index(); err != nil {
		log.Fatal(err)
	}
}
