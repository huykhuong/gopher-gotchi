package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LinearResponse struct {
	Data struct {
		Viewer struct {
			AssignedIssues struct {
				Nodes []struct {
					Identifier string `json:"identifier"`
					Title      string `json:"title"`
					State     struct {
						Name string `json:"name"`
					} `json:"state"`
				} `json:"nodes"`
			} `json:"assignedIssues"`
		} `json:"viewer"`
	} `json:"data"`
}

func FetchObjectives() (string, error) {
	apiKey := os.Getenv("DIANA_LINEAR_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing API key variable")
	}

	query := map[string]string{
		"query": `{
			viewer {
				assignedIssues(filter: { state: { type: { neq: "completed" } } }) {
					nodes {
						identifier
						title
						state { name }
					}
				}
			}
		}`,
	}

	jsonData, _ := json.Marshal(query)
	req, _ := http.NewRequest("POST", "https://api.linear.app/graphql", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5*time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result LinearResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	nodes := result.Data.Viewer.AssignedIssues.Nodes
	if len(nodes) == 0 {
		return "✨ All clear, Huy! No active objectives found in this sector.", nil
	}

	output := ""
	for _, issue := range nodes {
		output += fmt.Sprintf("> [%s] %s (%s)\n\n", issue.Identifier, issue.Title, issue.State.Name)
	}
	
	return output, nil
}