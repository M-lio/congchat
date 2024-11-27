package dto

import (
	"errors"
	"fmt"
	"strings"
)

type CreateCommentRequest struct {
	Contents string `json:"contents" binding:"required,min=1,max=255"` // 假设内容至少1个字符，最多255个字符
	UserID   uint   `json:"user_id" binding:"required"`                // 用户ID必填
	MomentID uint   `json:"moment_id" binding:"required"`              // 动态ID必填
}

type DeleteCommentRequest struct {
	Ids []int `json:"ids"`
}

// CreateClear 检查 Contents 和 UserID 和 MomentID是否都不为零
func CreateClear(d *CreateCommentRequest) error {
	if err := ValidateContents(d.Contents); err != nil {
		return errors.New("内容验证格式错误")
	}

	if d.MomentID == 0 {
		return errors.New("朋友圈ID不能为0")
	}
	if d.UserID == 0 {
		return errors.New("用户ID不能为0")
	}
	return nil
}

// ValidateContents 校验文本内容是否包含敏感词或不符合要求的内容
func ValidateContents(contents string) error {
	// 这里可以添加具体的敏感词检测逻辑，或者使用现成的敏感词库
	// 例如，如果内容包含“暴力”、“黄色”等关键词，则返回错误
	// 为了简化示例，这里只检查内容是否为空
	if contents == "" {
		return errors.New("内容不能为空")
	}
	// 假设我们有一个敏感词（关键字）列表，这里用简单的字符串切片模拟
	sensitiveWords := []string{"暴力", "黄色", "广告"}
	for _, word := range sensitiveWords {
		if strings.Contains(contents, word) {
			return fmt.Errorf("内容包含敏感词：%s", word)
		}
	}
	return nil
}
