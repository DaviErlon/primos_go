-- Criação do banco de dados
CREATE DATABASE IF NOT EXISTS numeros_primos;

-- Seleciona o banco
USE numeros_primos;

-- Criação da tabela numero para armazenar os numeros primos
CREATE TABLE IF NOT EXISTS numero (
    id INT NOT NULL AUTO_INCREMENT,
    valor INT NOT NULL,
    PRIMARY KEY (id)
);
