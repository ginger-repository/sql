package sql

import (
	"fmt"
	"time"

	"github.com/ginger-core/compound/registry"
	"github.com/ginger-core/errors"
	"github.com/ginger-core/log"
	_query "github.com/ginger-core/query"
	"github.com/ginger-repository/sql/query"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repo struct {
	logger log.Logger
	*gorm.DB

	config config
}

func New(logger log.Logger, registry registry.Registry) Repository {
	db := &repo{
		logger: logger,
	}

	if err := registry.Unmarshal(&db.config); err != nil {
		panic(err)
	}
	db.config.initialize()

	return db
}

func (repo *repo) Initialize() errors.Error {
	var dialect gorm.Dialector
	switch repo.config.Dialect {
	case "mysql":
		dialect = mysql.Open(repo.config.ConnectionString)
	case "postgres":
		dialect = postgres.Open(repo.config.ConnectionString)
	default:
		return errors.New().WithMessage(fmt.Sprintf("invalid dialect %s", repo.config.Dialect))
	}

	d, err := gorm.Open(dialect, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 repo.newLogger(repo.logger),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return errors.New(err)
	}
	sqlDB, err := d.DB()
	if err != nil {
		return errors.New(err)
	}

	if repo.config.ConnMaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(*repo.config.ConnMaxLifetime)
	}
	if repo.config.MaxIdleConnections != nil {
		sqlDB.SetMaxIdleConns(*repo.config.MaxIdleConnections)
	}
	if repo.config.MaxOpenConnections != nil {
		sqlDB.SetMaxOpenConns(*repo.config.MaxOpenConnections)
	}

	if repo.config.Debug {
		repo.logger.Debugf("switching into debug mode")
		d = d.Debug()
	}

	repo.DB = d

	return nil
}

func (repo *repo) GetDB(q _query.Query) any {
	if q != nil {
		db := q.GetDB()
		if db != nil {
			sqlDb := db.(*gorm.DB)
			if sqlDb != nil {
				return db
			}
		}
	}
	db := repo.DB
	if q != nil {
		model, _ := query.FindModel(q)
		if model != nil {
			db = db.Model(model)
		}
	}
	return db
}

func (repo *repo) getDB(query _query.Query) *gorm.DB {
	return repo.GetDB(query).(*gorm.DB)
}
