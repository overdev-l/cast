package config

import "strings"

var LiveId int
var PlatFrom string
var Url string // https://live.douyin.com/480552246257

func GetDouyinSocketUrl() string {
	liveid := Url[strings.LastIndex(Url, "/")+1:]
	return "wss://webcast5-ws-web-hl.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.0.14-beta.0&update_version_code=1.0.14-beta.0&compress=gzip&device_platform=web&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/130.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&cursor=r-1_d-1_u-1_fh-7436036470651474998_t-1731337266110&internal_ext=internal_src:dim|wss_push_room_id:" + liveid + "|wss_push_did:7419708364366611968|first_req_ms:1731337266033|fetch_time:1731337266110|seq:1|wss_info:0-1731337266110-0-0|wrds_v:7436036922931154702&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&endpoint=live_pc&support_wrds=1&user_unique_id=7419708364366611968&im_path=/webcast/im/fetch/&identity=audience&need_persist_msg_count=15&insert_task_id=&live_reason=&room_id=" + liveid + "&heartbeatDuration=0&signature=6B2uWMgrGlubFo/F"
}
