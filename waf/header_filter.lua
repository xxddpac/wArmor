local conf = require "config"

if ngx.status == 403 or ngx.status == 502 then
    ngx.header.server = conf.log.program
else
    ngx.header.content_length = nil
end
