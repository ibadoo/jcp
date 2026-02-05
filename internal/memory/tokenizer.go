package memory

import (
	"strings"

	"github.com/yanyiwu/gojieba"
)

// Tokenizer 分词器接口
type Tokenizer interface {
	Extract(text string, topK int) []string
	Cut(text string) []string
}

// JiebaTokenizer 基于 jieba 的分词器
type JiebaTokenizer struct {
	jieba     *gojieba.Jieba
	stopWords map[string]bool
}

// NewJiebaTokenizer 创建 jieba 分词器
func NewJiebaTokenizer() *JiebaTokenizer {
	return &JiebaTokenizer{
		jieba:     gojieba.NewJieba(),
		stopWords: defaultStopWords(),
	}
}

// Free 释放资源
func (t *JiebaTokenizer) Free() {
	if t.jieba != nil {
		t.jieba.Free()
	}
}

// Extract 提取关键词（使用 TF-IDF）
func (t *JiebaTokenizer) Extract(text string, topK int) []string {
	words := t.jieba.ExtractWithWeight(text, topK*2)
	result := make([]string, 0, topK)
	for _, w := range words {
		if !t.stopWords[w.Word] && len([]rune(w.Word)) >= 2 {
			result = append(result, w.Word)
			if len(result) >= topK {
				break
			}
		}
	}
	return result
}

// Cut 分词
func (t *JiebaTokenizer) Cut(text string) []string {
	words := t.jieba.Cut(text, true)
	result := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w != "" && !t.stopWords[w] && len([]rune(w)) >= 2 {
			result = append(result, w)
		}
	}
	return result
}

// defaultStopWords 默认停用词
func defaultStopWords() map[string]bool {
	words := []string{
		"的", "是", "在", "了", "和", "与", "或", "这", "那", "有",
		"个", "我", "你", "他", "她", "它", "们", "吗", "呢", "吧",
		"啊", "哦", "嗯", "呀", "哈", "哪", "什么", "怎么", "为什么",
		"可以", "可能", "应该", "需要", "能够", "已经", "正在",
		"一个", "一些", "这个", "那个", "这些", "那些", "如果",
		"但是", "因为", "所以", "虽然", "然后", "而且", "或者",
		"不是", "没有", "不会", "不能", "还是", "就是", "只是",
	}
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}
