local redis = require "resty.redis"
local cjson = require "cjson"
local mysql = require "lib.mysql"
local conf = require "config"

local _M = {}


function _M.redis_conn()
    local red, err = redis:new()
    if err or not red then
        ngx.log(ngx.ERR, "failed to new redis:", err)
        return nil
    end
    red:set_timeouts(conf.redis.connect_timeout, conf.redis.set_timeout, conf.redis.read_timeout)
    local ok, err = red:connect(conf.redis.host, conf.redis.port,
        { ssl = conf.redis.redis_ssl, pool_size = conf.redis.pool_size })
    if err or not ok then
        ngx.log(ngx.ERR, "failed to connect redis: ", err)
        return nil
    end
    if #conf.redis.password ~= 0 then
        local times = 0
        times, err = red:get_reused_times()
        if err then
            ngx.log(ngx.ERR, "failed get_reused_times: ", err)
            return nil
        end
        if times == 0 then
            local res, err = red:auth(conf.redis.password)
            if not res or err then
                ngx.log(ngx.ERR, "failed to authenticate: ", err)
                return nil
            end
        end
    end

    return red
end

function _M.subscribe()
    local red = _M.redis_conn()
    if not red then
        return
    end
    local res, err = red:subscribe(conf.redis.channel)
    if not res or err then
        ngx.log(ngx.ERR, "failed subscribe: ", err)
        return
    end
    while true do
        local msg, err = red:read_reply()
        if err then
            ngx.log(ngx.ERR, "failed read_reply: ", err)
            break
        end
        local result = cjson.decode(msg[3])
        if result.Event == "heartbeat" then
        else
            mysql.sync(result.Event)
        end
    end
    red:unsubscribe(_M.channel)
    _M.pool(red)
end

function _M.pool(red)
    local ok, err = red:set_keepalive(10000, 100)
    if not ok then
        ngx.log(ngx.ERR, "failed to set keepalive: ", err)
        return
    end
end

return _M
