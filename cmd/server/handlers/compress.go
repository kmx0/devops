package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type gzipGINWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w gzipGINWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func Compress() gin.HandlerFunc {
	return func(c *gin.Context) {

		logrus.Info(c.GetHeader("Accept-Encoding"))
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			// io.WriteString(c.Writer, err.Error())

			c.Status(http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		// access the status we are sending
		status := c.Writer.Status()
		c.Header("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных

		c.Writer = gzipGINWriter{ResponseWriter: c.Writer, Writer: gz}
		// c.Status(http.StatusOK)
		c.Next()
		logrus.Info(status)
	}
}

func Decompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reader io.Reader
		logrus.Info(c.GetHeader("Content-Encoding"))
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = c.Request.Body
			return
		}
		body, err := io.ReadAll(reader)
		// access the status we are sending
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header("Content-Encoding", "gzip")
		c.Header("Content-Length", fmt.Sprintf("%d", len(body)))

		c.String(http.StatusOK, fmt.Sprintf("%d", len(body)))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		c.Next()
	}
}

