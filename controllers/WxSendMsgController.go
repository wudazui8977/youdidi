package controllers

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"youdidi/models"
	"youdidi/redisClient"
)

type WxSendMsgController struct {
	beego.Controller
}

var (
	AccessTokenKey string
	Appid string
	AppSecret string
	TemplateMap map[int]string
	ItemMap map[int]int
)

type WxResult struct {
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
	Msgid int `json:"msgid"`
}

type Item struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type Items3 struct {
	First *Item `json:"first"`
	Keyword1 *Item `json:"keyword1"`
	Keyword2 *Item `json:"keyword2"`
	Keyword3 *Item `json:"keyword3"`
	Remark *Item `json:"remark"`
}

type Items4 struct {
	First *Item `json:"first"`
	Keyword1 *Item `json:"keyword1"`
	Keyword2 *Item `json:"keyword2"`
	Keyword3 *Item `json:"keyword3"`
	Keyword4 *Item `json:"keyword4"`
	Remark *Item `json:"remark"`
}

type Items5 struct {
	First *Item `json:"first"`
	Keyword1 *Item `json:"keyword1"`
	Keyword2 *Item `json:"keyword2"`
	Keyword3 *Item `json:"keyword3"`
	Keyword4 *Item `json:"keyword4"`
	Keyword5 *Item `json:"keyword5"`
	Remark *Item `json:"remark"`
}

type MsgTemplateItems3 struct {
	Touser string `json:"touser"`
	Template_id string `json:"template_id"`
	Url string `json:"url"`
	Topcolor string `json:"topcolor"`
	Data *Items3 `json:"data"`
}

type MsgTemplateItems4 struct {
	Touser string `json:"touser"`
	Template_id string `json:"template_id"`
	Url string `json:"url"`
	Topcolor string `json:"topcolor"`
	Data *Items4 `json:"data"`
}

type MsgTemplateItems5 struct {
	Touser string `json:"touser"`
	Template_id string `json:"template_id"`
	Url string `json:"url"`
	Topcolor string `json:"topcolor"`
	Data *Items5 `json:"data"`
}

func InitData () {
	Appid =  beego.AppConfig.String("weixin::AppId")
	AppSecret = beego.AppConfig.String("weixin::AppSecret")
	AccessTokenKey = "APP_ACCESS_TOKEN"

	TemplateMap = make(map[int]string)
	ItemMap = make(map[int]int)

	TemplateMap[0] = "SdyfFNTfGFGpxLXz9Es4xeNpToywPMsBC4r1NsftPg4" //乘客发起乘车申请推送 [接单成功通知]
	TemplateMap[1] = "WVpWH0teca_PwxBjiq7Im_hIcjIjsFC0MrH_gFVNb5Q" //车主拒绝请求推送 [拼车拒绝通知]
	TemplateMap[2] = "WVpWH0teca_PwxBjiq7Im_hIcjIjsFC0MrH_gFVNb5Q" //车主同意请求推送 [拼车拒绝通知]
	TemplateMap[3] = "NOFOaBaJYYZGa2feSuMWjmDeZR6Zqzg9DHV6N4jyRCs" //充值结果推送 [充值成功提醒]
	TemplateMap[4] = "GVZfbzaycoJdmB-zvjFDW1BW30b2ajn4lzoPdEIHhB8" //乘客取消推送 [订单取消提醒]

	ItemMap[0] = 5//乘客发起乘车申请推送
	ItemMap[1] = 3 //车主拒绝请求推送
	ItemMap[2] = 3 //车主同意请求推送
	ItemMap[3] = 3//充值结果推送
	ItemMap[4] = 5//乘客取消推送
}

func SendMsg (uid int, templateId int, data string, url string) bool{
	InitData()
	wxUrl := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="
	reqBody := ""
	accessToken := GetAccessToken()
	if (accessToken == "nil") {
		logs.Emergency("get access token fail")
		return false
	}
	wxUrl = wxUrl + accessToken

	openId, err := GetUserOpenId(uid)

	if (err != nil) {
		logs.Emergency("get user openId fail uid=%v err=%v", uid, err.Error())
		return false
	}

	itemNum := ItemMap[templateId]

	if (itemNum == 3) {
		itemBody := &Items3{}
		err1 := json.Unmarshal([]byte(data), &itemBody)
		if (err1 != nil) {
			logs.Error("pares data to item5 fail data=%v", data)
			return false
		}

		msgBody := &MsgTemplateItems3{}
		msgBody.Touser = openId
		msgBody.Template_id = TemplateMap[templateId]
		msgBody.Topcolor = "#FF0000"
		msgBody.Url = url
		msgBody.Data = itemBody

		reqBodyTmp, err2 := json.Marshal(&msgBody)
		reqBody = string(reqBodyTmp)

		if (err2 != nil) {
			logs.Error("pares struct to json fail err=%v", err2.Error())
			return false
		}
		logs.Debug("request body=%v", reqBody)
	} else if (itemNum ==5) {
		itemBody := &Items5{}
		err1 := json.Unmarshal([]byte(data), &itemBody)
		if (err1 != nil) {
			logs.Error("pares data to item5 fail data=%v", data)
			return false
		}

		msgBody := &MsgTemplateItems5{}
		msgBody.Touser = openId
		msgBody.Template_id = TemplateMap[templateId]
		msgBody.Topcolor = "#FF0000"
		msgBody.Url = url
		msgBody.Data = itemBody

		reqBodyTmp, err2 := json.Marshal(&msgBody)
		reqBody = string(reqBodyTmp)

		if (err2 != nil) {
			logs.Error("pares struct to json fail err=%v", err2.Error())
			return false
		}
		logs.Debug("request body=%v", reqBody)
	}


	req := httplib.Post(wxUrl)
	logs.Debug("wxurl=%v", wxUrl)

	req.Body(reqBody)
	result, errreq := req.String()

	if (errreq != nil) {
		logs.Error("http request fail err=%v", errreq.Error())
		return false
	}

	var wxRe WxResult
	errwxre := json.Unmarshal([]byte(result), &wxRe)

	if (errwxre != nil) {
		logs.Error("parse wx result to struct fail err=%v", errwxre.Error())
		return false
	}

	if (wxRe.Errcode != 0) {
		logs.Error("wx return err code =%v", wxRe.Errcode)
		return false
	}

	return true
}

func GetAccessToken () string {
	return redisClient.GetKey(AccessTokenKey)
}

func GetUserOpenId (uid int) (string, error) {
	var dbUser models.User
	var userInfo []*models.User

	num, err := orm.NewOrm().QueryTable(dbUser).Filter("Id", uid).All(&userInfo)

	if (num < 1 || err != nil) {
		logs.Error("get user info fail uid=%v", uid)
		return "", errors.New("can not get user info")
	}

	if (userInfo[0].OpenId == "") {
		logs.Error("openid of uid=%v is empty", uid)
		return "", errors.New("openid is empty")
	}

	return userInfo[0].OpenId, nil
}