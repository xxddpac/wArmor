local mysql = require "resty.mysql"
local conf = require "config"
local utils = require "lib.utils"

local _M = {}

function _M.conn()
    local db, err = mysql:new()
    if not db or err then
        ngx.log(ngx.ERR, "failed to instantiate mysql: ", err)
        return nil
    end
    db:set_timeout(conf.db.set_timeout)
    local ok, err, errcode, sqlstate = db:connect {
        host = conf.db.host,
        port = conf.db.port,
        database = conf.db.database,
        user = conf.db.username,
        password = conf.db.password,
        charset = conf.db.charset,
        max_packet_size = conf.db.max_packet_size,
    }

    if not ok or err then
        ngx.log(ngx.ERR, "failed to connect: ", err, ": ", errcode, " ", sqlstate)
        return nil
    end
    return db
end

function _M.sync(...)
    local db = _M.conn()
    if not db then
        return
    end
    local rule = conf.db_table.rule
    local config = conf.db_table.config
    local ip = conf.db_table.ip
    local args = { ... }
    local sql_list = {
        rule = string.format(
            [[
    SELECT
        rule.id,
        severity,
        rule_variable,
        rule_type,
        rules_operation,
        JSON_ARRAYAGG(rules) AS rules,
        rule_action
    FROM
        %s
    INNER JOIN
        %s
    ON
        rule.id = rules.rule_id
    WHERE
        status = 1
    GROUP BY
        rule.id, severity, rule_variable, rule_type, rules_operation, rule_action
    ORDER BY
        %s.updated_at DESC;
    ]], rule, rule .. "s", rule),
        config = string.format("SELECT mode FROM %s;", config),
        ip = string.format("SELECT id,ip_address,ip_type FROM %s;", ip),
    }

    local function execute_query(key, sql)
        local result, err, errcode, sqlstate = db:query(sql)
        if not result then
            ngx.log(ngx.ERR, "bad result: ", err, ": ", errcode, ": ", sqlstate, ".")
            return
        end
        utils.format_value_from_rule_engine(key, result)
    end

    if #args == 0 then
        for key, sql in pairs(sql_list) do
            execute_query(key, sql)
        end
    elseif args[1] == rule then
        execute_query(args[1], sql_list.rule)
    elseif args[1] == config then
        execute_query(args[1], sql_list.config)
    elseif args[1] == ip then
        execute_query(args[1], sql_list.ip)
    end
    _M.pool(db)
end

function _M.pool(db)
    local ok, err = db:set_keepalive(10000, 100)
    if not ok or err then
        ngx.log(ngx.ERR, "failed to set keepalive: ", err)
        return
    end
end

return _M
