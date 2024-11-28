package dto

import (
	"congchat-user/model"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// CreateMomentRequest 发送发布朋友圈动态请求
type CreateMomentRequest struct {
	Content string   `gorm:"type:text" binding:"required,min=1,max=255"`
	UserID  uint     `gorm:"index" binding:"required"`
	ImgURL  []string `gorm:"type:text[]" binding:"required,max=9,dive,required,url"` //图片或视频地址。用于前端加载视频
	Goods   int      `gorm:"not null"  binding:"required"`
	GoodsID []int    `gorm:"not null"  binding:"required"`
}

// EditMomentRequest 编辑朋友圈动态请求
type EditMomentRequest struct {
	model.Moment
}

// GetMomentRequest 查看朋友圈动态请求（）本质应该是一个列表（含有朋友动态）//滑动查询
type GetMomentRequest struct {
	StartIdx int `json:"start_idx,omitempty"` //0  -  9
	Limit    int `json:"limit,omitempty"`     //10
	//条件
	FriendID  uint   `json:"friend_id,omitempty"`
	SearchStr string `json:"search_str,omitempty"`
}

// DeleteMomentRequest 删除朋友圈动态请求
type DeleteMomentRequest struct {
	Ids []int `json:"ids"`
}

// ValidateImgURLs 自定义验证函数，用于验证ImgURLs切片的长度不超过9
func ValidateImgURLs(ImgURL []string) error {
	if len(ImgURL) > 9 {
		return errors.New("最多只能上传9张图片")
	}
	for _, imgURL := range ImgURL {
		_, err := url.ParseRequestURI(imgURL)
		if err != nil {
			return fmt.Errorf("无效的URL: %s", imgURL)
		}
	}
	return nil
}

// ValidateContent 校验文本内容是否包含敏感词或不符合要求的内容
func ValidateContent(content string) error {
	// 这里可以添加具体的敏感词检测逻辑，或者使用现成的敏感词库
	// 例如，如果内容包含“暴力”、“黄色”等关键词，则返回错误
	// 为了简化示例，这里只检查内容是否为空
	if content == "" {
		return errors.New("内容不能为空")
	}
	// 假设我们有一个敏感词（关键字）列表，这里用简单的字符串切片模拟
	sensitiveWords := []string{"暴力", "黄色", "广告"}
	for _, word := range sensitiveWords {
		if strings.Contains(content, word) {
			return fmt.Errorf("内容包含敏感词：%s", word)
		}
	}
	return nil
}

// CreateMomentClear 检查 Contents 和 UserID 和 MomentID是否都不为零
func CreateMomentClear(d *CreateMomentRequest) error {
	// 文本参数校验 查看是否符合格式
	if err := ValidateContent(d.Content); err != nil {
		return errors.New("内容验证格式错误")
	}

	//检验下图片的输入是否大于9张图片的处理函数
	if err := ValidateImgURLs(d.ImgURL); err != nil {
		return errors.New("图片验证错误")
	}
	// 其他参数校验，例如检查UserID是否为0（假设0是无效的用户ID）
	if d.UserID == 0 {
		return errors.New("用户ID不能为0")
	}
	return nil
}

// EditMomentClear 检查 Contents 和 UserID 和 MomentID是否都不为零
func EditMomentClear(d *EditMomentRequest) error {
	// 文本参数校验 查看是否符合格式
	if err := ValidateContent(d.Content); err != nil {
		return errors.New("内容验证格式错误")
	}

	//检验下图片的输入是否大于9张图片的处理函数
	if err := ValidateImgURLs(d.ImgURL); err != nil {
		return errors.New("图片验证错误")
	}
	return nil
}
