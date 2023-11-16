local conf = require "config"
local dispatch = require "dispatch"
local jit = require "lib.utils"

local function get_client_ip()
    local ips = {
        ngx.var.http_x_forwarded_for or conf.unknown,
        ngx.var.http_proxy_client_ip or conf.unknown,
        ngx.var.http_wl_proxy_client_ip or conf.unknown,
        ngx.var.http_http_client_ip or conf.unknown,
        ngx.var.http_http_x_forwarded_for or conf.unknown,
        ngx.var.remote_addr or conf.unknown
    }
    for _, ip in ipairs(ips) do
        if type(ip) == "table" then
            ip = ip[1]
        end
        if ip and #ip ~= 0 and string.lower(ip) ~= conf.unknown then
            return ip
        end
    end
    return conf.unknown
end

local function check()
    if dispatch.white_list() then
    elseif dispatch.black_list() then
    elseif dispatch.rule() then
    end
end

local function init()
    local config = dispatch.get_value_from_shared_dict(conf.shared_dict_waf_config, conf.db_table.config)
    if not config then
        return
    end
    local mode = config[1].mode
    if mode == 3 then --旁路模式
        return
    end
    ngx.ctx.ip = get_client_ip()
    ngx.ctx.ua = ngx.var.http_user_agent or conf.unknown
    ngx.ctx.uuid = jit.uuid()
    ngx.ctx.now = ngx.localtime()
    ngx.ctx.mode = mode
    check()
end

init()

