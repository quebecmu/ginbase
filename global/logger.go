package global

import (
	"bytes"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/dgrijalva/jwt-go"
	"github.com/mssola/useragent"
	"golang.org/x/sync/errgroup"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LOG *zap.Logger

// InitLogger 初始化lg
func InitLogger(MaxBackups, MaxAge, MaxSize int, path, Level string) (err error) {
	writeSyncer := getLogWriter(path, MaxSize, MaxBackups, MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(Level))
	if err != nil {
		return
	}
	var core zapcore.Core

	core = zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, l),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
	)

	LOG = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(LOG)
	zap.L().Info("init logger success")
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		wg := errgroup.Group{}
		wg.Go(func() error {
			start := time.Now()
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			ua := DeUa(c.Request.UserAgent())
			execTime := time.Since(start)
			c.Next()
			userid, err := enToken(c.Request.Header.Get("Authorization"))
			if err != nil {
				LOG.Info(path,
					zap.String("user", ""),
					zap.String("BrowserName", ua.Name),
					zap.String("BrowserVersion", ua.Version),
					zap.String("Platform", ua.Platform),
					zap.String("OS", ua.OS),
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("Body", ReadCloserToString(c.Request.Body)),
					zap.String("ip", c.ClientIP()),
					zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
					zap.Duration("execTime", execTime))
				c.Next()
				return nil
			}

			LOG.Info(path,
				zap.String("user", userid.UserId),
				zap.String("BrowserName", ua.Name),
				zap.String("BrowserVersion", ua.Version),
				zap.String("Platform", ua.Platform),
				zap.String("OS", ua.OS),
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("Body", ReadCloserToString(c.Request.Body)),
				zap.String("ip", c.ClientIP()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("execTime", execTime),
			)
			c.Next()
			return nil
		})
		c.Next()
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		wg := errgroup.Group{}
		wg.Go(func() error {
			defer func() {
				if err := recover(); err != nil {
					// Check for a broken connection, as it is not really a
					// condition that warrants a panic stack trace.
					var brokenPipe bool
					if ne, ok := err.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
								brokenPipe = true
							}
						}
					}

					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					if brokenPipe {
						LOG.Error(c.Request.URL.Path,
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
						// If the connection is dead, we can't write a status to it.
						c.Error(err.(error)) // nolint: errcheck
						c.Abort()
						return
					}

					if stack {
						LOG.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())),
						)
					} else {
						LOG.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
					}
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}()
			c.Next()
			return nil
		})

	}
}

// ReadCloserToString 将网络请求后得到的流转换为字符串
func ReadCloserToString(closer io.ReadCloser) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(closer)
	if err != nil {
		return ""
	}
	return buf.String()
}

