package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"

    users "github.com/erfnzmn/Library_Management_System/internal/users"
	"github.com/erfnzmn/Library_Management_System/pkg/redisclient"
	"github.com/erfnzmn/Library_Management_System/pkg/rate"
	books "github.com/erfnzmn/Library_Management_System/internal/books"
	loans "github.com/erfnzmn/Library_Management_System/internal/loans"

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
	    Redis struct {
        Enabled  bool   `mapstructure:"enabled"`
        Addr     string `mapstructure:"addr"`
        Password string `mapstructure:"password"`
        DB       int    `mapstructure:"db"`
    } `mapstructure:"redis"`

	JWT struct {
		Secret    string `mapstructure:"secret"`
		ExpiresIn string `mapstructure:"expires_in"`
	} `mapstructure:"jwt"`
}

func loadConfig() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("jwt.secret")
	_ = viper.BindEnv("jwt.expires_in")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	// نهایی‌سازی از ENV
	if s := viper.GetString("jwt.secret"); s != "" {
		cfg.JWT.Secret = s
	}
	if s := viper.GetString("jwt.expires_in"); s != "" {
		cfg.JWT.ExpiresIn = s
	}
	return &cfg, nil
}

func openDB(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.Database.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
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

    // Echo app
    e := echo.New()
    e.HideBanner = true

    // Base middlewares
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())
    e.Use(middleware.CORS())
    e.Use(middleware.Secure())

    // Health check
    e.GET("/healthz", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]any{
            "ok":   true,
            "time": time.Now().UTC(),
        })
    })

    // Database
    var db *gorm.DB
    if cfg.Database.Enabled {
        log.Printf("DB connecting to %s:%d (db=%s)...", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
        db, err = openDB(cfg)
        if err != nil {
            log.Fatalf("db error: %v", err)
        }
		if err := db.AutoMigrate(&books.Book{}, &books.Favorite{}, &loans.Loan{}); err != nil {
    log.Fatalf("failed to migrate database: %v", err)
}
        log.Printf("DB connected ✔")
        defer func() {
            if sqlDB, _ := db.DB(); sqlDB != nil {
                _ = sqlDB.Close()
            }
        }()
    }

var rdb *redis.Client
var loginLimiter *rate.Limiter
if cfg.Redis.Enabled {
	rdb, err = redisclient.New(redisclient.Config{
		Enabled:  cfg.Redis.Enabled,
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatalf("redis error: %v", err)
	}
	if rdb != nil {
		defer rdb.Close()
		log.Printf("Redis connected ✔")
	}
}

// login limiter setup
if rdb != nil {
	loginLimiter = rate.NewTokenBucket(rdb, 5, 1, 2*time.Minute, 20*time.Minute)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("loginLimiter", loginLimiter)
			return next(c)
		}
	})
}



    // JWT setup
    jwtSecret := cfg.JWT.Secret
    jwtTTL, err := time.ParseDuration(cfg.JWT.ExpiresIn)
    if err != nil || jwtTTL <= 0 {
        jwtTTL = time.Hour
    }

    // Register routes
	if db != nil {
        users.RegisterUserRoutes(e, db, jwtSecret, jwtTTL)
	}
    if db != nil {
		// Books
		booksRepo := books.NewRepository(db)
		booksService := books.NewService(booksRepo, rdb)
		booksHandler := books.NewHandler(booksService)
		booksHandler.RegisterRoutes(e)

		// Loans
		loansRepo := loans.NewRepository(db)
		loansService := loans.NewService(db, loansRepo, booksRepo)
		loansHandler := loans.NewHandler(loansService)
		loansHandler.RegisterRoutes(e)
    }
	
	

    // Start server
    addr := ":" + cfg.Server.Port
    log.Printf("server listening on %s", addr)
    if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
        e.Logger.Fatal("shutting down the server: ", err)
    }
}

