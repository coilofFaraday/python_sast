package reporter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sync"
)

type Issue struct {
	File    string
	Line    int
	Message string
}

type Reporter struct {
	mu     sync.Mutex
	issues []Issue
}

// NewReporter creates and returns a new Reporter instance
func NewReporter() *Reporter {
	return &Reporter{
		issues: make([]Issue, 0),
	}
}

// AddIssue adds an issue to the report
func (r *Reporter) AddIssue(file string, line int, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues = append(r.issues, Issue{File: file, Line: line, Message: message})
}

// Print prints the report's issues to standard output
func (r *Reporter) Print() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, issue := range r.issues {
		fmt.Printf("%s:%d: %s\n", issue.File, issue.Line, issue.Message)
	}
}

// Clear clears all the issues from the report
func (r *Reporter) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues = make([]Issue, 0)
}

// SaveToFile saves the report to the specified file
func (r *Reporter) SaveToFile(filePath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, issue := range r.issues {
		_, err := fmt.Fprintf(file, "%s:%d: %s\n", issue.File, issue.Line, issue.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

// ConvertToJSON converts the report to JSON format and returns it as a string
func (r *Reporter) ConvertToJSON() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	jsonBytes, err := json.Marshal(r.issues)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ConvertToXML converts the report to XML format and returns it as a string
func (r *Reporter) ConvertToXML() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	xmlBytes, err := xml.Marshal(r.issues)
	if err != nil {
		return "", err
	}
	return string(xmlBytes), nil
}

// FilterIssues filters the issues in the report by the specified condition and returns a new report
func (r *Reporter) FilterIssues(filterFunc func(Issue) bool) *Reporter {
	r.mu.Lock()
	defer r.mu.Unlock()
	newReporter := NewReporter()
	for _, issue := range r.issues {
		if filterFunc(issue) {
			newReporter.AddIssue(issue.File, issue.Line, issue.Message)
		}
	}
	return newReporter
}

// Issues 返回报告中的所有问题
func (r *Reporter) Issues() []Issue {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.issues
}

// NumIssues 返回报告中的问题数量
func (r *Reporter) NumIssues() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.issues)
}

// GetIssuesByFile 返回指定文件中的所有问题
func (r *Reporter) GetIssuesByFile(file string) []Issue {
	r.mu.Lock()
	defer r.mu.Unlock()
	issues := make([]Issue, 0)
	for _, issue := range r.issues {
		if issue.File == file {
			issues = append(issues, issue)
		}
	}
	return issues
}

// GetIssuesByLine 返回指定行号的所有问题
func (r *Reporter) GetIssuesByLine(line int) []Issue {
	r.mu.Lock()
	defer r.mu.Unlock()
	issues := make([]Issue, 0)
	for _, issue := range r.issues {
		if issue.Line == line {
			issues = append(issues, issue)
		}
	}
	return issues
}

// GetIssuesByMessage 返回包含指定字符串的所有问题
func (r *Reporter) GetIssuesByMessage(message string) []Issue {
	r.mu.Lock()
	defer r.mu.Unlock()
	issues := make([]Issue, 0)
	for _, issue := range r.issues {
		if contains(issue.Message, message) {
			issues = append(issues, issue)
		}
	}
	return issues
}

// contains 判断字符串s是否包含子字符串substr
func contains(s, substr string) bool {
	return len(substr) <= len(s) && s[len(s)-len(substr):] == substr
}
