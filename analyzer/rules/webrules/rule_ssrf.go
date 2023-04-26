package rules

import (
	"net"
	"net/url"
	"strings"

	"github.com/yourusername/yourproject/parser"
	"github.com/yourusername/yourproject/reporter"
)

type RuleSSRF struct {
	reporter *reporter.Reporter
}

// NewRuleSSRF 创建并返回一个新的RuleSSRF实例
func NewRuleSSRF(reporter *reporter.Reporter) *RuleSSRF {
	return &RuleSSRF{
		reporter: reporter,
	}
}

// Apply 应用规则，并将结果添加到报告中
func (r *RuleSSRF) Apply(ast *parser.ASTNode) {
	r.CheckCondition(ast)
}

// CheckCondition 检查SSRF的存在并报告
func (r *RuleSSRF) CheckCondition(ast *parser.ASTNode) {
	urlNodes := ast.FindAll(parser.NodeTypeURL)

	for _, node := range urlNodes {
		urlStr := node.Value()

		// 解析URL，检查是否为远程地址
		u, err := url.Parse(urlStr)
		if err != nil {
			continue
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		if u.Hostname() == "" {
			continue
		}
		hostIPs, err := net.LookupIP(u.Hostname())
		if err != nil {
			continue
		}
		hostIsLocal := false
		for _, hostIP := range hostIPs {
			if hostIP.IsLoopback() || hostIP.IsLinkLocalUnicast() || hostIP.IsLinkLocalMulticast() {
				hostIsLocal = true
				break
			}
		}
		if hostIsLocal {
			continue
		}

		// 检查协议是否为http或https
		if u.Scheme != "http" && u.Scheme != "https" {
			r.reporter.AddIssue(node.Position().Filename, node.Position().Line, "Potential SSRF vulnerability: "+urlStr)
			continue
		}

		// 检查端口是否在安全范围内
		port := u.Port()
		if port == "" {
			if u.Scheme == "http" {
				port = "80"
			} else {
				port = "443"
			}
		}
		if port != "80" && port != "443" && port != "8000" && port != "8080" && port != "8443" {
			r.reporter.AddIssue(node.Position().Filename, node.Position().Line, "Potential SSRF vulnerability: "+urlStr)
			continue
		}

		// 检查路径是否为可信任的
		if strings.HasPrefix(u.Path, "/") {
			continue
		}
		r.reporter.AddIssue(node.Position().Filename, node.Position().Line, "Potential SSRF vulnerability: "+urlStr)
	}
}
