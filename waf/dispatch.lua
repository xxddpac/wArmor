local conf = require "config"
local cjson = require "cjson"
local utils = require "lib.utils"
local logger = require "lib.log"
local http = require "lib.http"

local _M = {}

local function deny()
    ngx.exit(ngx.HTTP_FORBIDDEN)
end

local function pass()
end

local function redirect()
    ngx.status = ngx.HTTP_BAD_GATEWAY
    ngx.header.content_type = "text/html; charset=UTF-8"
    local html_content = utils.read_html_file_to_string()

    if html_content then
        local substitutions = {
            ["\\$client_ip"] = ngx.ctx.ip,
            ["\\$client_ua"] = ngx.ctx.ua,
            ["\\$request_id"] = ngx.ctx.uuid,
            ["\\$blocked_time"] = ngx.ctx.now,
        }
        for pattern, replacement in pairs(substitutions) do
            html_content = ngx.re.gsub(html_content, pattern, replacement, conf.regexp_option)
        end
        ngx.say(html_content)
    end

    ngx.exit(ngx.status)
end

function _M.match(rule, result)
    local rules = cjson.decode(rule.rules)
    if rule.rules_operation == "or" then
        for _, val in pairs(result) do
            local res
            if type(val) == "table" then
                res = table.concat(val, ",")
            else
                res = val
            end
            for _, regex_rule in ipairs(rules) do
                if ngx.re.find(res, regex_rule, conf.regexp_option) then
                    return 1024
                end
            end
        end
        return nil
    else
        for _, val in pairs(result) do
            local res
            if type(val) == "table" then
                res = table.concat(val, ",")
            else
                res = val
            end
            for _, regex_rule in ipairs(rules) do
                if not ngx.re.find(res, regex_rule, conf.regexp_option) then
                    return nil
                end
            end
        end
        return 1024
    end
end

local action = {
    [1] = deny,
    [2] = redirect,
    [3] = pass
}

local function log_and_exit(rule)
    local rule_type
    local severity
    local msg = {
        rule_id = rule.id,
        rule_action = conf.desc.rule_action[rule.rule_action],
        request_data = rule.request_data or ""
    }
    if rule.ip_type ~= nil then
        -- 匹配黑白名单
        rule_type = conf.desc.black_white_list[rule.ip_type]
        if rule.ip_type == 1 then
            severity = conf.desc.severity[5]
        else
            severity = conf.desc.severity[1]
        end
    else
        --匹配动态规则
        rule_type = conf.desc.rule_type[rule.rule_type]
        severity = conf.desc.severity[rule.severity]
    end
    msg.rule_type = rule_type
    msg.severity = severity
    logger.send(logger.new_log(msg))
    if ngx.ctx.mode == 1 then -- 阻断模式
        action[rule.rule_action]()
    end
end

local function request_args(rules)
    local result = ngx.req.get_uri_args()

    if not next(result) then
        return nil
    end

    for _, rule in ipairs(rules) do
        if _M.match(rule, result) then
            log_and_exit(rule)
            return 1024
        end
    end

    return nil
end


local function unsafe_http_method(rules)
    local result = ngx.req.get_method()

    if not result then
        return nil
    end

    for _, rule in ipairs(rules) do
        if _M.match(rule, { method = result }) then
            log_and_exit(rule)
            return 1024
        end
    end

    return nil
end

local function white_black_url(rules)
    local result = ngx.var.uri

    if not result then
        return nil
    end

    for _, rule in ipairs(rules) do
        if _M.match(rule, { white_black_url = result }) then
            if rule.rule_type == 7 then
                -- 命中黑名单url或robots.txt禁止访问的url,异步更新黑名单临时封禁IP 1天
                local request_body = {
                    ip_address = ngx.ctx.ip,
                    operator = conf.log.program,
                    block_type = 2,
                    ip_type = 2,
                    expire_time_tag = 3,
                    comment = "当命中黑名单url时waf自动拉黑请求IP并临时封禁1天",
                }
                http.send(conf.rule_engine.url, conf.rule_engine.method, cjson.encode(request_body))
            end
            log_and_exit(rule)
            return 1024
        end
    end

    return nil
end


local function request_headers(rules)
    local referer = ngx.var.http_referer

    if not referer and ngx.ctx.ua == conf.unknown then
        return nil
    end

    for _, rule in ipairs(rules) do
        if _M.match(rule, { referer = referer }) or _M.match(rule, { ua = ngx.ctx.ua }) then
            log_and_exit(rule)
            return 1024
        end
    end

    return nil
end


