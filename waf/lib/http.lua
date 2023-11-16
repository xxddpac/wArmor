local http = require "resty.http"
local conf = require "config"

local _M = {}

local function new()
    local httpc = http.new()
    if not httpc then
        return nil
    end
    httpc:set_keepalive(conf.rule_engine.max_idle_timeout, conf.rule_engine.pool_size)
    return httpc
end

function _M.send(url, method, body)
    local httpc = new()
    if not httpc then
        return
    end
    local res, err = httpc:request_uri(url, {
        method = method,
        headers = {
            ["Content-Type"] = "application/json",
        },
        body = body,
    })
    httpc:set_keepalive()
    if not res or res.status ~= 200 then
        ngx.log(ngx.ERR, "Request failed: ", err)
    end
end

return _M