const PriKey = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIFLTBXBgkqhkiG9w0BBQ0wSjApBgkqhkiG9w0BBQwwHAQI124qNyuxH/4CAggA
MAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAEqBBBnVdo+pXE/iWm4gIC9+CQzBIIE
0M1RViZFj/hXU0Czy+Sohu0JBHDMRUMt09lWEnOToT7wHJRqBMWbGtshn96EJehO
F32gzBITOP88oW71R38PsC306JGghZrRnGnMAASN0cb1ZxqwBpmXRdw7KPalLQcj
AF4Z0rhvOSWDcJpqZuC7svMs6Dbgl6fCKrSF6Z7IfRpYLrzdRn4xpuklDUTZqLij
PLC8TvW93Ry3/pJpgTwIqgjlSwDBQlK11bDirj9eAK3BSuc8z6XN0wd+oMGj9t3M
U3mIsutZhCfq/VXwPzNrBwWXy3/LT99exhEQCu6VAOyY4VtGddwaeWQelghrB4Nk
PiMcxzadQ7hYL1Ervj7UcinZV+mg2XAnbolkRni1GJwDhqQqnSHMDOzRVIbRWDGx
Mhgy34XH2brtIQgCEyStFXwQA2wDZ1DuK6BtJ/jnPogwUPGKLWKsvkItJfSVjnhh
v1Qci7d7lLKL11WTjUoSy2wJhL89IQkkjF8oAytJuGjQplVhYBdhNyMSCkvD5fXP
k0OK9AQ4bQHaianPzycTW3P6a5ZPB6kLk7xFRd5YdE5ToK6glQVYyUwThHGj9pRb
0c+aKI2ceT5t0FVdjxmKalrMB0F37Nhnkcr6h/VAmmZ8XIyjV/qMYoSYNxREtZox
Kj1VW9igyQCf72M0zaXnSASdBV+u4+Md69/mygSDO71qL4xYQL8VT+k8zJe3T5nh
GpwOQi1Dl1S4t0XRxNhgX5W8NeefsZThp6RomqULQen0+CfNqlz0Zge7/kMGAdFx
8SdMS93F3dgsY4wx+ufh/s4b1aPz1SI7gfgarND0Gmgvk5OWon1YbU4ozmb8Y4Cj
RNMha+JfPDasE6P4TjseESVGy1Cqw5dKEdchD3oX4n4vC8oasDPex54JDd5mIgf2
txOX67XrOR5LuwoWTusuhyjeAUWcUGkwExU299vYs4M1+y+dQZ8VgP3skVOT8kz4
1Mh3u9PFbbato4cLFeWS7YoCHvQOJ0vmENTOl5Z80ILI4AJHmbHuYTO17fyuy0pi
AZmZyfppCxsrWdmpos+NMKVsWJiyk6JOoltsszmE8kYXywuq46YH1ERx4wgBeXGy
blyeU6pVHXBSO4JDNDbQte58PEe1NAcWoUnyERBAbb85y/27+wrs4Ym6rvRjczCx
WovfB8DYC1A8iUk3wWWpaQu9MRk/PmaY142r/Pd/20D6QW/6KUGDCmBE+pf99L3+
CWVpvClFcHA0yKb2cVrFjZtm2IirhyWhivbouHZ/Fk0LYXWdsC8Uf8KCnN7is8V+
tocHP2ztCFN2pQWUV8b+lNDCL8tPRH++7OhaaQsOdLDOBUZa0b/88CqrWQHsoN3d
hpu6FbyL6+Nmahu5Q3E2pF1wU7Dqw3ZSTRzbdza/CLAud0aC8bf4j9gk5LFdipA4
8NCYfCSPlkT5c5qJ/ecgDDrHdpNnGP3ccjP5DEAUjfbKfpASGCd+G5fiwOy00uih
FI4NutfgISWWKfhwbbF5/Hb9dLiMVmrIBseDtAW4FnTv7KjvUiiXq8A2snbNkpi1
cFx2fgvLGFvPhxQTWdmnCrw5Fyet6A+IHA+hyT0xjwvodGOexojM1Uij9Oebt0jF
JqsGD0x33ZfgN9XOP4+oZCTyruiNrRIeSJk1+talEC1W
-----END ENCRYPTED PRIVATE KEY-----
`

type LogClaims struct {
	UserId             string `json:"userId"`
	jwt.StandardClaims        // 标准Claims结构体，可设置8个标准字段
}

// 通过jwt.ParseWithClaims返回的Token结构体取出Claims结构体
func enToken(t string) (*LogClaims, error) {
	token, err := jwt.ParseWithClaims(t, &LogClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(PriKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*LogClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil
}

type UAInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
	OS       string `json:"OSy"`
}

func DeUa(u string) *UAInfo {
	ch := make(chan *useragent.UserAgent, 1)
	gopool.Go(func() {
		ch <- useragent.New(u)
	})
	ua := <-ch
	close(ch)
	name, version := ua.Browser()
	return &UAInfo{
		Name:     name,
		Version:  version,
		Platform: ua.Platform(),
		OS:       ua.OS(),
	}
}
