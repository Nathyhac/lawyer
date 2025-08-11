package store

import (
	"database/sql"
	"fmt"
	"log"
)

type Lawyer struct {
	ID           int           `json:"id"`
	First_name   string        `json:"first_name"`
	Last_name    string        `json:"last_name"`
	Email        string        `json:"email"`
	Phone_number string        `json:"phone_number"`
	AddressID    sql.NullInt64 `json:"-"`
	Address      Address       `json:"addresses"`
}
type Address struct {
	ID      int            `json:"address_id"`
	City    string         `json:"city"`
	Street  string         `json:"street"`
	State   sql.NullString `json:"_"`
	Country string         `json:"country"`
}

type PostgresLawyerDB struct {
	DB *sql.DB
}

func NewPostgresDB(db *sql.DB) *PostgresLawyerDB {
	handler := &PostgresLawyerDB{
		DB: db,
	}
	return handler
}

type LawyerInterface interface {
	CreateLawyer(*Lawyer) (*Lawyer, error)
	GetLawyerById(int64) (*Lawyer, error)
	GetAllLawyers() ([]Lawyer, error)
	UpdateLawyer(*Lawyer) error
	Deletelawyer(int64) error
}

func (pg *PostgresLawyerDB) CreateLawyer(Lawyer *Lawyer) (*Lawyer, error) {

	tx, err := pg.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	AdQuery := `INSERT into addresses(country, city, street) VALUES($1,$2,$3) RETURNING address_id`
	err = tx.QueryRow(AdQuery, Lawyer.Address.Country, Lawyer.Address.City, Lawyer.Address.Street).Scan(&Lawyer.Address.ID)
	if err != nil {
		return nil, fmt.Errorf("error happend while copying into Lawyer struct: %v", err)
	}

	query := `INSERT into lawyer(first_name, last_name, email, phone_number, address_id) VALUES($1,$2,$3,$4, $5) RETURNING id`
	err = tx.QueryRow(query, Lawyer.First_name, Lawyer.Last_name, Lawyer.Email, Lawyer.Phone_number, Lawyer.Address.ID).Scan(&Lawyer.ID)
	if err != nil {
		return nil, fmt.Errorf("error happend while copying into Lawyer struct: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commiting:%v", err)
	}

	return Lawyer, nil
}

func (pg *PostgresLawyerDB) GetLawyerById(id int64) (*Lawyer, error) {
	Lawyer := &Lawyer{}

	query := `SELECT id, first_name, last_name, email , phone_number FROM lawyer WHERE id = $1 `
	err := pg.DB.QueryRow(query, id).Scan(&Lawyer.ID, &Lawyer.First_name, &Lawyer.Last_name, &Lawyer.Email, &Lawyer.Phone_number)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("nothing found with this id")
	}
	return Lawyer, nil
}

func (pg *PostgresLawyerDB) GetAllLawyers() ([]Lawyer, error) {
	tx, err := pg.DB.Begin()

	defer func() {
		tx.Rollback()
	}()

	if err != nil {
		return nil, fmt.Errorf("error starting the db")
	}
	// query := `SELECT  lawyer.first_name, lawyer.last_name, lawyer.email , lawyer.phone_number ,addresses.country, addresses.city, addresses.street FROM lawyer JOIN addresses ON lawyer.address_id = addresses.address_id`
	query := `SELECT * FROM lawyer
	 JOIN addresses
	 ON lawyer.address_id = addresses.address_id`

	rows, err := tx.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error happend during querying")
	}

	defer rows.Close()

	lawyers := []Lawyer{}

	for rows.Next() {
		var l Lawyer
		err := rows.Scan(&l.ID, &l.First_name, &l.Last_name, &l.Email, &l.Phone_number, &l.Address.ID, &l.AddressID, &l.Address.Street, &l.Address.City, &l.Address.State, &l.Address.Country)

		if err != nil {
			return nil, fmt.Errorf("GetAllLawyers: error scanning row: %w", err)
		}

		lawyers = append(lawyers, l)

	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllLawyers: error after iterating rows: %w", err)
	}

	return lawyers, nil
}

func (pg *PostgresLawyerDB) UpdateLawyer(lawyer *Lawyer) error {
	tx, err := pg.DB.Begin()
	if err != nil {
		fmt.Printf("error starting the db")
	}
	query := `UPDATE lawyer SET first_name =$1, last_name = $2, email =$3, phone_number=$4  WHERE id = $5`

	result, err := tx.Exec(query, lawyer.First_name, lawyer.Last_name, lawyer.Email, lawyer.Phone_number, lawyer.ID)
	if err != nil {
		fmt.Printf("the error occured in executing the query: %v ", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("the was some error in: %v", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if rowsAffected == 0 {
		return fmt.Errorf("error in commiting the changes:%v", err)
	}

	return nil
}

func (pg *PostgresLawyerDB) Deletelawyer(id int64) error {

	query := `DELETE FROM lawyer WHERE id = $1`

	result, err := pg.DB.Exec(query, id)
	if err != nil {
		fmt.Printf("error happended in: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("error happended in: %v", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