local function request_body(rules)
    local contentType = ngx.var.http_content_type
    local contentLength = tonumber(ngx.var.http_content_length)
    local boundary = nil
    if contentType then
        local bfrom, bto = ngx.re.find(contentType, "\\s*boundary\\s*=(\\S+)", conf.regexp_option, nil, 1)
        if bfrom then
            boundary = string.sub(contentType, bfrom, bto)
        end
    end
    if boundary then -- form-data
        local sock, err = ngx.req.socket()
        if not sock or err then
            ngx.log(ngx.ERR, "sock nil...,err: ", err)
            return nil
        end
        ngx.req.init_body(128 * 1024)
        local size = 0
        local delimiter = '--' .. boundary
        local delimiterEnd = '--' .. boundary .. '--'
        local body = ''
        local isFile = false
        while size < contentLength do
            local line, err = sock:receive()
            if not line or err then
                break
            end

            if line == delimiter or line == delimiterEnd then
                if body ~= '' then
                    body = string.sub(body, 1, -2)
                    if isFile then
                        for _, rule in ipairs(rules) do
                            if _M.match(rule, { file_body = body }) then
                                log_and_exit(rule)
                                return 1024
                            end
                        end
                        isFile = false
                    else
                        for _, rule in ipairs(rules) do
                            if _M.match(rule, { form_body = body }) then
                                rule.request_data = body
                                log_and_exit(rule)
                                return 1024
                            end
                        end
                    end
                    body = ''
                end
            elseif line ~= '' then
                if isFile then
                    if body == '' then
                        local fr = ngx.re.find(line, "Content-Type:\\s*\\S+/\\S+", "ijo")
                        if not fr then
                            body = body .. line .. '\n'
                        end
                    else
                        body = body .. line .. '\n'
                    end
                else
                    local from, to = ngx.re.find(line,
                        [[Content-Disposition:\s*form-data;[\s\S]+filename=["|'][\s\S]+\.(\w+)(?:"|')]],
                        "ijo", nil, 1)

                    if from then
                        local suffix = string.sub(line, from, to)
                        for _, rule in ipairs(rules) do
                            if _M.match(rule, { suffix = suffix }) then
                                rule.request_data = line
                                log_and_exit(rule)
                                return 1024
                            end
                        end
                        isFile = true
                    else
                        local fr = ngx.re.find(line, "Content-Disposition:\\s*form-data;\\s*name=", "ijo")
                        if fr == nil then
                            body = body .. line .. '\n'
                        end
                    end
                end
            end
            size = size + string.len(line)
            ngx.req.append_body(line .. '\n')
        end
        ngx.req.finish_body()
    elseif ngx.re.find(contentType, "\\s*x-www-form-urlencoded", conf.regexp_option) then -- x-www-form-urlencoded
        ngx.req.read_body()
        local args = ngx.req.get_post_args()
        if not next(args) then
            return nil
        end

        for _, rule in ipairs(rules) do
            if _M.match(rule, args) then
                rule.request_data = args
                log_and_exit(rule)
                return 1024
            end
        end
    else -- raw
        ngx.req.read_body()
        local body = ngx.req.get_body_data()
        if not body then
            return nil
        end
        for _, rule in ipairs(rules) do
            if _M.match(rule, cjson.decode(body)) then
                rule.request_data = body
                log_and_exit(rule)
                return 1024
            end
        end
    end

    return nil
end


local rule_process = {
    [1] = request_args,
    [2] = unsafe_http_method,
    [3] = white_black_url,
    [4] = request_body,
    [5] = request_headers,
}


function _M.get_value_from_shared_dict(dict_name, key)
    local json_str = tostring(dict_name:get(key))

    if not json_str or json_str == "nil" or #json_str == 0 then
        return nil
    end

    local success, value = pcall(cjson.decode, json_str)

    if not success or type(value) ~= "table" or not next(value) then
        return nil
    end

    return value
end

local function process_ip_list(rule_action, ip_type)
    local ip_list = _M.get_value_from_shared_dict(conf.shared_dict_waf_ip, conf.db_table.ip)
    if not ip_list then
        return nil
    end

    for _, item in ipairs(ip_list) do
        if ngx.ctx.ip == item.ip_address and item.ip_type == ip_type then
            item.rule_action = rule_action
            log_and_exit(item)
            return 1024
        end
    end

    return nil
end

function _M.white_list()
    return process_ip_list(3, 1)
end

function _M.black_list()
    return process_ip_list(1, 2)
end

function _M.rule()
    local rules = _M.get_value_from_shared_dict(conf.shared_dict_waf_rule, conf.db_table.rule)
    if not rules then
        return
    end

    for _, rule in ipairs(rules) do
        local check_func = rule_process[rule.rule_variable]
        if check_func then
            if check_func(rule.value) then
                return
            end
        end
    end
end

return _M
