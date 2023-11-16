local conf = require "config"
local log = require "lib.log"
local dispatch = require "dispatch"

if ngx.status ~= 403 and ngx.status ~= 502 then
    local response = ngx.arg[1]
    if not response then
        return
    end
    local rules = dispatch.get_value_from_shared_dict(conf.shared_dict_waf_rule, conf.db_table.rule)
    if not rules then
        return
    end
    local sensitive_rule
    for _, rule in ipairs(rules) do
        if rule.rule_variable == 6 then
            sensitive_rule = rule.value
        end
    end
    if not sensitive_rule then
        return
    end
    for _, rule in ipairs(sensitive_rule) do
        if dispatch.match(rule, { response = response }) then
            local msg = {
                rule_type = conf.desc.rule_type[rule.rule_type],
                rule_id = rule.id,
                response = response,
                severity = conf.desc.severity[rule.severity],
                program = conf.log.program,
                request_server_name = ngx.var.host,
                request_uri = ngx.var.request_uri,
                timestamp = ngx.ctx.now,
                request_method = ngx.req.get_method()
            }
            log.send(msg)
            return
        end
    end
end
