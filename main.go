package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	backdoorFile    = "lib/shell.json"
	userAgentsFile  = "lib/useragent.json"
	workerCount     = 50 // Jumlah maksimal worker berjalan bersamaan
	timeout         = 10 * time.Second
)

var httpClient = &http.Client{
	Timeout: timeout,
}

type Scanner struct {
	Targets    []string
	Paths      []string
	UserAgents []string
	Results    chan string
	Mutex      sync.Mutex
}

type UserAgentList struct {
	UserAgents []string `json:"user_agents"`
}

func printBanner() {
	banner := `
	███╗   ███╗ █████╗ ██╗    ██╗██████╗ ███████╗██╗██╗  ██╗
	████╗ ████║██╔══██╗██║    ██║╚════██╗██╔════╝██║╚██╗██╔╝
	██╔████╔██║███████║██║ █╗ ██║ █████╔╝███████╗██║ ╚███╔╝ 
	██║╚██╔╝██║██╔══██║██║███╗██║ ╚═══██╗╚════██║██║ ██╔██╗ 
	██║ ╚═╝ ██║██║  ██║╚███╔███╔╝██████╔╝███████║██║██╔╝ ██╗
	╚═╝     ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚═════╝ ╚══════╝╚═╝╚═╝  ╚═╝
	Webshell Finder Base On Go!    Hello Friend? @maw3six                                                 `
	fmt.Println(banner)
}

func (s *Scanner) LoadTargets(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		s.Targets = append(s.Targets, url)
	}
	return scanner.Err()
}

func (s *Scanner) LoadBackdoorPaths() error {
	file, err := os.Open(backdoorFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&s.Paths)
}

func (s *Scanner) LoadUserAgents() error {
	file, err := os.Open(userAgentsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var uaList UserAgentList
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&uaList); err != nil {
		return err
	}

	s.UserAgents = uaList.UserAgents
	return nil
}

func (s *Scanner) GetRandomUserAgent() string {
	if len(s.UserAgents) == 0 {
		return "Mozilla/5.0"
	}
	return s.UserAgents[rand.Intn(len(s.UserAgents))]
}

func (s *Scanner) CheckBackdoor(target string, path string, results chan<- string) {
	fullURL := target + path
	fmt.Printf("\033[33m[SCANNING]\033[0m %s\n", fullURL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Printf("\033[31m[ERROR]\033[0m Error creating request: %s\n", err)
		return
	}
	req.Header.Set("User-Agent", s.GetRandomUserAgent())
	req.Header.Set("Referer", target)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("\033[31m[FAILED]\033[0m %s - %s\n", fullURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("\033[32m[FOUND]\033[0m %s\n", fullURL)
		results <- fullURL
	} else {
		fmt.Printf("\033[31m[NOT FOUND]\033[0m %s (Status: %d)\n", fullURL, resp.StatusCode)
	}
}

func (s *Scanner) WorkerPool() {
	jobs := make(chan struct {
		target string
		path   string
	}, len(s.Targets)*len(s.Paths))
	results := make(chan string, len(s.Targets)*len(s.Paths))

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				s.CheckBackdoor(job.target, job.path, results)
			}
		}()
	}

	for _, target := range s.Targets {
		for _, path := range s.Paths {
			jobs <- struct {
				target string
				path   string
			}{target, path}
		}
	}
	close(jobs)
	wg.Wait()
	close(results)

	for result := range results {
		s.SaveResult(result)
	}
}

func (s *Scanner) SaveResult(result string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	file, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error saving result:", err)
		return
	}
	defer file.Close()
	file.WriteString(result + "\n")
}

func main() {
	printBanner()
	fmt.Println("========================================================================")
	fmt.Print("List Target: ")
	var targetFile string
	fmt.Scanln(&targetFile)

	scanner := &Scanner{}

	if err := scanner.LoadTargets(targetFile); err != nil {
		fmt.Println("Error loading targets:", err)
		return
	}

	if err := scanner.LoadBackdoorPaths(); err != nil {
		fmt.Println("Error loading backdoor paths from", backdoorFile, ":", err)
		return
	}

	if err := scanner.LoadUserAgents(); err != nil {
		fmt.Println("Error loading user agents from", userAgentsFile, ":", err)
		return
	}

	fmt.Println("Starting scan with", workerCount, "workers...")
	scanner.WorkerPool()
	fmt.Println("Scan selesai. Hasil disimpan di results.txt")
}
