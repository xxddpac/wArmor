local _M = {}

_M.shared_dict_waf_rule = ngx.shared.waf_rule
_M.shared_dict_waf_config = ngx.shared.waf_config
_M.shared_dict_waf_ip = ngx.shared.waf_ip

_M.redirect_html = "/usr/local/openresty/waf/redirect.html"

_M.regexp_option = "isjo"

_M.unknown = "unknown"

_M.rule_engine = {
    url = "http://X.X.X.X:9999/api/v1/ip", -- 替换规则引擎IP
    method = "post",
    max_idle_timeout = 100,
    pool_size = 10
}

_M.db_table = {
    rule = "rule",
    config = "config",
    ip = "ip"
}

_M.geoip_db = {
    language = "zh-CN",
    file = "/usr/local/share/GeoIP/GeoLite2-City.mmdb"
}

_M.desc = {
    mode = {
        [1] = "拦截模式",
        [2] = "监控模式",
        [3] = "旁路模式"
    },
    black_white_list = {
        [1] = "白名单IP",
        [2] = "黑名单IP",
    },
    rule_action = {
        [1] = "拒绝",
        [2] = "重定向",
        [3] = "允许"
    },
    rule_type = {
        [1] = "跨站脚本攻击",
        [2] = "webshell",
        [3] = "sql注入",
        [4] = "路径遍历",
        [5] = "不安全请求方法",
        [6] = "白名单URL",
        [7] = "黑名单URL",
        [8] = "爬虫",
        [9] = "敏感信息",
        [10] = "跨站请求伪造",
        [11] = "命令注入",
        [12] = "拒绝服务",
        [13] = "身份验证绕过",
        [14] = "逻辑缺陷",
        [15] = "其他"
    },
    severity = {
        [1] = "严重",
        [2] = "高危",
        [3] = "中危",
        [4] = "低危",
        [5] = "信息"
    }
}

_M.db = {
    set_timeout = 2000,
    host = "X.X.X.X",  -- 替换mysql服务器IP
    username = "root",        -- 替换用户名
    password = "XXXX", -- 替换密码
    port = "3306",
    database = "waf",
    charset = "utf8",
    max_packet_size = 1024 * 1024
}

_M.redis = {
    connect_timeout = 10000,
    set_timeout = 30000,
    read_timeout = 300000,
    redis_ssl = false,
    pool_size = 100,
    host = "X.X.X.X", -- 替换redis服务器IP
    password = "",
    port = "6379",
    channel = "waf"
}

_M.log = {
    host = "X.X.X.X", -- 替换syslog-NG服务器IP
    port = 514,
    sock_type = "udp",
    flush_limit = 1,
    drop_limit = 5678,
    pool_size = 100,
    timeout = 1000,
    program = "wArmor"
}

return _M
