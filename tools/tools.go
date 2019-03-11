package tools

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"shensuanzi/datastruct"
	"shensuanzi/log"
	"shensuanzi/thirdparty/tls-sig-api-golang"
	"strconv"
	"strings"
)

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func GetFileContentAsStringLines(filePath string) ([]string, error) {
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error("read file: %v error: %v", filePath, err)
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	return result, nil
}

func AccountGenForIM(user_identifier string, appid int) (string, datastruct.CodeType) {
	privateKey, _ := GetFileContentAsStringLines("datastruct/important/keys/private_key")
	publicKey, _ := GetFileContentAsStringLines("datastruct/important/keys/public_key")
	var privateKey_str string
	var publicKey_str string
	var buf1 bytes.Buffer

	for _, v := range privateKey {
		buf1.WriteString(v)
		buf1.WriteString("\n")
	}
	privateKey_str = buf1.String()

	var buf2 bytes.Buffer
	for _, v := range publicKey {
		buf2.WriteString(v)
		buf2.WriteString("\n")
	}
	publicKey_str = buf2.String()

	expire := 24 * 3600 * 8 //8天后失效
	userSig, err := TLSSigAPI.GenerateUsersigWithExpire(privateKey_str, appid, user_identifier, int64(expire))
	if err != nil {
		log.Error("GenerateUsersigWithExpire:%v", err.Error())
		return "", datastruct.GetDataFailed
	}
	err = TLSSigAPI.VerifyUsersig(publicKey_str, userSig, appid, user_identifier)
	if err != nil {
		log.Error("VerifyUsersig:%v", err.Error())
		return "", datastruct.GetDataFailed
	}

	return userSig, datastruct.NULLError
}

func StringToBool(value string) bool {
	tf := false
	if value == "1" {
		tf = true
	}
	return tf
}

func StringToInt(value string) int {
	rs, _ := strconv.Atoi(value)
	return rs
}

func StringToInt64(tmp string) int64 {
	rs, _ := strconv.ParseInt(tmp, 10, 64)
	return rs
}

func StringToFloat64(tmp string) float64 {
	rs, _ := strconv.ParseFloat(tmp, 64)
	return rs
}

func StringToIdCardState(value string) datastruct.IdCardState {
	rs := StringToInt(value)
	return datastruct.IdCardState(rs)
}

// func StringToOnlineState(value string) datastruct.FTOnlineState {
// 	rs := StringToInt(value)
// 	return datastruct.FTOnlineState(rs)
// }
