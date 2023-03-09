package sns

import "github.com/zedisdog/ty/sdk/wechat/official/response"

// {
// 	"openid": "OPENID",
// 	"nickname": NICKNAME,
// 	"sex": 1,
// 	"province":"PROVINCE",
// 	"city":"CITY",
// 	"country":"COUNTRY",
// 	"headimgurl":"https://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/46",
// 	"privilege":[ "PRIVILEGE1" "PRIVILEGE2"     ],
// 	"unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL"
//   }

type UserInfoRes struct {
	response.Error
	OpenID    string   `json:"openid"`
	Nickname  string   `json:"nickname"`
	Sex       int      `json:"sex"`
	Avatar    string   `json:"headimgurl"`
	Privilege []string `json:"privilege"`
	UnionID   string   `json:"unionid"`
}
