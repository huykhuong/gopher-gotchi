package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GistUpdate struct {
	Files map[string]map[string]string `json:"files"`
}

func getCloudCredentials() (string, string, string) {
	token := os.Getenv("DIANA_GITHUB_TOKEN")
	gistID := os.Getenv("DIANA_GIST_ID")
	logGistID := os.Getenv("DIANA_GIST_LOG_ID")

	if token == "" || gistID == "" {
		return "", "", ""
	}

	return token, gistID, logGistID
}

func patchGist(token, gistID string, files map[string]map[string]string) error {
	payload := GistUpdate{Files: files}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", "https://api.github.com/gists/"+gistID, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gist API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func fetchExistingMissionLog(token, gistID string) string {
	req, _ := http.NewRequest("GET", "https://api.github.com/gists/"+gistID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var gist struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}

	json.NewDecoder(resp.Body).Decode(&gist)

	if f, ok := gist.Files["mission_log.txt"]; ok {
		return f.Content
	}

	return ""
}

func (p *Pet) SyncAllToCloud(newMemory *Memory) error {
	token, gistID, logGistID := getCloudCredentials()

	if token == "" {
		return fmt.Errorf("missing cloud credentials")
	}

	data, _ := json.MarshalIndent(p, "", "  ")

	if err := patchGist(token, gistID, map[string]map[string]string{
		"diana.json": {"content": string(data)},
	}); err != nil {
		return err
	}

	if newMemory != nil {
		existing := fetchExistingMissionLog(token, logGistID)
		logEntry := fmt.Sprintf("[%s] Level %d: %s\n", newMemory.Timestamp, newMemory.Level, newMemory.Message)
		if err := patchGist(token, logGistID, map[string]map[string]string{
			"mission_log.txt": {"content": existing + logEntry},
		}); err != nil {
			return err
		}
	}

	return nil
}

func LoadFromCloud() (*Pet, error) {
	token, gistID, _ := getCloudCredentials()

	req, _ := http.NewRequest("GET", "https://api.github.com/gists/"+gistID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var gistData struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}

	json.NewDecoder(resp.Body).Decode(&gistData)

	var p Pet
	err = json.Unmarshal([]byte(gistData.Files["diana.json"].Content), &p)
	p.Tasks = make(chan Task, 10)
	return &p, err
}

func FetchArchiveData() string {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	req.Header.Set("Accept", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return "Sorry! I can't connect to the archive."
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
