package main

import (
	"reflect"
	"testing"
)

func TestImageSearch(t *testing.T) {
	tests := []struct {
		term string
		want []Hits
	}{
		{"santa",
			[]Hits{
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/05/4-19-23-Merrill-Spring-CL-002-1024x683.jpg", "", "https://www.ucsc.edu/about/", 0.002084088437839797},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2022/10/you-are-loved-1024x682.jpg", "", "https://www.ucsc.edu/about/", 0.002084088437839797},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2022/10/bridge-1024x768.jpg", "", "https://www.ucsc.edu/about/", 0.002084088437839797},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/05/Norcalmap.png", "", "https://www.ucsc.edu/about/", 0.004168176875679594},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2022/11/employment-1024x682.jpg", "", "https://www.ucsc.edu/about/", 0.004168176875679594},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2022/11/research-impact-1024x682.jpg", "", "https://www.ucsc.edu/about/", 0.004168176875679594},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/chancellor-edited-1024x577.jpg", "", "https://www.ucsc.edu/about/", 0.002084088437839797},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2022/11/sammy-move-in-1-1024x683.jpg", "", "https://www.ucsc.edu/about/", 0.002084088437839797},
				{"About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/monument-sign_ln-1024x768.jpg", "", "https://www.ucsc.edu/about/", 0.002084088437839797},
				{"Campus & Community – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/20171014-203857.jpg", "", "https://www.ucsc.edu/campus/", 0.0017623170638265268},
				{"Campus & Community – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/5-6-22-John-R-Lewis-CL-008.jpg", "", "https://www.ucsc.edu/campus/", 0.0035246341276530535},
				{"Our Mascot: Sammy the Banana Slug – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2022/12/Sammy-City-Hall-3-1-of-1-1024x772.jpg", "", "https://www.ucsc.edu/campus/mascot/", 0.00271771239513175},
				{"Overview – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/20220713-082437.jpg", "", "https://www.ucsc.edu/about/overview/", 0.002150537634408602},
				{"Overview – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/20220302-041553.jpg", "", "https://www.ucsc.edu/about/overview/", 0.002150537634408602},
				{"Overview – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/ucsc12_day1-63.jpg", "", "https://www.ucsc.edu/about/overview/", 0.002150537634408602},
				{"Overview – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/06/SAMMY-SILICON-VALLEY_RGB_300dpi-1024x1024.png", "", "https://www.ucsc.edu/about/overview/", 0.002150537634408602},
				{"Overview – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/IMG_3723.jpg", "", "https://www.ucsc.edu/about/overview/", 0.002150537634408602},
				{"Overview – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/06/6-10-23-Black-Grad-CL-029_SQR-1024x1024.jpg", "", "https://www.ucsc.edu/about/overview/", 0.002150537634408602},
				{"Visit UCSC – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/20171014-203857.jpg", "", "https://www.ucsc.edu/campus/visit/", 0.002385644642671922},
				{"Visit UCSC – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/07/5-6-22-John-R-Lewis-CL-008.jpg", "", "https://www.ucsc.edu/campus/visit/", 0.004771289285343844},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/cc-1-resize.jpg", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/ias-stage-photo.jpg", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/lick-lh-540.jpg", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/resources-icon.png", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/12-07-19-SCParade-CL-1-1.jpg", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/20220531-Aerial-Kresge-Construction-NEG-22-2-1024x682.jpg", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/water-icon.png", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/economy-icon.png", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/transportation-icon.png", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/house-icon.png", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
				{"UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/wp-content/uploads/2023/08/health-icon.png", "", "https://www.ucsc.edu/better-together/", 0.0009635525764558023},
			},
		},
	}

	url := "https://www.ucsc.edu/wp-sitemap.xml"

	go serve_web()
	crawl(url, true, false)

	for _, tc := range tests {
		term := "santa"
		got := search_image(term)

		if !reflect.DeepEqual(len(got), len(tc.want)) {
			t.Errorf("Want length: %v\nGot length: %v\n", len(tc.want), len(got))
		}

	}
}

func TestSnippet(t *testing.T) {
	tests := []struct {
		term string
		want []Hits
	}{
		{"thrive",
			[]Hits{
				{"About UC Santa Cruz – UC Santa Cruz", "", "Banana Slugs thrive in the redwoods", "https://www.ucsc.edu/campus/", 0.05463182897862233},
				{"Campus &amp; Community – UC Santa Cruz", "", "A place for all to thrive", "https://www.ucsc.edu/about/", 0.03230337078651685},
			},
		},
	}

	url := "https://www.ucsc.edu/wp-sitemap.xml"

	crawl(url, true, true)

	for _, tc := range tests {
		terms := []string{"thrive"}
		got := search(terms, false, false)

		for _, gotHit := range got {
			for _, wantHit := range tc.want {
				if gotHit.Url == wantHit.Url {
					if !reflect.DeepEqual(gotHit.Sentence, wantHit.Sentence) {
						t.Errorf("Want title: %v\nGot title: %v\n", wantHit.Sentence, gotHit.Sentence)
					}
				}
			}
		}
	}
}
