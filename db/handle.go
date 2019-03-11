package db

import (
	"fmt"
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/datastruct/important"
	"shensuanzi/log"
	"shensuanzi/tools"
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
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.Avatar = body.Avatar
	cold_ft.NickName = body.NickName
	_, err := engine.Where("id=?", ft_id).Cols("nick_name", "avatar").Update(cold_ft)
	if err != nil {
		log.Error("DBHandler->UpdateFtInfo err:%s", err.Error())
		return datastruct.UpdateDataFailed
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

	//get shop_id
	sql := "select id from shop_info where f_t_id = ?"
	results, err := session.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction get shop_id err")
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))

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
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	ft_cold := new(datastruct.ColdFTInfo)
	ft_cold.IdentityCard = body.Identity
	ft_cold.ActualName = body.ActualName
	ft_cold.IdBehindCover = body.IdBehindCover
	ft_cold.IdFrontCover = body.IdFrontCover
	_, err := session.Where("id=?", ft_id).Cols("identity_card", "actual_name", "id_behind_cover", "id_front_cover").Update(ft_cold)
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
	return auth.IdCardState, datastruct.NULLError
}

func (handle *DBHandler) GetAppraised(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select ap.user_id,ap.score,ap.mark,ap.appraised_type,ap.desc,ap.is_anonym,ap.created_at,pr.product_name from appraised_info ap join product_info pr on ap.product_id = pr.id join shop_info sh on sh.id = pr.shop_id where sh.f_t_id = ? LIMIT ? ?"
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
