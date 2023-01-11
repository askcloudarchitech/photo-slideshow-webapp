package main

import (
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gmorse81/party-slideshow/backendServer/photodb"
	"github.com/gmorse81/party-slideshow/backendServer/photoprocessing"
	"golang.org/x/sync/errgroup"
)

var (
	sessionSecret = getEnvDefault("SESSION_SECRET", "defaultSecret")
	appPassword   = getEnvDefault("APP_PASSWORD", "defaultPassword")
)

type Routes struct {
	Database photodb.DB
}

func main() {

	paths := []string{"/data/photos", "/data/database"}
	for _, v := range paths {
		if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(v, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}

	database := photodb.DB{}
	database.InitDB()

	routes := Routes{
		Database: database,
	}

	r := gin.Default()

	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("appSession", store))

	r.MaxMultipartMemory = 8 << 20
	r.GET("/ping", routes.ping)
	r.GET("/api/slideshow/next", routes.getNextPhoto)
	r.GET("/api/is-authenticated", routes.isAuthenticated)
	r.POST("/api/login", routes.login)
	r.POST("/api/upload", routes.uploadPhoto)
	r.Static("/assets", "/dist/static/assets")
	r.Static("/photos", "/data/photos")
	r.StaticFile("/", "/dist/static/index.html")
	r.StaticFile("/login", "/dist/static/index.html")
	r.StaticFile("/tv-slideshow", "/dist/static/index.html")
	r.StaticFile("/favicon.ico", "/dist/static/favicon.ico")

	errs, _ := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		err := r.Run()
		return err
	})

	err := errs.Wait()
	log.Fatal(err)

}

func getEnvDefault(env string, defaultVal string) string {
	val := os.Getenv(env)
	if val == "" {
		return defaultVal
	}
	return val
}

func (r *Routes) login(c *gin.Context) {
	if c.Request.FormValue("password") == appPassword {
		session := sessions.Default(c)
		session.Set("authenticated", "true")
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"authenticated": "true",
		})
	} else {
		c.JSON(401, gin.H{
			"authenticated": "false",
		})
	}
}

func (r *Routes) isAuthenticated(c *gin.Context) {
	err := authCheck(c)
	if err != nil {
		c.JSON(401, gin.H{
			"authenticated": "false",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": "true",
		})
	}
}

func authCheck(c *gin.Context) error {
	session := sessions.Default(c)
	if session.Get("authenticated") == "true" {
		return nil
	}
	return fmt.Errorf("not authenticated")
}

func (r *Routes) uploadPhoto(c *gin.Context) {
	if c.Query("p") != appPassword {
		err := authCheck(c)
		if err != nil {
			c.JSON(401, gin.H{
				"authenticated": "false",
				"error":         "not authenticated",
			})
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed file upload",
		})
		return
	}

	go func() {
		cleanFile, dateTaken, err := photoprocessing.DetectAndConvertImage(file)
		if err != nil {
			log.Print(err)
			// c.JSON(http.StatusBadRequest, gin.H{
			// 	"error": fmt.Sprintf("failed to convert image: %s", err),
			// })
			return
		}
		timestamp := time.Now().UnixMilli()

		out, _ := os.Create(fmt.Sprintf("/data/photos/%d.jpeg", timestamp))
		defer out.Close()

		var opts jpeg.Options
		opts.Quality = 90

		err = jpeg.Encode(out, cleanFile, &opts)
		if err != nil {
			log.Print(err)
			// c.Error(err)
			// c.JSON(http.StatusBadRequest, gin.H{
			// 	"error": "failed to write image",
			// })
			return
		}
		r.Database.AddPhoto(fmt.Sprintf("%d.jpeg", timestamp), int(dateTaken))
	}()

	c.Status(201)
}

func (r *Routes) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (r *Routes) getNextPhoto(c *gin.Context) {
	p, _ := r.Database.GetNextSlideshowPhoto()
	err := r.Database.UpdateLastViewed(p.Name)
	p.Name = fmt.Sprintf("/photos/%s", p.Name)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, p)
}
