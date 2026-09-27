package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CookieData struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Path   string `json:"path"`
	Domain string `json:"domain"`
}

type SessionData struct {
	StudentID string       `json:"student_id"`
	Cookies   []CookieData `json:"cookies"`
}

type VideoEntry struct {
	URL        string `json:"url"`
	Filename   string `json:"filename"`
	Downloaded bool   `json:"downloaded"`
	FoundAt    string `json:"found_at"`
}

// WriteCounter counts the number of bytes written to it and prints progress.
type WriteCounter struct {
	Total uint64
	Size  uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	if wc.Size > 0 {
		pct := float64(wc.Total) / float64(wc.Size) * 100
		fmt.Printf("\rDownloading... %s / %s (%.2f%%)", formatBytes(wc.Total), formatBytes(wc.Size), pct)
	} else {
		fmt.Printf("\rDownloading... %s", formatBytes(wc.Total))
	}
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	mode := os.Args[1]

	switch mode {
	case "find":
		fs := flag.NewFlagSet("find", flag.ExitOnError)
		username := fs.String("u", "", "Username")
		password := fs.String("p", "", "Password")
		maxID := fs.Int("max-id", 150, "Maximum classview ID to scan")
		_ = fs.Parse(os.Args[2:])
		runFindMode(username, password, *maxID)
	case "download":
		runDownloadMode()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Shaolin Course Video Scraper")
	fmt.Println("\nUsage:")
	fmt.Println("  shaolin-scraper <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  find        Logs in and crawls the portal to index all video resources")
	fmt.Println("  download    Downloads indexed videos (all or one-at-a-time)")
	fmt.Println("\nOptions for 'find':")
	fmt.Println("  -u <string>   Login username (optional, will prompt if omitted)")
	fmt.Println("  -p <string>   Login password (optional, will prompt if omitted)")
	fmt.Println("  -max-id <int> Maximum classview ID to scan (default 150)")
	fmt.Println("\nExamples:")
	fmt.Println("  ./shaolin-scraper find -u user@example.com -p secret")
	fmt.Println("  ./shaolin-scraper find -max-id 200")
	fmt.Println("  ./shaolin-scraper download")
}

func runFindMode(userOpt *string, passOpt *string, maxID int) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	var studentID string
	var sessionLoaded bool

	session, err := loadSession(jar, "session.json")
	if err == nil && session != nil {
		fmt.Println("[*] Loaded existing session from session.json")
		if isSessionValid(client, session.StudentID) {
			fmt.Println("[+] Session is valid! Resuming without login.")
			studentID = session.StudentID
			sessionLoaded = true
		} else {
			fmt.Println("[-] Session has expired or is invalid. Login required.")
		}
	}

	if !sessionLoaded {
		username := ""
		password := ""
		if userOpt != nil {
			username = *userOpt
		}
		if passOpt != nil {
			password = *passOpt
		}

		if username == "" || password == "" {
			fmt.Println("[*] Credentials not provided via flags. Please enter them interactively:")
			var err error
			username, password, err = readCredentialsInteractive()
			if err != nil {
				fmt.Printf("[-] Error reading credentials: %v\n", err)
				return
			}
		}

		fmt.Println("[*] Attempting login to learnshaolinonline.com...")
		id, err := attemptLogin(client, username, password)
		if err != nil {
			fmt.Printf("[-] Login failed: %v\n", err)
			return
		}
		studentID = id
		fmt.Printf("[+] Login successful! Student ID: %s\n", studentID)

		err = saveSession(studentID, client.Jar, "session.json")
		if err != nil {
			fmt.Printf("[-] Error saving session: %v\n", err)
		} else {
			fmt.Println("[*] Session saved to session.json")
		}
	}

	// Generate target seed URLs: Student page + all classview/## URLs
	var startURLs []string
	startURLs = append(startURLs, fmt.Sprintf("https://learnshaolinonline.com/index.php/student/%s", studentID))

	fmt.Printf("[*] Seeding crawler with student dashboard and classview IDs 1 to %d...\n", maxID)
	for i := 1; i <= maxID; i++ {
		startURLs = append(startURLs, fmt.Sprintf("https://learnshaolinonline.com/index.php/course/classview/%d", i))
	}

	discovered, err := crawl(startURLs, client)
	if err != nil {
		fmt.Printf("[-] Error during crawl: %v\n", err)
		return
	}

	existing, err := loadIndex("video_index.json")
	if err != nil {
		fmt.Printf("[-] Error loading existing index, starting fresh: %v\n", err)
		existing = nil
	}

	merged := mergeVideos(existing, discovered)

	err = saveIndex(merged, "video_index.json")
	if err != nil {
		fmt.Printf("[-] Error saving index: %v\n", err)
	} else {
		fmt.Println("[+] Video index successfully saved to video_index.json")
	}
}

