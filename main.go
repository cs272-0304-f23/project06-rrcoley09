package main

import (
	"flag"
	"time"
)

func main() {
	//exit := make(chan os.Signal, 1)
	//signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	url := "https://www.ucsc.edu/wp-sitemap.xml"
	//url := "https://www.npr.org/live-updates/sitemap.xml"

	var (
		siteMap = flag.Bool("sitemap", false, "Determines if input is sitemap")
		pastDB  = flag.Bool("pastDB", false, "Determines if past db should be used")
	)
	flag.Parse()

	go serve_web()

	if *siteMap {
		crawl(url, true, *pastDB)
	} else {
		crawl(url, false, *pastDB)
	}
	//<-exit

	for {
		time.Sleep(100 * time.Millisecond)
	}

}
