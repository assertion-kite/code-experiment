package main

import (
	"code/route"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/sharding"
	"log"
	"os"
	"time"
)

type A struct {
	Name string
}

type B struct {
	Name string
}

type C interface {
	A | B
}

func main() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
		"go",
		"golang",
		"120.55.98.12",
		"3306",
		"go",
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		}),
	})
	if err != nil {
		return
	}
	shardingConfig := sharding.Config{
		ShardingKey:         "date_time",
		NumberOfShards:      1024,
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return 0
		},
		ShardingAlgorithm: func(columnValue any) (suffix string, err error) {
			if dateTime, ok := columnValue.(time.Time); ok {
				return fmt.Sprintf("_%v", dateTime.Year()), nil
			}
			return "", errors.New("invalid data_time")
		},
	}
	err = db.Use(sharding.Register(shardingConfig, "user"))
	if err != nil {
		return
	}
	route.RegisterRoute([]route.RegisterRouteFunc{
		NewUserWeb(db).RegisterUserRoute,
		NewOrderWeb(db).RegisterOrderRoute,
	})
}