func runDownloadMode() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	session, err := loadSession(jar, "session.json")
	if err != nil {
		fmt.Println("[*] No saved session found. Attempting unauthenticated downloads...")
	} else {
		fmt.Printf("[*] Loaded saved session for Student: %s. Using it for downloads.\n", session.StudentID)
	}

	videos, err := loadIndex("video_index.json")
	if err != nil {
		fmt.Printf("[-] Error loading index: %v\n", err)
		return
	}

	if len(videos) == 0 {
		fmt.Println("[-] The video index (video_index.json) is empty. Please run the 'find' command first.")
		return
	}

	// Filter pending videos and check local disk
	var pending []*VideoEntry
	for i := range videos {
		if !videos[i].Downloaded {
			if info, err := os.Stat(videos[i].Filename); err == nil && info.Size() > 0 {
				videos[i].Downloaded = true
				continue
			}
			pending = append(pending, &videos[i])
		}
	}

	// Sync index file
	_ = saveIndex(videos, "video_index.json")

	if len(pending) == 0 {
		fmt.Println("[+] All videos in the index have already been successfully downloaded!")
		return
	}

	fmt.Printf("[*] Discovered %d pending videos:\n", len(pending))
	for i, v := range pending {
		fmt.Printf("  [%d] %s (Found at: %s)\n", i+1, v.Filename, v.FoundAt)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nOptions:")
		fmt.Println("  [a] Download ALL pending videos")
		fmt.Println("  [o] Download ONE-at-a-time (interactive)")
		fmt.Println("  [1-N] Download a specific video by its number")
		fmt.Println("  [q] Quit")
		fmt.Print("Choose option: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[-] Error reading input: %v\n", err)
			return
		}
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "q" {
			fmt.Println("Exiting download mode.")
			return
		}

		if input == "a" {
			fmt.Printf("[*] Starting download of all %d videos...\n", len(pending))
			for _, v := range pending {
				err := downloadVideo(client, v)
				if err != nil {
					fmt.Printf("[-] Error downloading %s: %v\n", v.Filename, err)
					_ = saveIndex(videos, "video_index.json")
					continue
				}
				_ = saveIndex(videos, "video_index.json")
			}
			fmt.Println("\n[+] Finished downloading all pending files.")
			return
		}

		if input == "o" {
			fmt.Println("[*] Starting one-at-a-time download...")
			for idx, v := range pending {
				err := downloadVideo(client, v)
				if err != nil {
					fmt.Printf("[-] Error downloading %s: %v\n", v.Filename, err)
					_ = saveIndex(videos, "video_index.json")
					continue
				}
				_ = saveIndex(videos, "video_index.json")

				if idx == len(pending)-1 {
					fmt.Println("\n[+] Downloaded the last pending video.")
					break
				}

				fmt.Print("\nDownload next pending video? [Y/n]: ")
				nextInput, _ := reader.ReadString('\n')
				nextInput = strings.TrimSpace(strings.ToLower(nextInput))
				if nextInput == "n" || nextInput == "no" {
					fmt.Println("Stopping one-at-a-time downloads.")
					break
				}
			}
			return
		}

		var num int
		_, err = fmt.Sscanf(input, "%d", &num)
		if err == nil && num >= 1 && num <= len(pending) {
			v := pending[num-1]
			err := downloadVideo(client, v)
			if err != nil {
				fmt.Printf("[-] Error downloading %s: %v\n", v.Filename, err)
			}
			_ = saveIndex(videos, "video_index.json")
			return
		}

		fmt.Println("[-] Invalid option, please try again.")
	}
}

