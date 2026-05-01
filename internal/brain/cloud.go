package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GistUpdate struct {
	Files map[string]map[string]string `json:"files"`
}

func getCloudCredentials() (string, string) {
	token := os.Getenv("DIANA_GITHUB_TOKEN")
	gistID := os.Getenv("DIANA_GIST_ID")

	if token == "" || gistID == "" {
		return "", ""
	}

	return token, gistID
}

func SaveToCloud(p *Pet) error {
	token, gistID := getCloudCredentials()
	if token == "" {
		return fmt.Errorf("Missing cloud credentials")
	}

	data, _ := json.MarshalIndent(p, "", "  ")

	payload := GistUpdate{
		Files: map[string]map[string]string{
			"diana.json": {"content": string(data)},
		},
	}

	jsonData, _ := json.Marshal(payload)
	url := "https://api.github.com/gists/" + gistID
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func LoadFromCloud() (*Pet, error) {
	token, gistID := getCloudCredentials()

	req, _ := http.NewRequest("GET", "https://api.github.com/gists/" + gistID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var gistData struct {
		Files map[string]struct{
			Content string `json:"content"`
		} `json:"files"`
	}

	json.NewDecoder(resp.Body).Decode(&gistData)

	var p Pet
	err = json.Unmarshal([]byte(gistData.Files["diana.json"].Content), &p)
	return &p, err
}