local jit = require "resty.jit-uuid"
local conf = require "config"
local cjson = require "cjson"

local _M = {}
jit.seed()


function _M.uuid()
    return jit.generate_v4()
end

function _M.read_html_file_to_string()
    local file = io.open(conf.redirect_html, "r")
    if not file then
        return nil
    end
    local text = file:read('*a')
    file:close()
    return text
end

function _M.format_value_from_rule_engine(key, result)
    local value
    if key == conf.db_table.rule then
        local mergedRules = {}

        for _, rule in ipairs(result) do
            local ruleVariable = rule.rule_variable
            if not mergedRules[ruleVariable] then
                mergedRules[ruleVariable] = { rule_variable = ruleVariable, value = {} }
            end
            table.insert(mergedRules[ruleVariable].value, rule)
        end

        local mergedRulesArray = {}
        for _, mergedRule in pairs(mergedRules) do
            table.insert(mergedRulesArray, mergedRule)
        end
        -- 通过排序保证白名单url规则优先检查
        local old_group = {}
        local new_group = {}
        for i = 1, #mergedRulesArray do
            local item = mergedRulesArray[i]
            if item.rule_variable == 3 and item.value[1].rule_type == 6 then
                table.insert(old_group, item)
            else
                table.insert(new_group, item)
            end
        end
        local sorted_result = {}
        for _, item in ipairs(old_group) do
            table.insert(sorted_result, item)
        end
        for _, item in ipairs(new_group) do
            table.insert(sorted_result, item)
        end
        value = sorted_result
    else
        value = result
    end
    conf["shared_dict_waf_" .. key]:set(key, cjson.encode(value))
end

return _M