func attemptLogin(client *http.Client, username, password string) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField("loginName")
	if err != nil {
		return "", err
	}
	if _, err = fw.Write([]byte(username)); err != nil {
		return "", err
	}

	fw, err = w.CreateFormField("loginPassword")
	if err != nil {
		return "", err
	}
	if _, err = fw.Write([]byte(password)); err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest("POST", "https://learnshaolinonline.com/index.php/sec/attemptlogin", &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	responseStr := string(bodyBytes)
	if !strings.Contains(responseStr, "Success") {
		return "", fmt.Errorf("server responded with: %s", responseStr)
	}

	parts := strings.Split(responseStr, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid success format: %s", responseStr)
	}

	return strings.TrimSpace(parts[1]), nil
}

func isSessionValid(client *http.Client, studentID string) bool {
	u := fmt.Sprintf("https://learnshaolinonline.com/index.php/student/%s", studentID)
	resp, err := client.Get(u)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if strings.Contains(resp.Request.URL.Path, "/sec/login") {
		return false
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	bodyStr := string(bodyBytes)
	if strings.Contains(bodyStr, "SignInButton") || strings.Contains(bodyStr, "LoginPassword") {
		return false
	}

	return true
}

func crawl(startURLs []string, client *http.Client) ([]VideoEntry, error) {
	var videos []VideoEntry
	visited := make(map[string]bool)
	queue := append([]string{}, startURLs...)

	fmt.Println("[*] Starting crawl...")

	maxPages := 500
	pagesScanned := 0

	for len(queue) > 0 && pagesScanned < maxPages {
		uStr := queue[0]
		queue = queue[1:]

		if visited[uStr] {
			continue
		}
		visited[uStr] = true

		u, err := url.Parse(uStr)
		if err != nil {
			continue
		}

		isClassView := strings.Contains(u.Path, "/course/classview/")

		req, err := http.NewRequest("GET", uStr, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			if !isClassView {
				fmt.Printf("  [-] Error fetching %s: %v\n", uStr, err)
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		pagesScanned++
		if isClassView {
			fmt.Printf("[%d] Scanning active classview page: %s\n", pagesScanned, uStr)
		} else {
			fmt.Printf("[%d] Scanning: %s\n", pagesScanned, uStr)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		htmlContent := string(bodyBytes)
		pageLinks, videoURLs := extractLinks(htmlContent, u)

		// Process videos
		for _, vURL := range videoURLs {
			duplicate := false
			for _, existing := range videos {
				if existing.URL == vURL {
					duplicate = true
					break
				}
			}
			if !duplicate {
				filename := filepath.Base(vURL)
				videos = append(videos, VideoEntry{
					URL:        vURL,
					Filename:   filename,
					Downloaded: false,
					FoundAt:    uStr,
				})
				fmt.Printf("  [+] Discovered video: %s\n", filename)
			}
		}

		// Process page links
		for _, pLink := range pageLinks {
			if !visited[pLink] && shouldCrawlString(pLink, u.Host) {
				queue = append(queue, pLink)
			}
		}
	}

	fmt.Printf("[*] Crawl finished. Visited %d active pages. Found %d videos.\n", pagesScanned, len(videos))
	return videos, nil
}

func shouldCrawlString(uStr string, targetHost string) bool {
	u, err := url.Parse(uStr)
	if err != nil {
		return false
	}
	path := strings.ToLower(u.Path)
	if strings.Contains(path, "/logout") || strings.Contains(path, "/login") || strings.Contains(path, "/sec/attemptlogin") {
		return false
	}
	// Support both subdomain and main domain
	if u.Host != targetHost && u.Host != "www."+targetHost && "www."+u.Host != targetHost {
		return false
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".css", ".js", ".ico", ".svg", ".woff", ".woff2", ".pdf", ".zip":
		return false
	}
	return true
}

func extractLinks(htmlContent string, baseURL *url.URL) ([]string, []string) {
	var pageLinks []string
	var videoURLs []string

	attrRegex := regexp.MustCompile(`(?i)(?:href|src)\s*=\s*["']([^"']+)["']`)
	matches := attrRegex.FindAllStringSubmatch(htmlContent, -1)
	for _, match := range matches {
		val := match[1]
		resolved, err := resolveURL(val, baseURL)
		if err != nil {
			continue
		}
		if strings.HasSuffix(strings.ToLower(resolved.Path), ".mp4") {
			videoURLs = append(videoURLs, resolved.String())
		} else {
			pageLinks = append(pageLinks, resolved.String())
		}
	}

	jsRegexes := []*regexp.Regexp{
		regexp.MustCompile(`(?i)location\s*=\s*["']([^"']+)["']`),
		regexp.MustCompile(`(?i)window\.open\(\s*["']([^"']+)["']`),
	}
	for _, r := range jsRegexes {
		jsMatches := r.FindAllStringSubmatch(htmlContent, -1)
		for _, match := range jsMatches {
			val := match[1]
			resolved, err := resolveURL(val, baseURL)
			if err != nil {
				continue
			}
			if strings.HasSuffix(strings.ToLower(resolved.Path), ".mp4") {
				videoURLs = append(videoURLs, resolved.String())
			} else {
				pageLinks = append(pageLinks, resolved.String())
			}
		}
	}

	return pageLinks, videoURLs
}

func resolveURL(ref string, baseURL *url.URL) (*url.URL, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}
	return baseURL.ResolveReference(u), nil
}

func downloadVideo(client *http.Client, video *VideoEntry) error {
	fmt.Printf("\n[*] Starting download for: %s\n", video.Filename)

	req, err := http.NewRequest("GET", video.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d %s", resp.StatusCode, resp.Status)
	}

	size := uint64(resp.ContentLength)

	tempFilename := video.Filename + ".tmp"
	out, err := os.Create(tempFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	counter := &WriteCounter{Size: size}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		os.Remove(tempFilename)
		return err
	}

	out.Close()

	err = os.Rename(tempFilename, video.Filename)
	if err != nil {
		return err
	}

	video.Downloaded = true
	fmt.Printf("\n[+] Successfully downloaded: %s\n", video.Filename)
	return nil
}

func saveIndex(videos []VideoEntry, filepath string) error {
	for i := range videos {
		if !videos[i].Downloaded {
			if info, err := os.Stat(videos[i].Filename); err == nil && info.Size() > 0 {
				videos[i].Downloaded = true
				fmt.Printf("[*] Sync: Detected local file for %s, marking as downloaded in index.\n", videos[i].Filename)
			}
		}
	}

	data, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

func loadIndex(filepath string) ([]VideoEntry, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, nil
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var videos []VideoEntry
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func readCredentialsInteractive() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	fmt.Print("Enter Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	password = strings.TrimSpace(password)

	return username, password, nil
}

func saveSession(studentID string, jar http.CookieJar, filepath string) error {
	u, _ := url.Parse("https://learnshaolinonline.com")
	cookies := jar.Cookies(u)

	var cookieData []CookieData
	for _, c := range cookies {
		cookieData = append(cookieData, CookieData{
			Name:   c.Name,
			Value:  c.Value,
			Path:   c.Path,
			Domain: c.Domain,
		})
	}

	session := SessionData{
		StudentID: studentID,
		Cookies:   cookieData,
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0600)
}

func loadSession(jar http.CookieJar, filepath string) (*SessionData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	for _, c := range session.Cookies {
		cookies = append(cookies, &http.Cookie{
			Name:   c.Name,
			Value:  c.Value,
			Path:   c.Path,
			Domain: c.Domain,
		})
	}

	u, _ := url.Parse("https://learnshaolinonline.com")
	jar.SetCookies(u, cookies)

	return &session, nil
}

func mergeVideos(existing []VideoEntry, discovered []VideoEntry) []VideoEntry {
	videoMap := make(map[string]VideoEntry)
	for _, v := range existing {
		videoMap[v.URL] = v
	}

	for _, v := range discovered {
		if val, exists := videoMap[v.URL]; exists {
			v.Downloaded = val.Downloaded
		}
		videoMap[v.URL] = v
	}

	var result []VideoEntry
	for _, v := range videoMap {
		result = append(result, v)
	}
	return result
}
