package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
)

type BancoImagens struct{}

type Args struct {
	Arquivo []byte
	Nome    string
}

type Vazio struct{}

func (b *BancoImagens) ListaImagens(argumentos *Vazio, resposta *[]string) error {
	files, err := os.ReadDir("imgs/")
	if err != nil {
		log.Println(err)
		return err
	}

	var lista []string
	for _, file := range files {
		if !file.IsDir() {
			lista = append(lista, file.Name())
		}
	}
	*resposta = lista
	return nil
}

func (b *BancoImagens) RecebeImagens(argumentos *Args, resposta *Vazio) error {
	caminho := filepath.Join("imgs", argumentos.Nome)

	err := os.WriteFile(caminho, argumentos.Arquivo, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar imagem:", err)
		return err
	}
	return nil
}

func (b *BancoImagens) EnviaImagens(argumentos *Args, resposta *Args) error {
	caminho := filepath.Join("imgs", argumentos.Nome)

	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		fmt.Println("Erro ao ler imagem:", err)
		return err
	}

	resposta.Arquivo = conteudo
	resposta.Nome = argumentos.Nome
	return nil
}

func main() {

	bancoImg := new(BancoImagens)
	rpc.Register(bancoImg)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Erro no listen:", err)
	}

	defer listener.Close()

	fmt.Println("Servidor iniciado na porta 8080")
	rpc.Accept(listener)
}
