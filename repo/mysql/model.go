package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/novan/golang-api-server/repo/mysql/schema"
	"github.com/novan/golang-api-server/util"
)

// var DB *sqlx.DB
func Connect() *sqlx.DB {
	driverName := os.Getenv("MYSQL_DRIVER")
	usernameAndPassword := fmt.Sprintf("%s:%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"))
    hostName := fmt.Sprintf("tcp(%s:%s)", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
    urlConnection := fmt.Sprintf("%s@%s/%s?charset=%s&collation=%s&parseTime=true&loc=%s", 
    	usernameAndPassword, 
    	hostName, 
    	os.Getenv("MYSQL_DATABASE"), 
    	os.Getenv("MYSQL_CHARSET"), 
    	os.Getenv("MYSQL_COLLATION"), 
    	url.QueryEscape(os.Getenv("APP_TIMEZONE")),
    )

	db, err :=  sqlx.Connect(driverName, urlConnection)
	if err != nil {
		log.Panicf("⇨ %s Data source %s , Failed : %s \n", os.Getenv("MYSQL_DRIVER"), urlConnection, err.Error())
	}
	db.SetMaxOpenConns(util.AtoI(os.Getenv("MYSQL_MAX_CONN")))
	db.SetMaxIdleConns(util.AtoI(os.Getenv("MYSQL_MAX_IDLE_CONN")))
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("MySQL | Database %s connection to %s is successfully connected!", driverName, hostName)
	return db
}

func ConnectCommerce() *sqlx.DB {
	driverName := os.Getenv("MYSQL_COMMERCE_DRIVER")
	usernameAndPassword := fmt.Sprintf("%s:%s", os.Getenv("MYSQL_COMMERCE_USER"), os.Getenv("MYSQL_COMMERCE_PASSWORD"))
    hostName := fmt.Sprintf("tcp(%s:%s)", os.Getenv("MYSQL_COMMERCE_HOST"), os.Getenv("MYSQL_COMMERCE_PORT"))
    urlConnection := fmt.Sprintf("%s@%s/%s?charset=%s&collation=%s&parseTime=true&loc=%s", 
    	usernameAndPassword, 
    	hostName, 
    	os.Getenv("MYSQL_COMMERCE_DATABASE"), 
    	os.Getenv("MYSQL_COMMERCE_CHARSET"), 
    	os.Getenv("MYSQL_COMMERCE_COLLATION"), 
    	url.QueryEscape(os.Getenv("APP_TIMEZONE")),
    )

	db, err :=  sqlx.Connect(driverName, urlConnection)
	if err != nil {
		log.Panicf("⇨ %s Data source %s , Failed : %s \n", os.Getenv("MYSQL_COMMERCE_DRIVER"), urlConnection, err.Error())
	}
	db.SetMaxOpenConns(util.AtoI(os.Getenv("MYSQL_COMMERCE_MAX_CONN")))
	db.SetMaxIdleConns(util.AtoI(os.Getenv("MYSQL_COMMERCE_MAX_IDLE_CONN")))
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("MySQL Commerce | Database %s connection to %s is successfully connected!", driverName, hostName)
	return db
}

type MysqlInterface interface {
	SelectContext(ctx context.Context, db *sqlx.DB, data interface{}, query string, args ...interface{}) error
	List(context.Context, *sqlx.DB, interface{}, *util.Query, string) error
	Get(context.Context, *sqlx.DB, interface{}, *util.Query, string) error
	GetTransaction(context.Context, *sqlx.Tx, interface{}, *util.Query, string) error
	CreateUpdate(context.Context, *sqlx.DB, interface{}, string) (sql.Result, error)
	PreparedCreateUpdate(ctx context.Context, db *sqlx.DB, data interface{}, queryString string) (tableID interface{}, err error)
	CreateUpdateTransaction(context.Context, *sqlx.Tx, interface{}, string) (sql.Result, error)
	ListWithPagination(context.Context, *sqlx.DB, interface{}, *util.Query, string, string) (*schema.Pagination, error)
}

type Model struct{}

func NewModel() *Model {
	return &Model{}
}

func (c *Model) SelectContext(ctx context.Context, db *sqlx.DB, data interface{}, query string, args ...interface{}) error {
	err := db.SelectContext(ctx, data, db.Rebind(query), args...)
	return err
}

func (c *Model) List(ctx context.Context, db *sqlx.DB, data interface{}, query *util.Query, queryString string) (err error) {
	where, args := query.Where()
	q := queryString
	q += where
	q += query.Order()

	if err = db.SelectContext(ctx, data, db.Rebind(q), args...); err != nil {
		return
	}
	return
}

func (c *Model) Get(ctx context.Context, db *sqlx.DB, data interface{}, query *util.Query, queryString string) (err error) {
	where, args := query.Where()
	q := queryString
	q += where

	if err = db.GetContext(ctx, data, db.Rebind(q), args...); err != nil {
		return
	}
	return
}

func (c *Model) CreateUpdate(ctx context.Context, db *sqlx.DB, data interface{}, queryString string) (result sql.Result, err error) {
	result, err = db.NamedExecContext(ctx, queryString, data)
	if err != nil {
		return
	}
	return
}

func (c *Model) PreparedCreateUpdate(ctx context.Context, db *sqlx.DB, data interface{}, queryString string) (tableID interface{}, err error) {
	stmt, err := db.PrepareNamedContext(ctx, queryString)
	if err != nil {
		return
	}

	_ = stmt.GetContext(ctx, &tableID, data)
	return
}

func (c *Model) ListWithPagination(ctx context.Context, db *sqlx.DB, data interface{}, query *util.Query, queryString string, queryCount string) (paginate *schema.Pagination, err error) {
	where, args := query.Where()
	sort := query.Order()
	q := queryString
	q += where
	q += sort
	q += " " + query.Limit()

	// util.Log.WithContext(ctx).Debugf("ListWithPagination | Query: %s | Where: %s | Args: %s", q, where, util.ToString(args))

	if err = db.SelectContext(ctx, data, db.Rebind(q), args...); err != nil {
		return
	}
	
	var count int
	err = c.Get(ctx, db, &count, query, queryCount)
	if err != nil {
		return
	}

	paginate = &schema.Pagination{
		CurrentPage: query.Page,
		PageSize:    query.Count,
		TotalPage:   count,
		TotalResult: count,
	}
	paginate.SetTotalPage(count)
	return
}

func (c *Model) CreateUpdateTransaction(ctx context.Context, tx *sqlx.Tx, data interface{}, queryString string) (result sql.Result, err error) {
	result, err = tx.NamedExecContext(ctx, queryString, data)
	if err != nil {
		return
	}
	return
}

func (c *Model) GetTransaction(ctx context.Context, tx *sqlx.Tx, data interface{}, query *util.Query, queryString string) (err error) {
	where, args := query.Where()
	q := queryString
	q += where

	if err = tx.GetContext(ctx, data, tx.Rebind(q), args...); err != nil {
		return
	}
	return
}
