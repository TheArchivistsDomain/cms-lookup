package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/*
Drupal: / ---> '/sites/default/'
OpenCart: / ---> 'catalog/view/'
PrestaShop: / ---> 'content="PrestaShop"'
vBulletin: / ---> 'window.vBulletin' Laravel: /
(cookies) ---> 'laravel_session'
Wordpress: /wp-includes/css/buttons.css ---> 'WordPress-style Buttons'
Joomla: /media/system/js/core.js ---> 'window.Joomla' */

func check_url(input_url string) string {
    client := http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get(input_url)
    if err == nil {
        defer resp.Body.Close()
                body := new(strings.Builder)
        io.Copy(body, resp.Body)
        if strings.Contains(body.String(), "/sites/default"){
                        return "Drupal"
                }
                if strings.Contains(body.String(), "catalog/view/") {
            return "OpenCart"
                }
                if strings.Contains(body.String(), "content=\"PrestaShop\"") {
                        return "PrestaShop"
                }
                if strings.Contains(body.String(), "window.vBulletin") {
            return "vBulletin"
                }
                for _, cookie := range resp.Cookies() { if cookie.Name == "laravel_session" {
                                return "Laravel"
                        }
                }

                resp, err := client.Get(input_url + "wp-includes/css/buttons.css")
                if err == nil {
                        defer resp.Body.Close()
            body := new(strings.Builder)
                        io.Copy(body, resp.Body)

                        if strings.Contains(body.String(), "WordPress-style Buttons"){
                                return "Wordpress"
                        }
                }
                resp, err = client.Get(input_url + "media/system/js/core.js")
        if
                err == nil {
                        defer resp.Body.Close()
            body := new(strings.Builder)
                        io.Copy(body, resp.Body)
            if strings.Contains(body.String(), "window.Joomla"){
                                return "Joomla"
                        }
                }
                return "Other"
        }
        return "Invalid"
}
func check_url_chunk(chunk chan string) {
    for url := range chunk {
        if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
            url = "http://" + url
        }
    
        if !strings.HasSuffix(url, "/") {
            url = url + "/"
        }

        result := check_url(url)
        if result != "Invalid" {
            append_to_file("cmslookup_" + result + ".txt", url + "\n")
            fmt.Printf("[✓] %v: %v\n", url, result)
        }
    }
}
func read_urls(file_name string) []string {
    file, err := os.Open(file_name)
    if err != nil {
                log.Fatal(err)
        }
        defer file.Close()
    sc := bufio.NewScanner(file)
    lines := make([]string, 0)
        // Read through 'urls' until an EOF is encountered.
        for sc.Scan() {
        lines = append(lines, sc.Text())
        }
        if err := sc.Err(); err != nil {
        log.Fatal(err)
        }
        return lines
}

func append_to_file(filename string, text string){
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
    }
    defer f.Close()
    if _, err := f.WriteString(text); err != nil {
        log.Println(err)
    }
}

func main() {
    fmt.Println("[!] CMSLookup by The Archivist")
    fmt.Println("[!] Contact @The_Archivist_01")
        // define the command line arguments for the utility.
        arg_filename := flag.String("filename", "", "File containing the list of URL's to process.")
    arg_threads := flag.Int("threads", 5, "Number of threads to run.")
        // parse the given arguments and check for mandatory parameters
        flag.Parse()
    if *arg_filename == "" {
        flag.Usage()
        os.Exit(1)
        }
        // Attempt to read the specified file
        urls := read_urls(*arg_filename)
        // Create a channel which will take one string
        urls_channel := make(chan string)
        // Start populating this channel with the urls in a go routine, so that
        // the rest of the application can continue.
        go func() {
        for _, url := range urls {
            urls_channel <- url
                }
                // Once all the urls have been put onto the channel, close the
                // channel so that the checker routines can finish their for-loop
                // and stop.
                close(urls_channel)
        }()
        // Create a waitgroup to use in the scanning go routines. This will allow
        // us to wait and close down the results channel later on. This will have
        // the effect of letting the main thread which is looping over the results
        // channel exit gracefully.
        var wg sync.WaitGroup
    fmt.Printf("[!] Starting %d threads\n", *arg_threads)
    for i := 1; i <= *arg_threads; i++ {
                wg.Add(1)
        go func() {
            defer wg.Done()
                        check_url_chunk(urls_channel)
                }()
        }
        // meanwhile, in the main thread we simply wait for the waitgroup to be
        // empty. the url list reader will close the urls channel when its done,
        // meaning the go routines will also close when there are no more items on
        // the urls channel. each of the wg.Done() calls will be made when the go
        // routines exit, and the wait here will unblock. if a result is found,
        // the call to cancel() will stop the go routines, unblocking us here as
        // well.
        wg.Wait()
    fmt.Println("[!] Finished checking...")
}
