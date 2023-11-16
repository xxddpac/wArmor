local redis = require "lib.redis"
local mysql = require "lib.mysql"
local executed = false

local function timer_callback(premature)
    if not premature then
        if not executed then
            ngx.thread.spawn(redis.subscribe)
            ngx.thread.spawn(function()
                mysql.sync()
            end)
            executed = true
        end
    end
end

if ngx.worker.id() == 0 then
    ngx.timer.at(0, timer_callback)
end


