package db

import (
	"fmt"
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/datastruct/important"
	"shensuanzi/log"
	"shensuanzi/tools"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetDirectDownloadApp() string {
	engine := handle.mysqlEngine
	column_name := "down_load_u_i_addr"
	sql := "select down_load_u_i_addr from domain_info where id = ?"
	results, err := engine.Query(sql, datastruct.DefaultId)
	if err != nil || len(results) <= 0 {
		return ""
	}
	return string(results[0][column_name][:])
}

func (handle *DBHandler) IsExistPhone(phone string, isFT bool) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	var sql string
	if isFT {
		sql = "select 1 from cold_f_t_info where phone = ? limit 1"
	} else {
		sql = "select 1 from cold_user_info where phone = ? limit 1"
	}
	results, err := engine.Query(sql, phone)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(results)
	isExist := false
	if count > 0 {
		isExist = true
	}
	return isExist, datastruct.NULLError
}

func (handle *DBHandler) IsExistNickName(nickname string) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select 1 from cold_f_t_info where nick_name = ? limit 1"
	results, err := engine.Query(sql, nickname)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(results)
	isExist := false
	if count > 0 {
		isExist = true
	}
	return isExist, datastruct.NULLError
}
func (handle *DBHandler) GetFtMarkInfo() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select ftm.desc from f_t_mark_info ftm order by id desc"
	results, err := engine.Query(sql)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["desc"][:]))
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) FtRegister(body *datastruct.FTRegisterBody) datastruct.CodeType {
	return ftRegister(handle.mysqlEngine, body, nil)
}

func (handle *DBHandler) FtRegisterWithID(body *datastruct.FTRegisterWithIDBody) datastruct.CodeType {
	return ftRegister(handle.mysqlEngine, &body.FTRegisterBody, body)
}

