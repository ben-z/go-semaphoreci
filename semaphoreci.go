package semaphoreci

import (
	"encoding/json"
	"fmt"
	"time"
)

type Project struct {
	Id        int
	HashId    string `json:"hash_id"`
	Name      string
	Owner     string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	branches  []Branch
	Client    *Client
}

type BriefBranchInfo struct {
	ID         int    `json:"id"`
	BranchName string `json:"name"`
	BranchURL  string `json:"branch_url"`
}

type Branch struct {
	BranchName       string `json:"branch_name"`
	BranchURL        string `json:"branch_url"`
	BranchStatusURL  string `json:"branch_status_url"`
	BranchHistoryURL string `json:"branch_history_url"`
	ProjectName      string `json:"project_name"`
	BuildURL         string `json:"build_url"`
	BuildInfoURL     string `json:"build_info_url"`
	BuildNumber      int    `json:"build_number"`
	Result           string `json:"result"`
	StartedAt        string `json:"started_at"`
	FinishedAt       string `json:"finished_at"`
}

type Commit struct {
	Id         string
	URL        string
	AuthorName string `json:"author_name"`
	AuthorMail string `json:"author_mail"`
	Message    string
	Timestamp  string
}

type BranchStatus struct {
	Branch
	Commit Commit
}

type BranchHistory struct {
	Branch
	Pagination Pagination
	Builds     []Build
}

type Pagination struct {
	TotalEntries int  `json:"total_entries"`
	TotalPages   int  `json:"total_pages"`
	PerPage      int  `json:"per_page"`
	CurrentPage  int  `json:"current_page"`
	IsFirstPage  bool `json:"first_page"`
	IsLastPage   bool `json:"last_page"`
}

type Build struct {
	BranchInfo
	Commit Commit
}

type BuildInfo struct {
	ProjectName string `json:"project_name"`
	BrancName   string `json:"branch_name"`
	Number      int
	Result      string
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	StartedAt   string `json:"started_at"`
	FinishedAt  string `json:"finished_at"`
	HtmlURL     string `json:"html_url"`
	Commits     []Commit
}

type BranchInfo struct {
	BuildURL     string `json:"build_url"`
	BuildInfoURL string `json:"build_info_url"`
	BuildNumber  int    `json:"build_number"`
	Result       string
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
}

type BuildLog struct {
	Threads      []Thread
	BuildInfoURL string `json:"build_info_url"`
}

type Thread struct {
	Number   int
	Commands []Command
}

type Command struct {
	Name       string
	Result     string
	Output     string
	Duration   string
	StartTime  *time.Time `json:"start_time"`
	FinishTime *time.Time `json:"finish_time"`
}

func (c *Client) Projects() ([]Project, error) {
	data := []Project{}
	body, _, _ := c.GetRequest("projects", nil)
	err := json.Unmarshal(body, &data)
	return data, err
}

func (c *Client) Project(hash_id string) *Project {
	return &Project{HashId: hash_id, Client: c}
}

func (p *Project) Branches() ([]BriefBranchInfo, error) {
	var data []BriefBranchInfo
	url := fmt.Sprintf("projects/%v/branches", p.HashId)
	body, _, _ := p.Client.GetRequest(url, nil)
	err := json.Unmarshal(body, &data)
	return data, err
}

func (p *Project) BranchStatus(branch_id interface{}) (BranchStatus, error) {
	data := BranchStatus{}
	url := fmt.Sprintf("projects/%v/%v/status", p.HashId, branch_id)
	body, _, _ := p.Client.GetRequest(url, nil)
	err := json.Unmarshal(body, &data)
	return data, err
}

func (p *Project) BranchHistory(branch_id interface{}) (*BranchHistory, error) {
	data := BranchHistory{}
	url := fmt.Sprintf("projects/%v/%v", p.HashId, branch_id)
	body, header, _ := p.Client.GetRequest(url, nil)
	err := json.Unmarshal(body, &data)
	if err != nil {
		return &data, err
	}
	err = json.Unmarshal([]byte(header.Get("pagination")), &(data.Pagination))
	return &data, err
}

func (p *Project) BranchHistoryNextPage(branchHistory *BranchHistory) (*BranchHistory, error) {
	data := BranchHistory{}
	currentPage := branchHistory.Pagination.CurrentPage
	url := fmt.Sprintf("projects/%v/%v", p.HashId, branchHistory.BranchName)
	params := map[string]interface{}{"page": currentPage + 1}
	body, header, _ := p.Client.GetRequest(url, &params)
	err := json.Unmarshal(body, &data)
	if err != nil {
		return &data, err
	}
	err = json.Unmarshal([]byte(header.Get("Pagination")), &(data.Pagination))
	return &data, err
}

func (p *Project) BuildInfo(branch_id interface{}, build_num int) (BuildInfo, error) {
	data := BuildInfo{}
	url := fmt.Sprintf("projects/%v/%v/builds/%v", p.HashId, branch_id, build_num)
	body, _, _ := p.Client.GetRequest(url, nil)
	err := json.Unmarshal(body, &data)
	return data, err
}

func (p *Project) BuildLog(branch_id interface{}, build_num int) (BuildLog, error) {
	data := BuildLog{}
	url := fmt.Sprintf("projects/%v/%v/builds/%v/log", p.HashId, branch_id, build_num)
	body, _, _ := p.Client.GetRequest(url, nil)
	err := json.Unmarshal(body, &data)
	return data, err
}
