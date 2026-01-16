package main

import (
	"Modulo/databaseprim"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// abertura do banco de dados
	password := os.Getenv("MYSQLPASSWORD")
	port := os.Getenv("MYSQLPORT")

	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:%s)/numeros_primos", password, port)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// coleta dos argumentos
	calculate := flag.Bool("calc", false, "")
	show := flag.Int64("show", 0, "")
	indice := flag.Int64("i", 0, "")
	densi := flag.Int64("densi", 0, "")
	count := flag.Bool("count", false, "")
	t := flag.Int64("timeml", 0, "")

	flag.Parse()

	// interface para comunicação com os primos no banco de dados
	repo := databaseprim.NewRepo(db)
	
	// flag para retorno de um primo por um índice específico
	if *indice > 0 {
		p, _ := repo.GetPrim(*indice)
		fmt.Println("Primo de índice ", *indice, ": ", p)
	}
	
	// flag que exibe a quantidade de primos dentro do banco de dados
	if *count {
		c, _ := repo.GetCountPrim()
		fmt.Println("Quantidade de primos no BD: ", c)
	}

	// flag que exibe a densidade dos primos até um número específico
	if *densi > 2 {
		d, err := repo.GetDensiPrim(*densi) 
		if err != nil {
			panic(err)
		}

		fmt.Println("Densidade de primos até", *densi, ": ", d * 100, "%")
	}

	// flag que exibe os primos até o índice escolhido com um intervalo de tempo em ml
	// passado pelo flag timeml
	if *show > 0 {

		list, _ := repo.GetAtePrim(*show)

		i := int64(1)
		for p := range list.IterPrim() {
			fmt.Println(i,"° Primo: ", p)
			i++
			if *t > 0 {
				time.Sleep(time.Millisecond * time.Duration(*t))
			}
		}

		return
	}

	// caso não receba a flag de calculo o programa deve terminar aqui
	if !*calculate {
		return
	}

	// recuperar o contexto ja calculado no banco de dados
	mp, _ := repo.GetMaxPrim()
	if mp == 1 {
		// se mp eh 1 entao a lista sairá vazia, 
		// então já adicionamos o primeiro primo para que isso nao aconteça
		repo.SetNewPrim(2)
	}
	
	list, _ := repo.GetAllPrim()
	
	// criar um canal para não bloquear os calculos
	c := make(chan int64, 1024)
	defer close(c)

	go func() {
		for v := range c {
			fmt.Println("Novo primo: ", v)
			repo.SetNewPrim(v)
		}
	}()

	// inicio dos calculos
	for i := mp + 2; ; i++ {
		
		flag := true

		for p := range list.IterPrim() {
			
			// primo
			if p*p > i {
				break
			}

			// nao primo
			if i % p == 0 {
				flag = false
				break
			}
		}

		// gatilhos de ser primo
		if flag {
			c <- i
			list.InsertEnd(i)
		}
	}
}
