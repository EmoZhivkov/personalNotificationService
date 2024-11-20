package repositories

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tablesToCreate = []any{&User{}, &UserNotificationChannels{}, &Template{}, &Notification{}}

type DbClient interface {
	Create(value interface{}) (tx *gorm.DB)
	Take(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
}

type dbClient struct {
	*gorm.DB
}

func NewDbClient(dbConnectionUrl string) (DbClient, error) {
	client, err := getPgClient(dbConnectionUrl)
	if err != nil {
		return nil, err
	}

	if err := client.AutoMigrate(tablesToCreate...); err != nil {
		return nil, err
	}

	return &dbClient{DB: client}, err
}

func getPgClient(dbConnectionUrl string) (*gorm.DB, error) {
	dsn, err := pgURLToDSN(dbConnectionUrl)
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func pgURLToDSN(pgURL string) (string, error) {
	u, err := url.Parse(pgURL)
	if err != nil {
		return "", err
	}

	user := u.User.Username()
	password, _ := u.User.Password()

	host := u.Hostname()
	port := u.Port()

	dbName := u.Path[1:]

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port)

	return dsn, nil
}
