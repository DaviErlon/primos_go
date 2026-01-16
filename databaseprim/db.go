package databaseprim

import (
	"Modulo/utilprim"
	"database/sql"
	"errors"
)

type Repo struct {
	DB *sql.DB
}

// Construtor do repositório
func NewRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

// Insere um novo valor primo
func (r *Repo) SetNewPrim(p int64) error {
	query := "INSERT INTO numero (valor) VALUES (?)"
	_, err := r.DB.Exec(query, p)
	return err
}

// Retorna o maior valor da tabela
func (r *Repo) GetMaxPrim() (int64, error) {
	query := "SELECT MAX(valor) FROM numero"

	var max sql.NullInt64
	if err := r.DB.QueryRow(query).Scan(&max); err != nil {
		return 0, err
	}

	if !max.Valid {
		return 0, nil
	}

	return max.Int64, nil
}

// Retorna a quantidade de registros da tabela
func (r *Repo) GetCountPrim() (int64, error) {
	query := "SELECT COUNT(*) FROM numero"

	var count int64
	if err := r.DB.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// Calcula a densidade: limite / quantidade
func (r *Repo) GetDensiPrim(l int64) (float64, error) {

	countall, _ := r.GetCountPrim()

	// limites para o argumento
	if l < 2  || l > countall{
		return 0, errors.New("Não é possível calcular a densidade")
	}

	query := "SELECT COUNT(*) FROM numero WHERE numero.valor <= ?"

	var count int64
	if err := r.DB.QueryRow(query, l).Scan(&count); err != nil {
		return 0, err
	}

	return float64(count) / float64(l), nil
}

// Gera uma lista até o primo de indice l
func (r *Repo) GetAtePrim(l int64) (*utilprim.List[int64], error)  {
	
	countall, _ := r.GetCountPrim()

	// limites
	if l < 0  || l > countall{
		return nil, errors.New("Não é possível calcular a densidade")
	}
	
	query := "SELECT valor FROM numero WHERE numero.id <= ? ORDER BY numero.id"

	rows, err := r.DB.Query(query, l)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := utilprim.NewList[int64]()

	// Gera uma lista de primos
	for rows.Next() {
		var p int64

		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}

		list.InsertEnd(p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// Busca por um primo a partir de um índice
func (r *Repo) GetPrim(i int64) (int64, error)  {
	query := "SELECT valor FROM numero WHERE numero.id = ?"

	var p int64
	if err := r.DB.QueryRow(query, i).Scan(&p); err != nil {
		return 0, err
	}
	
	return p, nil 
}

// Calcula uma lista com todos os primos do BD
func (r *Repo) GetAllPrim() (*utilprim.List[int64], error)  {
	query := "SELECT valor FROM numero ORDER BY numero.id"

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := utilprim.NewList[int64]()

	for rows.Next() {
		var p int64

		err := rows.Scan(&p)
		if err != nil {
			return nil, err
		}

		list.InsertEnd(p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}