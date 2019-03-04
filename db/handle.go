package db

import (
	"shensuanzi/datastruct"
	"shensuanzi/log"
)

func (handle *DBHandler) Test() datastruct.CodeType {
	engine := handle.mysqlEngine
	ad := new(datastruct.AdInfo)
	ad.ImgUrl = "sad"
	ad.IsHidden = true
	ad.IsJumpTo = false
	ad.JumpTo = "SS"
	ad.Platform = datastruct.H5
	_, err := engine.InsertOne(ad)
	if err != nil {
		log.Error("Test InsertOne Ad err:%", err.Error())
	}
	return datastruct.NULLError
}
