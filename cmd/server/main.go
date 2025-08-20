package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Dialect  string `mapstructure:"dialect"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		DSN      string `mapstructure:"dsn"`
		Enabled  bool   `mapstructure:"enabled"`
	} `mapstructure:"database"`
	JWT struct {
		Secret    string `mapstructure:"secret"`
		ExpiresIn string `mapstructure:"expires_in"`
	} `mapstructure:"jwt"`
}

func loadConfig() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// اتصال به MySQL با GORM
func openDB(cfg *Config) (*gorm.DB, error) {
	// اگر dsn مستقیم در config ندادی، از فیلدها بساز
	dsn := cfg.Database.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// تنظیمات کانکشن‌پول و تست اتصال
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil) // برای dev هشدار پروکسی را می‌گیرد

	// health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true, "time": time.Now().UTC()})
	})

	// اگر در config فعال باشد، به DB وصل شو
	if cfg.Database.Enabled {
		log.Printf("DB connecting to %s:%d (db=%s)...", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
		db, err := openDB(cfg)
		if err != nil {
			log.Fatalf("db error: %v", err)
		}
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		log.Printf("DB connected ✔")
	}

	log.Printf("server listening on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
