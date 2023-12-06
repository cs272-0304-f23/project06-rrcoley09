package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/kljensen/snowball"
	"github.com/neurosnap/sentences/english"
	"golang.org/x/net/html"
)

type ExtractResult struct {
	words, sentences  []string
	url, webPageTitle string
	images            map[string][]string
}

func renderStr(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)

	title := strings.Trim(buf.String(), "<title>")
	title = strings.TrimSuffix(title, "</")

	return title
}

func extract(inCh chan string, dOutCh chan DownloadResult, exOutCh chan ExtractResult, wg *sync.WaitGroup, n *int32, siteMap bool) {
	// Loop over downloadResults for each url in in channel (inCh) to download
	for downloadResult := range dOutCh {
		io_reader := bytes.NewReader(downloadResult.body)
		doc, err := html.Parse(io_reader)
		if err != nil {
			fmt.Printf("error parsing HTML data: %v\n", err.Error())
			return
		}

		// Initialize string slices and title variable
		var hrefsSlice, wordsSlice, captionWordSlice, sliceSentence []string
		imagesWCaptions := make(map[string][]string)
		var titleExtracted string

		// Obtain map of stop words, calling func stopWords()
		mapStopWords := stopWords()

		// Closure function to extract words and hrefs
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode {
				// Go to anchor attributes
				if n.Data == "a" {
					// Loop over anchor attributes looking for hrefs
					for _, value := range n.Attr {
						// If href found append to slice of hrefs
						if value.Key == "href" {
							hrefsSlice = append(hrefsSlice, value.Val)
						}
					}
				} else if n.Data == "style" || n.Data == "script" { // Skip style and script nodes
					return
				} else if n.Data == "title" { // Extract and trim webpage title
					titleExtracted = renderStr(n)
				} else if n.Data == "img" {
					var src, alt string
					for _, value := range n.Attr {
						switch value.Key {
						case "src":
							src = value.Val
						case "alt":
							alt = value.Val
						}
					}
					f := func(c rune) bool {
						return !unicode.IsLetter(c) && !unicode.IsNumber(c)
					}

					// Strip bunched words
					captionWords := strings.Fields(alt)

					// Add stripped words to slice to return
					for _, word := range captionWords {
						captionSplit := strings.FieldsFunc(word, f)

						for _, splitWord := range captionSplit {
							stemmed, err := snowball.Stem(splitWord, "english", true)

							if err != nil {
								fmt.Println("error stemming word for caption")
							}

							if _, exists := mapStopWords[stemmed]; !exists {
								captionWordSlice = append(captionWordSlice, stemmed)
							}
						}
					}
					if src != "" || alt != "" {
						imagesWCaptions[src] = captionWordSlice
					}
				}
			} else if n.Type == html.TextNode {
				tokenizer, err := english.NewSentenceTokenizer(nil)
				if err != nil {
					panic(err)
				}

				sentences := tokenizer.Tokenize(n.Data)
				for _, s := range sentences {
					if trimmedSentence := strings.TrimSpace(s.Text); trimmedSentence != "" {
						sliceSentence = append(sliceSentence, trimmedSentence)
					}
				}

				// Strip punctuation from words
				f := func(c rune) bool {
					return !unicode.IsLetter(c) && !unicode.IsNumber(c)
				}

				// Strip bunched words
				wordsNoSpace := strings.Fields(n.Data)

				// Add stripped words to slice to return
				for _, word := range wordsNoSpace {
					wordsSplit := strings.FieldsFunc(word, f)
					wordsSlice = append(wordsSlice, wordsSplit...)
				}

			}
			for c := n.FirstChild; c != nil; c = c.NextSibling { // Recursive call to loop over all data
				f(c)
			}
		}
		// Call function on body
		f(doc)

		postStopStemSlice := []string{}

		// Stem and filter out stop words from wordsSlice
		for _, word := range wordsSlice {
			stemmed, err := snowball.Stem(word, "english", true)
			if err != nil {
				fmt.Println("error stemming word")
			}

			if _, exists := mapStopWords[stemmed]; !exists {
				postStopStemSlice = append(postStopStemSlice, stemmed)
			}
		}

		if !siteMap {
			// Loop of hrefs in seed url
			for _, href := range hrefsSlice {
				cleanURL, err := clean(downloadResult.url, href)

				// Check if url has already been crawled
				if err == nil {
					// Increment waitGroup value by 1 for each new url
					wg.Add(1)
					atomic.AddInt32(n, 1)

					// Insert subsequent clean urls into inCh to be downloaded
					inCh <- cleanURL
				}
			}
		} else {
			// Extract urls from sitemap
			extractSitemapUrl(downloadResult.body, n, wg, inCh)
		}
		// Return extract result for use in add_index, which populates index
		exOutCh <- ExtractResult{postStopStemSlice, sliceSentence, downloadResult.url, titleExtracted, imagesWCaptions}
	}
}
