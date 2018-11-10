package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
)

type cmdLineBrowser struct {
	application string
	params      []string
	url         string
}

func renderPage(URL string, file bool) []byte {
	chrome := cmdLineBrowser{
		application: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		params:      []string{"--headless", "--disable-gpu", "--dump-dom"},
		url:         URL}
	cmd := exec.Command(chrome.application, append(chrome.params, chrome.url)...)
	log.Printf("Running Chrome to render the page for %s ...", chrome.url)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Chrome finished with error: %v", err)
	} else {
		if file {
			u, errP := url.Parse(URL)
			if errP != nil {
				log.Printf("Error parsing URL: %v", errP)
			}
			log.Println("Writing file", u.Hostname()+"_out.html")
			ioutil.WriteFile(u.Hostname()+"_out.html", []byte(out), 0666)
		}
	}
	return out
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	content := renderPage(url, false)
	fmt.Fprintf(w, "{'URL':'%s'}{'content':'%s'}", url, content)
}

func testPages() {
	pages := []string{"http://index.hu/", "https://www.theguardian.com/uk"}
	for _, p := range pages {
		renderPage(p, true)
	}
}

func main() {
	testPages()
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
