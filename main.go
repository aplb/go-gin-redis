package main

import (
    "net/http"
    "fmt"
    "log"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis"
)

var redisClient *redis.Client

func setupRedis() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })

    pong, err := redisClient.Ping().Result()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(pong)
}

func setupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })

    authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
        "foo": "123",
    }))

    authorized.POST("admin", func(c *gin.Context) {
        user := c.MustGet(gin.AuthUserKey).(string)

        var json struct {
            Value string `json:"value" binding:"required"`
        }

        err := c.Bind(&json)
        if err == nil {
            err = redisClient.Set(user, json.Value, 0).Err()
            if err != nil {
                panic(err)
            }

            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })

    return r
}

func main() {
    setupRedis()
    r := setupRouter()
    r.Run(":8080")
}

