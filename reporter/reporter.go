package reporter

import (
	"fmt"
	"os"
	"strings"
)

// ReportItem 表示报告中的一个问题项
type ReportItem struct {
	RuleID      string
	Description string
	Severity    string
	Location    string
}

// Reporter 用于生成报告
type Reporter struct {
	reportItems []ReportItem
}

// NewReporter 创建一个新的 Reporter 实例
func NewReporter() *Reporter {
	return &Reporter{}
}

// AddReportItem 向报告中添加一个新的问题项
func (r *Reporter) AddReportItem(item ReportItem) {
	r.reportItems = append(r.reportItems, item)
}

// GenerateReport 生成报告并将其输出到指定的文件
func (r *Reporter) GenerateReport(outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// 生成报告标题
	title := "Python SAST Report"
	hr := strings.Repeat("=", len(title))
	_, _ = fmt.Fprintf(f, "%s\n%s\n%s\n\n", hr, title, hr)

	// 输出每个问题项
	for _, item := range r.reportItems {
		_, _ = fmt.Fprintf(f, "Rule ID: %s\n", item.RuleID)
		_, _ = fmt.Fprintf(f, "Description: %s\n", item.Description)
		_, _ = fmt.Fprintf(f, "Severity: %s\n", item.Severity)
		_, _ = fmt.Fprintf(f, "Location: %s\n", item.Location)
		_, _ = fmt.Fprintln(f, strings.Repeat("-", 80))
	}

	return nil
}
