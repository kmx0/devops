package handlers

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/devops/internal/crypto"
	"github.com/sirupsen/logrus"
)

type aesWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// type gzipGINWriter struct {
// 	gin.ResponseWriter
// 	Writer io.Writer
// }

func (w aesWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// Сжимает отправляемые данные, если клиент поддерживает сжатие
// поле Accept-Encoding.
func Decrypt(privateKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!111")
		logrus.Info("decryptinhg")
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		rawdata, err := c.GetRawData()
		if err != nil {
			logrus.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		decryptedbytes, err := crypto.DecryptData(privateKey, rawdata)
		if err != nil {
			logrus.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header("Content-Length", fmt.Sprintf("%d", len(decryptedbytes)))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedbytes))
		//  Go through the process
		c.Next()
	}
}

// func handleAes(c *gin.Context, md5key string) {
// 	contentType := c.Request.Header.Get("Content-Type")
// 	isJsonRequest := strings.Contains(contentType, "application/json")
// 	isFileRequest := strings.Contains(contentType, "multipart/form-data")
// 	isFormUrl := strings.Contains(contentType, "application/x-www-form-urlencoded")

// 	if c.Request.Method == "GET" {
// 		err := parseQuery(c, md5key)
// 		if err != nil {
// 			log("handleAes parseQuery  err:%v", err)
// 			// The output here should be ciphertext   Once the encryption and decryption are debugged   You won't come in here
// 			response(c, 2001, " System error ", err.Error())
// 			return
// 		}
// 	} else if isJsonRequest {
// 		err := parseJson(c, md5key)
// 		if err != nil {
// 			log("handleAes parseJson err:%v", err)
// 			// The output here should be ciphertext   Once the encryption and decryption are debugged   You won't come in here
// 			response(c, 2001, " System error ", err.Error())
// 			return
// 		}
// 	} else if isFormUrl {
// 		err := parseForm(c, md5key)
// 		if err != nil {
// 			log("handleAes parseForm err:%v", err)
// 			// The output here should be ciphertext   Once the encryption and decryption are debugged   You won't come in here
// 			response(c, 2001, " System error ", err.Error())
// 			return
// 		}
// 	} else if isFileRequest {
// 		err := parseFile(c, md5key)
// 		if err != nil {
// 			log("handleAes parseFile err:%v", err)
// 			// The output here should be ciphertext   Once the encryption and decryption are debugged   You won't come in here
// 			response(c, 2001, " System error ", err.Error())
// 			return
// 		}
// 	}

// 	/// Intercept  response body
// 	oldWriter := c.Writer
// 	blw := &aesWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
// 	c.Writer = blw

// 	//  Go through the process
// 	c.Next()

// 	/// Get the returned data
// 	responseByte := blw.body.Bytes()

// 	// journal
// 	c.Writer = oldWriter
// 	// If not json Format   Then go straight back , It should be for file download and so on. It should not be encrypted
// 	if !isJsonResponse(c) {
// 		_, _ = c.Writer.Write(responseByte)
// 		return
// 	}

// 	/// encryption
// 	encryptStr, err := aes.GcmEncrypt(md5key, string(responseByte))
// 	if err != nil {
// 		log("handleAes GcmEncrypt err:%v", err)
// 		response(c, 2001, " System error ", err.Error())
// 		return
// 	}

// 	_, _ = c.Writer.WriteString(encryptStr)
// }

// // Распаковывет принятые данные, в случае если клиент поддерживает сжатие
// // поле Accept-Encoding.
// func Decrypt() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var reader io.Reader
// 		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
// 			gz, err := gzip.NewReader(c.Request.Body)
// 			if err != nil {
// 				c.Status(http.StatusInternalServerError)
// 				return
// 			}
// 			reader = gz
// 			defer gz.Close()
// 		} else {
// 			return
// 		}
// 		body, err := io.ReadAll(reader)
// 		// access the status we are sending
// 		if err != nil {
// 			c.Status(http.StatusInternalServerError)
// 			return
// 		}
// 		c.Header("Content-Encoding", "gzip")
// 		c.Header("Content-Length", fmt.Sprintf("%d", len(body)))

// 		c.String(http.StatusOK, fmt.Sprintf("%d", len(body)))
// 		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
// 		c.Next()
// 	}
// }

// func AesGcmDecrypt() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer func() {
// 			if e := recover(); e != nil {
// 				stack := debug.Stack()
// 				log("AesGcmDecrypt Recovery: err:%v, stack:%v", e, string(stack))
// 			}
// 		}()

// 		if c.Request.Method == "OPTIONS" {
// 			c.Next()
// 		} else {
// 			md5key := aes.GetAesKey("gavin12345678")
// 			log("AesGcmDecrypt start url:%s  ,md5key:%s, Method:%s, Header:%+v", c.Request.URL.String(), md5key, c.Request.Method, c.Request.Header)
// 			handleAes(c, md5key)
// 		}
// 	}
// }
