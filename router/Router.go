package router

import (
	"net/http"
	"p4_web/constant"
	"p4_web/controller/UserController"

	// db "p4_web/database"
	"p4_web/middleware"
	"p4_web/tools/auth"
	"p4_web/tools/exception"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/go-oauth2/mysql/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	// _ "github.com/go-sql-driver/mysql"
)

var srv *server.Server

func InitRouter() *gin.Engine {
	r := gin.Default()
	CorsInit(r)
	r.Use(middleware.MyRecovery())

	r.GET("/panic", func(c *gin.Context) {
		// var data map[string]interface{}
		// c.BindJSON(&data)
		panic(exception.ApiException{
			Code:    []int{1000},
			Message: "hello~",
		})
	})

	// oauth
	OauthInit(r)

	oauthGroup := r.Group("/oauth2")
	{
		r.GET("/oauth2/token", func(c *gin.Context) {
			err := srv.HandleTokenRequest(c.Writer, c.Request)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		})
		oauthGroup.Use(middleware.Oauth(srv))
		oauthGroup.GET("/test", func(c *gin.Context) {
			token, _ := c.Get(constant.OAUTH)
			c.Set(constant.RESPONSE, token)
		})
		oauthGroup.GET("users", UserController.GetList)
	}
	api := r.Group("/api")
	{
		api.Use(middleware.Api())
		api.POST("register", UserController.Register)
		api.POST("login", UserController.Login)
		authGroup := api.Group("/")
		{
			authGroup.Use(auth.Middleware())
			authGroup.GET("self", UserController.GetSelf)
		}

	}
	return r
}

func CorsInit(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           90 * time.Hour,
	}))
}

func OauthInit(r *gin.Engine) {
	manager := manage.NewDefaultManager()
	// mysql_store := mysql.NewDefaultStore(
	// 	mysql.NewConfig(db.Dsn()),
	// )

	// defer mysql_store.Close()
	// manager.MapTokenStorage(mysql_store)			e
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		// Domain: "http://localhost",

	})
	clientStore.Set("000001", &models.Client{
		ID:     "000001",
		Secret: "999999",
		// Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)
	srv = server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
}
