package service

import (
	"database/sql"
	"github.com/Russiancold/testApp/database"
	"github.com/jmoiron/sqlx"
	"log"
)

func Init() {
	db := database.New()
	ds, err := db.Preparex(`update test.accounts set is_del = TRUE where username = $1 and not is_del`)
	if err != nil {
		log.Fatal(err)
	}
	stbn, err := db.Preparex(`select api_token from test.accounts where username = $1 and not is_del`)
	if err != nil {
		log.Fatal(err)
	}
	sbn, err := db.Preparex(`select id from test.accounts where username = $1 and not is_del`)
	if err != nil {
		log.Fatal(err)
	}
	check, err := db.Preparex(`select id from test.accounts where (username = $1 or email = $2) and not is_del`)
	if err != nil {
		log.Fatal(err)
	}
	scbn, err := db.Preparex(`select count(id) from test.accounts where username = $1 and not is_del`)
	if err != nil {
		log.Fatal(err)
	}
	cr, err := db.Preparex(`insert into test.accounts(username, email) values ($1, $2)`)
	if err != nil {
		log.Fatal(err)
	}
	updt, err := db.Preparex(`update test.accounts set api_token = $1 where username = $2 and not is_del`)
	if err != nil {
		log.Fatal(err)
	}
	service = &accountService{db: db}
	service.delete = ds
	service.getToken = stbn
	service.create = cr
	service.selectByName = sbn
	service.updateByName = updt
	service.countByName = scbn
	service.checkExist = check
}

type accountService struct {
	db           *sqlx.DB
	create       *sqlx.Stmt
	updateByName *sqlx.Stmt
	delete       *sqlx.Stmt
	selectByName *sqlx.Stmt
	getToken     *sqlx.Stmt
	countByName  *sqlx.Stmt
	checkExist   *sqlx.Stmt
}

var service *accountService

func GetService() *accountService {
	return service
}

func (s *accountService) DeleteAccount(email string) error {
	_, err := s.delete.Exec(email);
	return err
}

func (s *accountService) CreateAccount(name, email string) error {
	if res, err := s.checkExist.Query(name, email); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		_, err := s.create.Exec(name, email)
		return err
	} else {
		if !res.Next() {
			_, err := s.create.Exec(name, email)
			return err
		}
	}
	return AlreadyExist
}

func (s *accountService) GetToken(name string) (string, error) {
	var id int
	if err := s.selectByName.Get(&id, name); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		return "", NoUser
	}
	var token sql.NullString
	if err := s.getToken.Get(&token, name); err != nil {
		if err == sql.ErrNoRows {
			return "", NoContent
		}
		return "", err
	}
	if !token.Valid {
		return "", NoContent
	}
	return token.String, nil
}

func (s *accountService) UpdateToken(name, token string) error {
	var count int
	if err := s.countByName.Get(&count, name); err != nil {
		return err
	}
	if count == 0 {
		return NoUser
	}
	_, err := s.updateByName.Exec(token, name)
	return err
}

func (s *accountService) Close() {
	if err := s.create.Close(); err != nil {
		log.Println(err)
	}
	if err := s.updateByName.Close(); err != nil {
		log.Println(err)
	}
	if err := s.delete.Close(); err != nil {
		log.Println(err)
	}
	if err := s.selectByName.Close(); err != nil {
		log.Println(err)
	}
	if err := s.getToken.Close(); err != nil {
		log.Println(err)
	}
	if err := s.countByName.Close(); err != nil {
		log.Println(err)
	}
	if err := s.db.Close(); err != nil {
		log.Println(err)
	}
}