func ftRegister(engine *xorm.Engine, body *datastruct.FTRegisterBody, IDbody *datastruct.FTRegisterWithIDBody) datastruct.CodeType {
	nowTime := time.Now().Unix()
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.CreatedAt = nowTime
	cold_ft.Introduction = body.Desc
	cold_ft.NickName = body.NickName
	cold_ft.Phone = body.Phone
	cold_ft.Pwd = body.Pwd
	cold_ft.Avatar = body.Avatar
	cold_ft.Registration = body.Platform
	cold_ft.AuthState = datastruct.Authing

	if IDbody != nil {
		cold_ft.IdentityCard = IDbody.Identity
		cold_ft.IdFrontCover = IDbody.IdFrontCover
		cold_ft.IdBehindCover = IDbody.IdBehindCover
		cold_ft.ActualName = IDbody.ActualName
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.InsertOne(cold_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister insert ColdFTInfo :%s", err.Error())
		rollbackError(str, session)
		return datastruct.NickNamePhoneIsUsed
	}

	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.FTId = cold_ft.Id
	hot_ft.Mark = body.Mark

	auth := new(datastruct.Authentication)
	auth.FTId = cold_ft.Id
	if IDbody != nil {
		auth.IdCardState = datastruct.IdCardSubmited
	} else {
		auth.IdCardState = datastruct.IdCardNotSubmit
	}
	auth.IsCP = false
	auth.IsHR = false

	//create shop
	shop := new(datastruct.ShopInfo)
	shop.FTId = cold_ft.Id

	_, err = session.Insert(hot_ft, auth, shop)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister insert HotFTInfo and Authentication :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func rollbackError(err_str string, session *xorm.Session) {
	log.Error("will rollback,err_str:%v", err_str)
	session.Rollback()
}

func (handle *DBHandler) FtLogin(body *datastruct.FtLoginBody) (*datastruct.RespFtLogin, datastruct.CodeType) {
	engine := handle.mysqlEngine
	cold_ft := new(datastruct.ColdFTInfo)
	has, err := engine.Where("phone=?", body.Phone).Get(cold_ft)
	if err != nil {
		log.Error("FtLogin get err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	if !has {
		return nil, datastruct.NotRegisterPhone
	}
	if cold_ft.Pwd != body.Pwd {
		return nil, datastruct.PwdError
	}
	var code datastruct.CodeType
	var rs *datastruct.RespFtLogin
	rs = nil
	switch cold_ft.AuthState {
	case datastruct.Authing:
		code = datastruct.AuthingCode
	case datastruct.AuthFailed:
		code = datastruct.AuthFailedCode
	default:
		rs, code = ftLoginSucceed(cold_ft.Id, engine)
	}
	return rs, code
}

func ftLoginSucceed(ft_id int, enging *xorm.Engine) (*datastruct.RespFtLogin, datastruct.CodeType) {
	token := commondata.UniqueId()
	nowTime := time.Now().Unix()
	im_id := fmt.Sprintf("ft_%d", ft_id)
	im_privatekey, code := tools.AccountGenForIM(im_id, important.IM_SDK_APPID)
	if code != datastruct.NULLError {
		return nil, code
	}
	sql := "update hot_f_t_info set token = ?,login_time = ?,i_m_private_key = ? where f_t_id=?"
	_, err := enging.Exec(sql, token, nowTime, im_privatekey, ft_id)
	if err != nil {
		log.Error("ftLoginSucceed err:%s", err.Error())
		return nil, datastruct.UpdateDataFailed
	}
	respFtInfo, code := getFtInfo(ft_id, enging)
	rs := new(datastruct.RespFtLogin)
	rs.FtInfo = respFtInfo
	rs.Token = token
	rs.IMPrivateKey = im_privatekey
	return rs, code
}

func (handle *DBHandler) GetFtInfo(ft_id int) (*datastruct.RespFtInfo, datastruct.CodeType) {
	return getFtInfo(ft_id, handle.mysqlEngine)
}

func getFtInfo(ft_id int, engine *xorm.Engine) (*datastruct.RespFtInfo, datastruct.CodeType) {
	sql := "select cold.avatar,cold.nick_name,hot.enable_free,hot.mark,auth.is_c_p,auth.is_h_r,auth.id_card_state from cold_f_t_info cold join hot_f_t_info hot on cold.id = hot.f_t_id join authentication auth on auth.f_t_id = cold.id where id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil {
		log.Error("GetFtInfo err0: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	if len(results) <= 0 {
		log.Error("GetFtInfo err1: len(results) is zero")
		return nil, datastruct.GetDataFailed
	}
	mp := results[0]
	info := new(datastruct.RespFtInfo)
	info.Id = ft_id
	info.Avatar = string(mp["avatar"][:])
	info.NickName = string(mp["nick_name"][:])
	info.Mark = string(mp["mark"][:])
	info.EnableFree = tools.StringToBool(string(mp["enable_free"][:]))
	authIcon := new(datastruct.AuthIcon)
	authIcon.IsCP = tools.StringToBool(string(mp["is_c_p"][:]))
	authIcon.IsHR = tools.StringToBool(string(mp["is_h_r"][:]))
	authIcon.IdCard = tools.StringToIdCardState(string(mp["id_card_state"][:]))
	info.AuthIcon = authIcon
	return info, datastruct.NULLError
}

func (handle *DBHandler) GetFtDataWithToken(token string) (*datastruct.FtRedisData, bool) {
	engine := handle.mysqlEngine
	sql := "select f_t_id,account_state from hot_f_t_info where token = ?"
	results, err := engine.Query(sql, token)
	if err != nil || len(results) <= 0 {
		return nil, false
	}
	redis_data := new(datastruct.FtRedisData)
	redis_data.Token = token

	mp := results[0]
	redis_data.FtId = tools.StringToInt(string(mp["f_t_id"][:]))
	redis_data.AccountState = datastruct.AccountState(tools.StringToInt(string(mp["f_t_id"][:])))
	sql = "update hot_f_t_info set login_time = ? where token = ?"
	_, err = engine.Exec(sql, time.Now().Unix(), token)
	if err != nil {
		log.Error("GetFtDataWithToken err:%s", err.Error())
		return nil, false
	}

	return redis_data, true
}

func (handle *DBHandler) UpdateFtInfo(body *datastruct.UpdateFtInfoBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	var err error
	var results []map[string][]byte
	sql := "select avatar from cold_f_t_info where id = ?"
	results, err = engine.Query(sql, ft_id)
	if err != nil {
		log.Error("DBHandler->UpdateFtInfo get avatar err:", err.Error())
		return datastruct.UpdateDataFailed
	}
	avatar := ""
	if len(results) > 0 {
		avatar = string(results[0]["avatar"][:])
	}
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.Avatar = body.Avatar
	cold_ft.NickName = body.NickName
	_, err = engine.Where("id=?", ft_id).Cols("nick_name", "avatar").Update(cold_ft)
	if err != nil {
		log.Error("DBHandler->UpdateFtInfo err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	if avatar != "" {
		commondata.DeleteOSSFileWithUrl(avatar)
	}
	return datastruct.NULLError
}

func (handle *DBHandler) UpdateFtMark(body *datastruct.UpdateFtMarkBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.Mark = body.Mark
	_, err := engine.Where("f_t_id=?", ft_id).Cols("mark").Update(hot_ft)
	if err != nil {
		log.Error("DBHandler->UpdateFtMark err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) UpdateFtIntroduction(body *datastruct.UpdateFtIntroductionBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.Introduction = body.Desc
	_, err := session.Where("id=?", ft_id).Cols("introduction").Update(cold_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction err0: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	shop_id, code := getShopIdWithSession(session, ft_id)
	if code != datastruct.NULLError {
		return code
	}

	sql := "select img_url from shop_imgs where shop_id = ?"
	results, err := session.Query(sql, shop_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction get imgs_url err")
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	will_DeleteImgs := make([]string, 0, len(results))
	for _, v := range results {
		will_DeleteImgs = append(will_DeleteImgs, string(v["img_url"][:]))
	}

	sql = "delete from shop_imgs where shop_id = ?"
	_, err = session.Exec(sql, shop_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction err1: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = "insert into shop_imgs(shop_id,img_url) values "
	values := ""
	tmp := ""
	for i := 0; i < len(body.Imgs); i++ {
		if i == 0 {
			tmp = fmt.Sprintf("(%d,'%s')", shop_id, body.Imgs[i])
		} else {
			tmp = fmt.Sprintf(",(%d,'%s')", shop_id, body.Imgs[i])
		}
		values += tmp
	}

	sql = sql + values
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction err2: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	if len(will_DeleteImgs) > 0 {
		for _, v := range will_DeleteImgs {
			commondata.DeleteOSSFileWithUrl(v)
		}
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetFtIntroduction(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select introduction from cold_f_t_info where id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtIntroduction err0")
		return nil, datastruct.GetDataFailed
	}
	desc := string(results[0]["introduction"][:])

	sql = "select id from shop_info where f_t_id = ?"
	results, err = engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtIntroduction err1")
		return nil, datastruct.GetDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))

	sql = "select img_url from shop_imgs where shop_id = ? order by id asc"
	results, err = engine.Query(sql, shop_id)
	if err != nil || len(results) < 0 {
		log.Error("DBHandler->GetFtIntroduction err2")
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["img_url"][:]))
	}
	resp := new(datastruct.UpdateFtIntroductionBody)
	resp.Desc = desc
	resp.Imgs = arr
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetFtAutoReply(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select auto_reply from hot_f_t_info where f_t_id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtAutoReply err0")
		return nil, datastruct.GetDataFailed
	}
	autoreply := string(results[0]["auto_reply"][:])

	sql = "select reply from f_t_quick_reply where f_t_id = ? order by id asc"
	results, err = engine.Query(sql, ft_id)
	if err != nil || len(results) < 0 {
		log.Error("DBHandler->GetFtIntroduction err1")
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["reply"][:]))
	}
	resp := new(datastruct.UpdateFtAutoReplyBody)
	resp.AutoReply = autoreply
	resp.QuickReply = arr
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateFtAutoReply(body *datastruct.UpdateFtAutoReplyBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.AutoReply = body.AutoReply
	_, err := session.Where("f_t_id=?", ft_id).Cols("auto_reply").Update(hot_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtAutoReply err0: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql := "delete from f_t_quick_reply where f_t_id = ?"
	_, err = session.Exec(sql, ft_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtAutoReply err1: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = "insert into f_t_quick_reply(f_t_id,reply) values "
	values := ""
	tmp := ""
	for i := 0; i < len(body.QuickReply); i++ {
		if i == 0 {
			tmp = fmt.Sprintf("(%d,'%s')", ft_id, body.QuickReply[i])
		} else {
			tmp = fmt.Sprintf(",(%d,'%s')", ft_id, body.QuickReply[i])
		}
		values += tmp
	}

	sql = sql + values
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtAutoReply err2: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) FtSubmitIdentity(body *datastruct.FtIdentity, ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine

	var err error
	var results []map[string][]byte
	sql := "select id_front_cover,id_behind_cover from cold_f_t_info where id = ?"
	results, err = engine.Query(sql, ft_id)
	if err != nil {
		log.Error("DBHandler->FtSubmitIdentity get id_cover err:", err.Error())
		return nil, datastruct.UpdateDataFailed
	}
	last_front_cover := ""
	last_behind_cover := ""
	if len(results) > 0 {
		last_front_cover = string(results[0]["id_front_cover"][:])
		last_behind_cover = string(results[0]["id_behind_cover"][:])
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	ft_cold := new(datastruct.ColdFTInfo)
	ft_cold.IdentityCard = body.Identity
	ft_cold.ActualName = body.ActualName
	ft_cold.IdBehindCover = body.IdBehindCover
	ft_cold.IdFrontCover = body.IdFrontCover

	_, err = session.Where("id=?", ft_id).Cols("identity_card", "actual_name", "id_behind_cover", "id_front_cover").Update(ft_cold)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtSubmitIdentity err0: %s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}

	auth := new(datastruct.Authentication)
	auth.IdCardState = datastruct.IdCardSubmited
	_, err = session.Where("f_t_id=?", ft_id).Cols("id_card_state").Update(auth)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtSubmitIdentity err1: %s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtSubmitIdentity Commit :%s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}
	if last_front_cover != "" || last_behind_cover != "" {
		commondata.DeleteOSSFileWithUrl(last_front_cover)
		commondata.DeleteOSSFileWithUrl(last_behind_cover)
	}
	return auth.IdCardState, datastruct.NULLError
}

func (handle *DBHandler) GetAppraised(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select ap.user_id,ap.score,ap.mark,ap.appraised_type,ap.desc,ap.is_anonym,ap.created_at,pr.product_name from appraised_info ap join product_info pr on ap.product_id = pr.id join shop_info sh on sh.id = pr.shop_id where sh.f_t_id = ? LIMIT ?,?"
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetAppraised err0: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	rs := make([]*datastruct.RespAppraise, 0, len(results))
	for _, v := range results {
		ap := new(datastruct.RespAppraise)
		is_anonym := tools.StringToBool(string(v["is_anonym"][:]))
		if is_anonym {
			user_id := tools.StringToInt64(string(v["user_id"][:]))
			sql = "select nick_name,avatar from cold_user_info where id = ?"
			sub_results, sub_err := engine.Query(sql, user_id)
			if sub_err != nil || len(sub_results) <= 0 {
				log.Error("DBHandler->GetAppraised err1: the user is not exist, id is %v", user_id)
				return nil, datastruct.GetDataFailed
			}
			ap.Name = string(sub_results[0]["nick_name"][:])
			ap.Avatar = string(sub_results[0]["avatar"][:])
		}
		ap.IsAnonym = is_anonym
		ap.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		ap.Desc = string(v["desc"][:])
		ap.Mark = string(v["mark"][:])
		ap.ProductName = string(v["product_name"][:])
		ap.Score = tools.StringToFloat64(string(v["score"][:]))
		rs = append(rs, ap)
	}
	return rs, datastruct.NULLError
}

func (handle *DBHandler) GetFtUnReadMsgCount(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	where_str := "f_t_read = 0 and f_t_id = ?"
	arr := make([]interface{}, 0)
	arr = append(arr, new(datastruct.OrderInfoMsg))
	arr = append(arr, new(datastruct.OrderRightsFinishedMsg))
	arr = append(arr, new(datastruct.FTRegisterMsg))
	arr = append(arr, new(datastruct.FTOrderRefundMsg))
	var count int64
	count = 0
	for _, v := range arr {
		tmp, err := engine.Where(where_str, ft_id).Count(v)
		if err != nil {
			log.Error("DBHandler->GetFtUnReadMsgCount err:", err.Error())
			return nil, datastruct.GetDataFailed
		}
		count += tmp
	}
	return count, datastruct.NULLError
}

func (handle *DBHandler) GetUserUnReadMsgCount(user_id int) (interface{}, datastruct.CodeType) {

	engine := handle.mysqlEngine
	where_str := "user_read = 0 and user_id = ?"
	arr := make([]interface{}, 0)
	arr = append(arr, new(datastruct.OrderInfoMsg))
	arr = append(arr, new(datastruct.OrderRightsFinishedMsg))
	arr = append(arr, new(datastruct.UserRegisterMsg))
	arr = append(arr, new(datastruct.UserOrderRefundMsg))

	var count int64
	count = 0
	for _, v := range arr {
		tmp, err := engine.Where(where_str, user_id).Count(v)
		if err != nil {
			log.Error("DBHandler->GetUserUnReadMsgCount err:", err.Error())
			return nil, datastruct.GetDataFailed
		}
		count += tmp
	}
	return count, datastruct.NULLError
}

func (handle *DBHandler) GetFtSystemMsg(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	//type 1 f_t_register_msg
	//type 2 f_t_order_refund_msg
	//type 3 order_info_msg
	//type 4 order_rights_finished_msg
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	sql := "select * from (SELECT -1 as id,ftr.f_t_read,ftr.created_at,ftr.f_t_id,1 as type,-1 as user_nick_name,-1 as product_name,-1 as handle from f_t_register_msg ftr union all select ftor.id,ftor.f_t_read,ftor.created_at,ftor.f_t_id,2 as type,user_nick_name,product_name,refund_type as handle from f_t_order_refund_msg ftor union all select oi.id,oi.f_t_read,oi.created_at,oi.f_t_id,3 as type,user_nick_name,product_name,-1 as handle from order_info_msg oi union all select orf.id,orf.f_t_read,orf.created_at,orf.f_t_id,4 as type,user_nick_name,product_name,is_agree as handle from order_rights_finished_msg orf) as tmp_msg where f_t_id = ? ORDER BY created_at DESC LIMIT ?,?"
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetFtSystemMsg err0: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]interface{}, 0, len(results))
	var rs interface{}
	for _, v := range results {
		tableType := tools.StringToInt(string(v["type"][:]))
		isRead := tools.StringToBool(string(v["f_t_read"][:]))
		switch tableType {
		case 1:
			rrm := new(datastruct.RespRegisterMsg)
			rrm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rrm.Type = tableType
			rrm.GZH_Name = important.WX_GZH_NAME
			rrm.QRCode = important.WX_GZH_QRCode
			rs = rrm
			if !isRead {
				ftrm := new(datastruct.FTRegisterMsg)
				ftrm.FTRead = true
				engine.Where("f_t_id=?", ft_id).Cols("f_t_read").Update(ftrm)
			}
		case 2:
			rrfm := new(datastruct.RespRefundFTMsg)
			rrfm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rrfm.OrderId = tools.StringToInt64(string(v["id"][:]))
			rrfm.NickName = string(v["user_nick_name"][:])
			rrfm.ProductName = string(v["product_name"][:])
			rrfm.RefundType = datastruct.UserOrderRefundType(tools.StringToInt(string(v["handle"][:])))
			rrfm.Type = tableType
			rs = rrfm
			if !isRead {
				ftor := new(datastruct.FTOrderRefundMsg)
				ftor.FTRead = true
				engine.Where("id=?", rrfm.OrderId).Cols("f_t_read").Update(ftor)
			}
		case 3:
			toi := new(datastruct.TmpOrderInfoFT)
			toi.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			toi.OrderId = tools.StringToInt64(string(v["id"][:]))
			toi.NickName = string(v["user_nick_name"][:])
			toi.ProductName = string(v["product_name"][:])
			toi.Type = tableType
			rs = toi
			if !isRead {
				oi := new(datastruct.OrderInfoMsg)
				oi.FTRead = true
				engine.Where("id=?", toi.OrderId).Cols("f_t_read").Update(oi)
			}
		case 4:
			rfm := new(datastruct.RespRightsFinishedFTMsg)
			rfm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rfm.OrderId = tools.StringToInt64(string(v["id"][:]))
			rfm.NickName = string(v["user_nick_name"][:])
			rfm.ProductName = string(v["product_name"][:])
			rfm.Type = tableType
			rfm.IsAgree = tools.StringToBool(string(v["handle"][:]))
			rs = rfm
			if !isRead {
				orf := new(datastruct.OrderRightsFinishedMsg)
				orf.FTRead = true
				engine.Where("id=?", rfm.OrderId).Cols("f_t_read").Update(orf)
			}
		}
		arr = append(arr, rs)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) GetUserSystemMsg(user_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	//type 1 register_msg
	//type 2 user_order_refund_msg
	//type 3 order_info_msg
	//type 4 order_rights_finished_msg
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	sql := "select * from (select -1 as id,ur.user_read,ur.created_at,ur.user_id,1 as type,-1 as f_t_nick_name,-1 as product_name,-1 as handle from user_register_msg ur union all select uor.id,uor.user_read,uor.created_at,uor.user_id,2 as type,f_t_nick_name,product_name,refund_result_type as handle from user_order_refund_msg uor union all select oi.id,oi.user_read,oi.created_at,oi.user_id,3 as type,f_t_nick_name,product_name,-1 as handle from order_info_msg oi union all select orf.id,orf.user_read,orf.created_at,orf.user_id,4 as type,f_t_nick_name,product_name,is_agree as handle from order_rights_finished_msg orf) as tmp_msg where user_id = ? ORDER BY created_at DESC LIMIT ?,?"
	results, err := engine.Query(sql, user_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetUserSystemMsg err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]interface{}, 0, len(results))
	var rs interface{}
	for _, v := range results {
		tableType := tools.StringToInt(string(v["type"][:]))
		isRead := tools.StringToBool(string(v["user_read"][:]))
		switch tableType {
		case 1:
			rrm := new(datastruct.RespRegisterMsg)
			rrm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rrm.Type = tableType
			rrm.GZH_Name = important.WX_GZH_NAME
			rrm.QRCode = important.WX_GZH_QRCode
			rs = rrm
			if !isRead {
				urm := new(datastruct.UserRegisterMsg)
				urm.UserRead = true
				engine.Where("user_id=?", user_id).Cols("user_read").Update(urm)
			}
		case 2:
			rrfu := new(datastruct.RespRefundUserMsg)
			rrfu.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rrfu.OrderId = tools.StringToInt64(string(v["id"][:]))
			rrfu.NickName = string(v["f_t_nick_name"][:])
			rrfu.ProductName = string(v["product_name"][:])
			rrfu.RefundResultType = datastruct.UserOrderRefundResultType(tools.StringToInt(string(v["handle"][:])))
			rrfu.Type = tableType
			rs = rrfu
			if !isRead {
				uor := new(datastruct.UserOrderRefundMsg)
				uor.UserRead = true
				engine.Where("id=?", rrfu.OrderId).Cols("user_read").Update(uor)
			}
		case 3:
			toi := new(datastruct.TmpOrderInfoFT)
			toi.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			toi.OrderId = tools.StringToInt64(string(v["id"][:]))
			toi.NickName = string(v["f_t_nick_name"][:])
			toi.ProductName = string(v["product_name"][:])
			toi.Type = tableType
			rs = toi
			if !isRead {
				oi := new(datastruct.OrderInfoMsg)
				oi.UserRead = true
				engine.Where("id=?", toi.OrderId).Cols("user_read").Update(oi)
			}
		case 4:
			rfm := new(datastruct.RespRightsFinishedFTMsg)
			rfm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rfm.OrderId = tools.StringToInt64(string(v["id"][:]))
			rfm.NickName = string(v["f_t_nick_name"][:])
			rfm.ProductName = string(v["product_name"][:])
			rfm.Type = tableType
			rfm.IsAgree = tools.StringToBool(string(v["handle"][:]))
			rs = rfm
			if !isRead {
				orf := new(datastruct.OrderRightsFinishedMsg)
				orf.UserRead = true
				engine.Where("id=?", rfm.OrderId).Cols("user_read").Update(orf)
			}
		}
		arr = append(arr, rs)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) GetFtDndList(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select dnd.id,dnd.created_at,u.avatar,u.nick_name from f_t_dnd_list dnd join cold_user_info u on dnd.user_id = u.id where dnd.f_t_id = ? order by dnd.created_at desc LIMIT ?,?"
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetFtDndList err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	rs := make([]*datastruct.RespDndList, 0, len(results))
	for _, v := range results {
		dnd := new(datastruct.RespDndList)
		dnd.Id = tools.StringToInt(string(v["id"][:]))
		dnd.Avatar = string(v["avatar"][:])
		dnd.CreatedAt = tools.StringToInt64(string(v["avatar"][:]))
		dnd.NickName = string(v["nick_name"][:])
		rs = append(rs, dnd)
	}
	return rs, datastruct.NULLError
}

func (handle *DBHandler) RemoveFtDndList(id int, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	_, err := engine.Where("id=? and f_t_id=?", id, ft_id).Delete(new(datastruct.FTDndList))
	if err != nil {
		log.Error("DBHandler->RemoveFtDndList err: %s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func getShopIdWithSession(session *xorm.Session, ft_id int) (int, datastruct.CodeType) {
	sql := "select id from shop_info where f_t_id = ?"
	results, err := session.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		str := fmt.Sprintf("DBHandler->getShopIdWithSession err")
		rollbackError(str, session)
		return -1, datastruct.UpdateDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))
	return shop_id, datastruct.NULLError
}

func getShopId(engine *xorm.Engine, ft_id int) (int, datastruct.CodeType) {
	sql := "select id from shop_info where f_t_id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->getShopId err")
		return -1, datastruct.UpdateDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))
	return shop_id, datastruct.NULLError
}

func (handle *DBHandler) EditProduct(body *datastruct.EditProductBody, ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine

	sql := "select word from sensitive_word order by id asc"
	results, err := engine.Query(sql)
	if err != nil {
		log.Error("EditProduct get words err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}

	for _, v := range results {
		tmp := string(v["word"][:])
		if strings.Contains(body.ProductName, tmp) || strings.Contains(body.ProductDesc, tmp) {
			return tmp, datastruct.Sensitive
		}
	}

	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}

	count, err := engine.Where("shop_id=?", shop_id).Count(new(datastruct.ProductInfo))
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	if count >= datastruct.MAX_PRODUCT_COUNT {
		return nil, datastruct.MaxCreateCount
	}

	isAdd := true
	if body.Id > 0 {
		isAdd = false
	}
	nowTime := time.Now().Unix()
	pro := new(datastruct.ProductInfo)
	pro.UpdatedAt = nowTime
	pro.IsHidden = body.IsHidden
	pro.Price = body.Price
	pro.ProductDesc = body.ProductDesc
	pro.ProductName = body.ProductName
	//pro.SortId = 0
	if isAdd {
		pro.CreatedAt = nowTime
		pro.ShopId = shop_id
		_, err := engine.InsertOne(pro)
		if err != nil {
			log.Error("DBHandler->EditProduct insert err:%v", err.Error())
			return nil, datastruct.UpdateDataFailed
		}

	} else {

		_, err := engine.Where("id = ? and shop_id = ?", body.Id, shop_id).Cols("product_name", "product_desc", "price", "is_hidden", "updated_at").Update(pro)
		if err != nil {
			log.Error("DBHandler->EditProduct update err:%v", err.Error())
			return nil, datastruct.UpdateDataFailed
		}
	}
	return nil, datastruct.NULLError
}

func (handle *DBHandler) RemoveProduct(id int, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return code
	}
	_, err := engine.Where("id=? and shop_id=?", id, shop_id).Delete(new(datastruct.ProductInfo))
	if err != nil {
		log.Error("DBHandler->RemoveProduct err: %s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetProduct(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}
	resp := new(datastruct.RespProductInfo)
	onsale_count, onsale, code := getSaleInfo(shop_id, true, engine)
	if code != datastruct.NULLError {
		return nil, code
	}
	offsale_count, offsale, code := getSaleInfo(shop_id, false, engine)
	if code != datastruct.NULLError {
		return nil, code
	}
	resp.OnSaleCount = onsale_count
	resp.OnSale = onsale
	resp.OffSaleCount = offsale_count
	resp.OffSale = offsale
	resp.Total = onsale_count + offsale_count
	return resp, datastruct.NULLError
}

func getSaleInfo(shop_id int, isOnSale bool, engine *xorm.Engine) (int, []*datastruct.TmpProductInfo, datastruct.CodeType) {
	isHidden := 1
	if isOnSale {
		isHidden = 0
	}
	sql := "select id,product_name,product_desc,price from product_info where shop_id = ? and is_hidden = ? order by sort_id asc,created_at asc"
	results, err := engine.Query(sql, shop_id, isHidden)
	if err != nil {
		log.Error("DBHandler->GetProduct err0: %s", err.Error())
		return -1, nil, datastruct.GetDataFailed
	}
	sale_count := len(results)
	sale := make([]*datastruct.TmpProductInfo, 0, sale_count)
	for _, v := range results {
		tmp := new(datastruct.TmpProductInfo)
		tmp.Id = tools.StringToInt(string(v["id"][:]))
		tmp.Price = tools.StringToFloat64(string(v["price"][:]))
		tmp.ProductDesc = string(v["product_desc"][:])
		tmp.ProductName = string(v["product_name"][:])
		sale = append(sale, tmp)
	}
	return sale_count, sale, datastruct.NULLError

}
