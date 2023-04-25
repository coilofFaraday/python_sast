package reporter

import (
	"fmt"
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

// NewReporter 创建并返回一个新的Reporter实例
func NewReporter() *Reporter {
	return &Reporter{
		issues: make([]Issue, 0),
	}
}

// AddIssue 添加一个问题到报告中
func (r *Reporter) AddIssue(file string, line int, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues = append(r.issues, Issue{File: file, Line: line, Message: message})
}

// Print 打印报告的问题到标准输出
func (r *Reporter) Print() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, issue := range r.issues {
		fmt.Printf("%s:%d: %s\n", issue.File, issue.Line, issue.Message)
	}
}

// Clear 清除报告中的所有问题
func (r *Reporter) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues = make([]Issue, 0)
}

// SaveToFile 保存报告到指定的文件
func (r *Reporter) SaveToFile(filePath string) error {
	// ...实现保存报告到文件的逻辑
}

// ConvertToJSON 将报告转换为JSON格式并返回
func (r *Reporter) ConvertToJSON() (string, error) {
	// ...实现将报告转换为JSON格式的逻辑
}

// ConvertToXML 将报告转换为XML格式并返回
func (r *Reporter) ConvertToXML() (string, error) {
	// ...实现将报告转换为XML格式的逻辑
}

// FilterIssues 按条件过滤报告中的问题并返回新的报告
func (r *Reporter) FilterIssues(filterFunc func(Issue) bool) *Reporter {
	// ...实现过滤问题的逻辑
}
