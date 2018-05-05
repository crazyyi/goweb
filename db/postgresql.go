package db

import (
	"log"
	"database/sql"
	"fmt"

	"github.com/crazyyi/goweb/model"
	
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // no need to assign a name
)

type Config struct {
	ConnectString string
}

type pgDb struct {
	dbConn *sqlx.DB
	sqlSelectPeople *sqlx.Stmt
	sqlInsertPerson *sqlx.NamedStmt
	sqlSelectPerson *sql.Stmt
	sqlDeletePerson *sql.Stmt
	sqlUpdatePerson *sql.Stmt
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.createTablesIfNotExist(); err != nil {
			return nil, err
		}

		if err := p.prepareSQLStatements(); err != nil {
			return nil, err
		}

		return p, nil
	}
}

func (p *pgDb) createTablesIfNotExist() error {
	createSQL := `
		CREATE TABLE IF NOT EXISTS people (
		id SERIAL NOT NULL PRIMARY KEY,
		first TEXT NOT NULL,
		last TEXT NOT NULL);
	`

	if rows, err := p.dbConn.Query(createSQL); err != nil {
		return err
	} else {
		rows.Close()
	}

	return nil
}

func (p *pgDb) prepareSQLStatements() (err error) {
	if p.sqlSelectPeople, err = p.dbConn.Preparex(
		"SELECT id, first, last FROM people",
	); err != nil {
		return err
	}

	if p.sqlInsertPerson, err = p.dbConn.PrepareNamed(
		"INSERT INTO people (first, last) VALUES (:first, :last) RETURNING id",
	); err != nil {
		return err
	}

	if p.sqlSelectPerson, err = p.dbConn.Prepare(
		"SELECT id, first, last FROM people WHERE id = $1",
	); err != nil {
		return err
	}

	if p.sqlDeletePerson, err = p.dbConn.Prepare(
		"DELETE FROM people WHERE id = $1",
	); err != nil {
		return err
	}

	if p.sqlUpdatePerson, err = p.dbConn.Prepare(
		"UPDATE people SET first = $2, last = $3 WHERE id = $1",
	); err != nil {
		return err
	}

	return nil
}

func (p *pgDb) SelectPeople() ([]*model.Person, error) {
	people := make([]*model.Person, 0)
	if err := p.sqlSelectPeople.Select(&people); err != nil {
		return nil, err
	}

	return people, nil
}

func (p *pgDb) CreateNewRecord(record *model.Person) (int, error) {
	var data int
	err := p.sqlInsertPerson.Get(&data, record)
	if err != nil {
		return -1, err
	} 

	return data, nil
}

func (p *pgDb) DeleteRowAt(id int64) (int64, error) {
	res, err := p.sqlDeletePerson.Exec(id)

	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()

	if err != nil {
		panic(err)
	}

	return count, nil
}

func (p *pgDb) UpdateRowAt(id int64, firstname string, lastname string) (int64, error) {
	res, err := p.sqlUpdatePerson.Exec(id, firstname, lastname)

	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()

	if err != nil {
		panic(err)
	}

	return count, nil
}

