## Algoritmo para cálculo de Números Primos

Um número primo é aquele número divisível apenas por dois outros números: 1 e ele mesmo. Apesar de já ser conhecido há milhares de anos que esse conjunto de números é infinito, sua dispersão até hoje é dita como caótica, apesar de ser computável por algoritmos como este.

Este algoritmo é somente uma otimização de um modelo trivial de fatoração. Em suma: navegar pelos números naturais consecutivos construindo um conjunto de primos, testando sua divisibilidade pelos elementos do conjunto; caso seja primo, adicionar à lista, senão avançar para o próximo natural.  
A otimização vem em duas partes: não navegar por todos os naturais, apenas os ímpares, e testar sua divisibilidade até o maior primo cujo quadrado é menor que o número testado, além de persistir os dados em uma tabela no MySQL, cujo script se encontra em `scriptsql/`.

---

## Guia para uso

- É necessário possuir Golang 1.23 ou superior, pois ele faz uso do pacote "iter";  
- É necessário possuir o MySQL rodando, ou outro gerenciador, desde que se corrija a sintaxe das queries e suas dependências;  
- Inserir nas variáveis de ambiente `MYSQLUSER`, `MYSQLPASSWORD` e `MYSQLPORT` suas respectivas informações referentes ao banco de dados, ou adaptar o código para já ser compilado com essas informações.

Para configuração do banco de dados:

```bash
mysql -u usuario -p < scriptsql/script.sql
```

Em seguida, para compilação da aplicação:
```bash
  go mod tidy
  go build -o primos
```

Agora, pode-se executar o programa fazendo uso de suas flags:
```bash
  ./primos -calc                    \\ Começa a calcular e imprimir na tela os novos primos e também a guardar no banco de dados
  ./primos -show=10                 \\ Imprime os 10 primeiros primos
  ./primos -show=12 -timems=1000    \\ Imprime os 12 primeiros primos com um intervalo de 1s (1000ms)
  ./primos -i=20                    \\ Imprime o primo de índice 20
  ./primos -densi=80                \\ Imprime a densidade dos primos até o natural 80                
  ./primos -count                   \\ Imprime a quantidade de primos no banco de dados
```