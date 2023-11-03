package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	// "atm5_microservices/logger"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gopkg.in/Shopify/sarama.v1"

	_transactionHandler "github.com/gamepkw/transactions-banking-microservice/internal/handlers"
	_transactionRepostitory "github.com/gamepkw/transactions-banking-microservice/internal/repositories"
	_transactionService "github.com/gamepkw/transactions-banking-microservice/internal/services"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// logger.Info("start program...")

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	dbconnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "true")
	val.Add("loc", "Asia/Bangkok")
	dsn := fmt.Sprintf("%s?%s", dbconnection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	redisHost := viper.GetString(`redis.host`)
	redisdbPort := viper.GetString(`redis.port`)
	redisdbPass := viper.GetString(`redis.pass`)

	addr := fmt.Sprintf("%s:%s", redisHost, redisdbPort)
	password := redisdbPass

	redis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	fmt.Println(redis)

	config := sarama.NewConfig()
	config.ClientID = "my-kafka-client"
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	brokers := []string{viper.GetString("kafka.broker_address")}
	kafkaClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaClient.Close()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	transactionRepo := _transactionRepostitory.NewTransactionRepository(dbConn, redis)
	transactionService := _transactionService.NewTransactionService(transactionRepo, timeoutContext, redis)
	_transactionHandler.NewTransactionHandler(e, transactionService, redis)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
