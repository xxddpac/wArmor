local geo = require "lib.maxminddb"
local conf = require "config"
local _M = {}

local pcall = pcall
local dbFile = conf.geoip_db.file

local hongkong = "香港"
local taiwan = "台湾"
local macao = "澳门"
local cn = "CN"
local china = "中国"

function _M.lookup(ip)
    local default_ip_info = {
        city = conf.unknown,
        country = conf.unknown,
        latitude = conf.unknown,
        longitude = conf.unknown,
        iso_code = conf.unknown
    }
    if not geo.initted() then
        geo.init(dbFile)
    end

    local ok, result, err = pcall(geo.lookup, ip)
    if not ok or not result then
        ngx.log(ngx.ERR, "Failed to lookup " .. ip .. ", err: " .. tostring(err))
        return default_ip_info
    end
    local language = conf.geoip_db.language
    local city = result.city and result.city.names and result.city.names[language] or default_ip_info.city
    local country = result.country and result.country.names and result.country.names[language] or default_ip_info
        .country
    local latitude = result.location and result.location.latitude or default_ip_info.latitude
    local longitude = result.location and result.location.longitude or default_ip_info.longitude
    local iso_code = result.country and result.country.iso_code or default_ip_info.iso_code
    if country == hongkong or country == taiwan or country == macao then
        city = country
        country = china
        iso_code = cn
    end
    return {
        country = country,
        city = city,
        longitude = longitude,
        latitude = latitude,
        iso_code = iso_code
    }
end

return _M
