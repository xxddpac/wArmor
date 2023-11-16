local conf = require "config"
local logger = require "resty.socket"
local cjson = require "cjson"
local geoip = require "geoip"
local initialized = false

_M = {}

function _M.send(msg)
    if not initialized then
        local ok, err = logger.init {
            host = conf.log.host,
            port = conf.log.port,
            flush_limit = conf.log.flush_limit,
            drop_limit = conf.log.drop_limit,
            pool_size = conf.log.pool_size,
            sock_type = conf.log.sock_type,
            timeout = conf.log.timeout
        }
        if not ok then
            ngx.log(ngx.ERR, "failed to initialize the logger: ", err)
            return
        end
        initialized = true
    end
    if #msg > conf.log.drop_limit then
        msg = string.sub(msg, 1, conf.log.drop_limit - 1)
    end
    msg = cjson.encode(msg) .. "\n"
    local _, err = logger.log(msg)
    if err then
        ngx.log(ngx.ERR, "Failed to log message: ", err)
    end
end

function _M.new_log(msg)
    if type(msg) ~= "table" then
        return nil
    end
    local geoip = geoip.lookup(ngx.ctx.ip)
    msg.city = geoip.city
    msg.country = geoip.country
    msg.latitude = geoip.latitude
    msg.longitude = geoip.longitude
    msg.iso_code = geoip.iso_code
    msg.request_server_name = ngx.var.host
    msg.request_uri = ngx.var.request_uri
    msg.request_ua = ngx.ctx.ua
    msg.timestamp = ngx.ctx.now
    msg.request_id = ngx.ctx.uuid
    msg.request_ip = ngx.ctx.ip
    msg.request_method = ngx.req.get_method()
    msg.mode = conf.desc.mode[ngx.ctx.mode]
    msg.program = conf.log.program -- 远程syslog-ng标识匹配存储至指定文件(/var/log/waf.log)
    return msg
end

return _M
