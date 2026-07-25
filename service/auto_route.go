package service

import (
	"strings"

	"github.com/QuantumNous/new-api/dto"
)

// 虚拟模型 auto：启发式路由（v1）
// 零额外上游调用、零额外延迟，按任务特征把 auto 改写为最合适的实际模型。
// 改写发生在 relay 控制器验证请求之后、生成 relayInfo 与计费之前，
// 因此计费、渠道、日志全部按实际路由到的模型结算。
const (
	AutoModelName          = "auto"
	AutoRouteModelMultimodal = "qwen3.5-omni-flash" // 含图片/音频输入
	AutoRouteModelCode       = "kimi-k2.7-code"     // 编程任务
	AutoRouteModelComplex    = "qwen3.7-max"        // 长文/深度任务
	AutoRouteModelSimple     = "qwen3.6-flash"      // 简短闲聊
	AutoRouteModelDefault    = "qwen3.7-plus"       // 默认
)

var autoCodeSignals = []string{
	"```", "def ", "function", "import ", "class ", "traceback", "stacktrace",
	"代码", "函数", "编程", "写个脚本", "报错", "bug", "debug", "编译", "算法实现",
	"sql", "python", "java", "golang", "javascript", "typescript", "rust", "c++", "regex", "正则",
}

var autoComplexSignals = []string{
	"深度分析", "研究报告", "论文", "合同", "商业计划", "可行性", "逐条", "长篇",
	"详细分析", "综合分析", "系统设计", "架构设计", "comprehensive", "research", "whitepaper",
}

func IsAutoModel(modelName string) bool {
	return strings.EqualFold(strings.TrimSpace(modelName), AutoModelName)
}

// ResolveAutoRequestModel 若请求模型为 auto，按启发式规则改写为实际模型并返回该模型；否则返回原模型。
func ResolveAutoRequestModel(req *dto.GeneralOpenAIRequest) string {
	if req == nil || !IsAutoModel(req.Model) {
		if req == nil {
			return ""
		}
		return req.Model
	}

	hasMedia := false
	var textBuilder strings.Builder
	for i := range req.Messages {
		message := &req.Messages[i]
		for _, part := range message.ParseContent() {
			switch part.Type {
			case dto.ContentTypeText:
				if message.Role == "user" {
					textBuilder.WriteString(part.Text)
					textBuilder.WriteString("\n")
				}
			case dto.ContentTypeImageURL, dto.ContentTypeInputAudio:
				hasMedia = true
			}
		}
	}
	if hasMedia {
		req.Model = AutoRouteModelMultimodal
		return req.Model
	}

	text := textBuilder.String()
	lower := strings.ToLower(text)
	runeLen := len([]rune(text))

	codeHits := 0
	for _, signal := range autoCodeSignals {
		if strings.Contains(lower, signal) {
			codeHits++
			if codeHits >= 2 {
				req.Model = AutoRouteModelCode
				return req.Model
			}
		}
	}
	if codeHits == 1 && runeLen > 200 {
		req.Model = AutoRouteModelCode
		return req.Model
	}

	if runeLen > 1200 {
		req.Model = AutoRouteModelComplex
		return req.Model
	}
	for _, signal := range autoComplexSignals {
		if strings.Contains(lower, signal) {
			req.Model = AutoRouteModelComplex
			return req.Model
		}
	}

	if runeLen > 0 && runeLen <= 30 {
		req.Model = AutoRouteModelSimple
		return req.Model
	}

	req.Model = AutoRouteModelDefault
	return req.Model
}
