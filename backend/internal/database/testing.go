package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB initializes a test database connection
func SetupTestDB() error {
	// 设置测试数据库配置
	viper.Set("database.host", "localhost")
	viper.Set("database.port", 5432)
	viper.Set("database.user", "postgres")
	viper.Set("database.password", "postgres")
	viper.Set("database.dbname", "share_ai_platform_test")
	viper.Set("database.sslmode", "disable")

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	// 连接数据库
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %v", err)
	}

	log.Println("Test database connection established")
	return nil
}

// TeardownTestDB cleans up the test database
func TeardownTestDB() {
	if db != nil {
		// 获取底层的 *sql.DB 对象
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting database instance: %v", err)
			return
		}
		sqlDB.Close()
	}
}

// TestMain is used to setup and teardown the test database
func TestMain(m *testing.M) {
	// 设置测试环境
	if err := SetupTestDB(); err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}

	// 运行测试
	code := m.Run()

	// 清理测试环境
	TeardownTestDB()

	os.Exit(code)
}
