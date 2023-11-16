package global

type WafRuleType int

const (
	Xss WafRuleType = iota + 1
	WebShell
	SQLInjection
	PathTraversal
	UnsafeHttpMethod
	WhiteURL
	BlackURL
	Spider
	SensitiveInformationMonitor
	CSRF
	CommandInjection
	DoSDenialOfService
	AuthenticationBypass
	LogicFlaw
	Other
)

func (w WafRuleType) String() string {
	switch w {
	case Xss:
		return "跨站脚本攻击"
	case SQLInjection:
		return "sql注入"
	case WebShell:
		return "恶意文件上传"
	case PathTraversal:
		return "路径遍历"
	case UnsafeHttpMethod:
		return "不安全的请求方法"
	case WhiteURL:
		return "白名单URL"
	case BlackURL:
		return "黑名单URL"
	case Spider:
		return "爬虫"
	case SensitiveInformationMonitor:
		return "敏感信息"
	case CSRF:
		return "跨站请求伪造"
	case CommandInjection:
		return "命令注入"
	case DoSDenialOfService:
		return "拒绝服务"
	case AuthenticationBypass:
		return "身份验证绕过"
	case LogicFlaw:
		return "逻辑缺陷"
	case Other:
		return "其他"
	default:
		return "未知"
	}
}

type WafRuleAction int

const (
	Deny WafRuleAction = iota + 1
	Redirect
	Pass
)

func (w WafRuleAction) String() string {
	switch w {
	case Deny:
		return "拒绝"
	case Redirect:
		return "重定向"
	case Pass:
		return "允许"
	default:
		return "未知"
	}
}

type WafConfigType int

const (
	Block WafConfigType = iota + 1
	Alert
	Bypass
)

func (w WafConfigType) String() string {
	switch w {
	case Block:
		return "阻断模式"
	case Alert:
		return "监控模式"
	case Bypass:
		return "旁路模式"
	default:
		return "未知"
	}
}

func (w WafConfigType) Comment() string {
	switch w {
	case Block:
		return "【当前模式为阻断模式】阻断命中规则的恶意请求并记录"
	case Alert:
		return "【当前模式为监控模式】仅记录命中规则的恶意请求"
	case Bypass:
		return "【当前模式为旁路模式】无记录无阻断"
	default:
		return "未知"
	}
}

type WafIpType int

const (
	WhiteList WafIpType = iota + 1
	BlackList
)

func (w WafIpType) String() string {
	switch w {
	case WhiteList:
		return "白名单"
	case BlackList:
		return "黑名单"
	default:
		return "未知"
	}
}

type WafBlockType int

const (
	Permanent WafBlockType = iota + 1
	Temporary
)

func (w WafBlockType) String() string {
	switch w {
	case Permanent:
		return "永久封禁"
	case Temporary:
		return "临时封禁"
	default:
		return "未知"
	}
}

type WafRuleSeverity int

const (
	Serious WafRuleSeverity = iota + 1
	High
	Medium
	Low
	Info
)

func (w WafRuleSeverity) String() string {
	switch w {
	case Serious:
		return "严重"
	case High:
		return "高危"
	case Medium:
		return "中危"
	case Low:
		return "低危"
	case Info:
		return "提示"
	default:
		return "未知"
	}
}

type WafRuleVariable int

const (
	HttpRequestArgs WafRuleVariable = iota + 1
	HttpRequestMethod
	HttpRequestURI
	HttpRequestBody
	HttpRequestHeaders
	Response
)

func (w WafRuleVariable) String() string {
	switch w {
	case HttpRequestArgs:
		return "请求参数"
	case HttpRequestBody:
		return "请求体"
	case HttpRequestMethod:
		return "请求方法"
	case HttpRequestURI:
		return "请求路径"
	case HttpRequestHeaders:
		return "请求头"
	case Response:
		return "响应体"
	default:
		return "未知"
	}
}
