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
	usuario := os.Getenv("MYSQLUSER")
	password := os.Getenv("MYSQLPASSWORD")
	port := os.Getenv("MYSQLPORT")

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/numeros_primos", usuario, password, port)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// coleta dos argumentos
	calculate := flag.Bool("calc", false, "")
	show := flag.Int64("show", 0, "")
	index := flag.Int64("i", 0, "")
	densi := flag.Int64("densi", 0, "")
	count := flag.Bool("count", false, "")
	timems := flag.Int64("timems", 0, "")

	flag.Parse()

	// interface para comunicação com os primos no banco de dados
	repo := databaseprim.NewRepo(db)
	
	// flag para retorno de um primo por um índice específico
	if *index > 0 {
		p, err := repo.GetPrim(*index)
		if err != nil {
			panic(err)
		}

		fmt.Println("Primo de índice ", *index, ": ", p)
	}
	
	// flag que exibe a quantidade de primos dentro do banco de dados
	if *count {
		c, err := repo.GetCountPrim()
		if err != nil {
			panic(err)
		}
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

	// flag que exibe os primos até o índice escolhido com um intervalo de tempo em ms
	// passado pelo flag timems
	if *show > 0 {

		list, err := repo.GetAtePrim(*show)
		if err != nil {
			panic(err)
		}

		i := int64(1)
		for p := range list.IterPrim() {
			fmt.Println(i,"° Primo: ", p)
			i++
			if *timems > 0 {
				time.Sleep(time.Millisecond * time.Duration(*timems))
			}
		}

		return
	}

	// caso não receba a flag de calculo o programa deve terminar aqui
	if !*calculate {
		return
	}

	// recuperar o contexto ja calculado no banco de dados
	mp, err := repo.GetMaxPrim()
	if err != nil {
		panic(err)
	}

	if mp == 1 {
		// se mp eh 1 entao a lista sairá vazia, 
		// então já adicionamos o primeiro primo para que isso nao aconteça
		repo.SetNewPrim(2)
	}
	
	list, err := repo.GetAllPrim()
	if err != nil {
		panic(err)
	}
	
	// criar um canal para não bloquear os calculos
	c := make(chan int64, 1024)
	defer close(c)

	go func() {
		for v := range c {
			fmt.Println("Novo primo: ", v)
			repo.SetNewPrim(v)
		}
	}()

	// inicio dos cálculos, pulando de 2 em dois para navegar pelos ímpares
	for i := mp + 2; ; i++ {
		
		// flag interna
		isPrime := true

		for p := range list.IterPrim() {
			
			// é primo, quebra de loop antecipada
			if p*p > i {
				break
			}

			// nao primo
			if i % p == 0 {
				isPrime = false
				break
			}
		}

		// gatilhos de ser primo
		if isPrime {
			c <- i
			list.InsertEnd(i)
		}
	}
}
