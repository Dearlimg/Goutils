package utils

import (
	"crypto/md5"
	"encoding/hex"
)

/*
该方法用于针对上传后的文件名格式化，简单来说，将文件名 MD5 后再进行写入，防止直接把原始名称暴露出去
*/

func EncodeMD5(value string) string {
	m := md5.New()                        // 创建一个新的 MD5 哈希对象
	m.Write([]byte(value))                // 向哈希对象中写入数值的字节表示
	return hex.EncodeToString(m.Sum(nil)) // 返回哈希值的十六进制编码表示
}
