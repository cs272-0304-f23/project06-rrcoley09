package main

import (
	"reflect"
	"testing"
)

func TestNGram(t *testing.T) {
	tests := []struct {
		term string
		want []Hits
	}{
		{"uc santa",
			[]Hits{
				{"https://www.ucsc.edu/about/", "About UC Santa Cruz – UC Santa Cruz", 0.02896769662921348},
				{"https://www.ucsc.edu/about/overview/", "Overview – UC Santa Cruz", 0.026902173913043476},
				{"https://www.ucsc.edu/programs-and-units/", "Programs and Units – UC Santa Cruz", 0.026258680555555552},
				{"https://www.ucsc.edu/campus/", "Campus & Community – UC Santa Cruz", 0.02449524940617577},
				{"https://www.ucsc.edu/research/", "Research – UC Santa Cruz", 0.023174157303370788},
				{"https://www.ucsc.edu/campus/mascot/", "Our Mascot: Sammy the Banana Slug – UC Santa Cruz", 0.022664835164835168},
				{"https://www.ucsc.edu/better-together/", "UC + Santa Cruz. Better together – UC Santa Cruz", 0.02142857142857143},
				{"https://www.ucsc.edu/", "UC Santa Cruz – A world-class public research institution comprised of ten residential college communities nestled in the redwood forests and meadows overlooking central California's Monterey Bay.", 0.0207286432160804},
				{"https://www.ucsc.edu/land-acknowledgment/", "Land Acknowledgment – UC Santa Cruz", 0.019908301158301157},
				{"https://www.ucsc.edu/principles-community/", "Principles of Community – UC Santa Cruz", 0.01939655172413793},
				{"https://www.ucsc.edu/admissions/", "Admissions & Aid – UC Santa Cruz", 0.018284574468085107},
				{"https://www.ucsc.edu/campus/visit/", "Visit UCSC – UC Santa Cruz", 0.01657958199356913},
				{"https://www.ucsc.edu/about/achievements-facts-and-figures/", "Achievements, Facts, and Figures – UC Santa Cruz", 0.015345982142857142},
				{"https://www.ucsc.edu/mission-and-vision/", "Mission and Vision – UC Santa Cruz", 0.015054744525547444},
				{"https://www.ucsc.edu/residential-colleges/", "Residential Colleges – UC Santa Cruz", 0.014732142857142857},
				{"https://www.ucsc.edu/campus-destinations/", "Campus Destinations – UC Santa Cruz", 0.013323643410852713},
				{"https://www.ucsc.edu/people/", "Find people – UC Santa Cruz", 0.012890625000000001},
				{"https://www.ucsc.edu/about/leadership/", "Leading the change – UC Santa Cruz", 0.011752136752136752},
				{"https://www.ucsc.edu/research/undergraduate-research/", "Undergraduate Research – UC Santa Cruz", 0.010522959183673469},
				{"https://www.ucsc.edu/feedback/", "Feedback – UC Santa Cruz", 0.0103125},
				{"https://www.ucsc.edu/search/", "Search results – UC Santa Cruz", 0.0103125},
				{"https://www.ucsc.edu/author/milpowelucsc-edu/", "Miranda Powell – UC Santa Cruz", 0.01021039603960396},
				{"https://www.ucsc.edu/author/raknightucsc-edu/", "raknight@ucsc.edu – UC Santa Cruz", 0.01021039603960396},
				{"https://www.ucsc.edu/author/lmnielseucsc-edu/", "Lisa Nielsen – UC Santa Cruz", 0.01021039603960396},
				{"https://www.ucsc.edu/author/gwenjucsc-edu/", "Gwen Jourdonnais – UC Santa Cruz", 0.01021039603960396},
				{"https://www.ucsc.edu/campus/campus-galleries-and-theaters/", "Campus Galleries and Theaters – UC Santa Cruz", 0.009915865384615386},
				{"https://www.ucsc.edu/campus/visit/maps-directions/", "Campus maps and directions – UC Santa Cruz", 0.008152173913043478},
				{"https://www.ucsc.edu/calendars/", "Calendars – UC Santa Cruz", 0.0074190647482014396},
				{"https://www.ucsc.edu/academics/", "Academics – UC Santa Cruz", 0.004209183673469388},
				{"https://www.ucsc.edu/address-and-phone/", "Address and phone – UC Santa Cruz", 0.003763686131386861},
				{"https://www.ucsc.edu/azindex/", "A–Z index – UC Santa Cruz", 0.003020650263620387},
				{"https://www.ucsc.edu/privacy-policy/", "Privacy Policy – UC Santa Cruz", 0.0027871621621621623},
			},
		},
	}

	url := "https://www.ucsc.edu/wp-sitemap.xml"

	go serve_web()
	crawl(url, true, false)

	for _, tc := range tests {
		terms := []string{"uc", "santa"}
		got := search(terms, false)

		for _, gotHit := range got {
			//fmt.Println(gotHit)
			for _, wantHit := range tc.want {
				if gotHit.url == wantHit.url {
					if !reflect.DeepEqual(gotHit.title, wantHit.title) {
						t.Errorf("Want title: %v\nGot title: %v\n", wantHit.title, gotHit.title)
					}
				}
			}
		}
	}
}

func TestWildCard(t *testing.T) {
	tests := []struct {
		term string
		want []Hits
	}{
		{"farm",
			[]Hits{
				{"https://www.ucsc.edu/campus-destinations/", "Campus Destinations – UC Santa Cruz", 0.036544850498338874},
				{"https://www.ucsc.edu/campus/visit/", "Visit UCSC – UC Santa Cruz", 0.01515847496554892},
				{"https://www.ucsc.edu/about/overview/", "Overview – UC Santa Cruz", 0.013664596273291927},
				{"https://www.ucsc.edu/about/", "About UC Santa Cruz – UC Santa Cruz", 0.013242375601926164},
				{"https://www.ucsc.edu/better-together/", "UC + Santa Cruz. Better together – UC Santa Cruz", 0.006122448979591837},
				{"https://www.ucsc.edu/azindex/", "A–Z index – UC Santa Cruz", 0.0027617373838814963},
			},
		},
	}

	url := "https://www.ucsc.edu/wp-sitemap.xml"

	crawl(url, true, true)

	for _, tc := range tests {
		terms := []string{"farm"}
		got := search(terms, true)

		for _, gotHit := range got {
			for _, wantHit := range tc.want {
				if gotHit.url == wantHit.url {
					if !reflect.DeepEqual(gotHit.title, wantHit.title) {
						t.Errorf("Want title: %v\nGot title: %v\n", wantHit.title, gotHit.title)
					}
				}
			}
		}
	}
}
