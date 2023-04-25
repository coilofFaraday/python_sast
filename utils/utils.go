package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

// ReadFile 读取文件并返回其内容
func ReadFile(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IsPythonFile 检查给定的文件路径是否为Python文件
func IsPythonFile(path string) bool {
	return strings.HasSuffix(path, ".py")
}

// ListPythonFiles 递归列出指定目录下的所有Python文件
func ListPythonFiles(dir string) ([]string, error) {
	var pythonFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && IsPythonFile(path) {
			pythonFiles = append(pythonFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return pythonFiles, nil
}
func Indent(level int, code string) string {
	indentation := strings.Repeat("    ", level)
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%s%s", indentation, line)
	}
	return strings.Join(lines, "\n")
}

// EscapeString 用于转义字符串中的特殊字符
func EscapeString(s string) string {
	var builder strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			builder.WriteString(`\\`)
		case '"':
			builder.WriteString(`\"`)
		case '\n':
			builder.WriteString(`\n`)
		case '\t':
			builder.WriteString(`\t`)
		case '\r':
			builder.WriteString(`\r`)
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// DeepCopyASTNode 创建并返回一个新的ASTNode，它是原始节点的深度拷贝
func DeepCopyASTNode(node *ASTNode) *ASTNode {
	// ...实现深度拷贝逻辑
}

// TraverseASTNode 遍历给定的AST节点（以及其子节点），对每个节点调用指定的回调函数
func TraverseASTNode(node *ASTNode, callback func(*ASTNode)) {
	callback(node)
	// ...实现递归遍历逻辑
}