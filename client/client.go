package main

import (
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
)

type Args struct {
	Arquivo []byte
	Nome    string
}

type Vazio struct{}

func main() {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	parar := false

	for !parar {
		pastaLocal := ""
		pastaLocal = "imgs"
	
		fmt.Println("\nO que deseja fazer?")
		fmt.Println("(1) Listar arquivos do servidor")
		fmt.Println("(2) Listar seus arquivos locais")
		fmt.Println("(3) Receber um arquivo do servidor (Download)")
		fmt.Println("(4) Enviar um arquivo ao servidor (Upload)")
		fmt.Println("(5) Sair")
		fmt.Print("Escolha uma opção: ")

		var valor int
		fmt.Scanln(&valor)
		switch valor {
		case 1:
			var lista []string
			err := client.Call("BancoImagens.ListaImagens", &Vazio{}, &lista)
			if err != nil {
				fmt.Println("Erro ao listar arquivos do servidor:", err)
			} else {
				fmt.Println("-> Arquivos no servidor:", lista)
			}

		case 2:
			files, err := os.ReadDir(pastaLocal)
			if err != nil {
				fmt.Println("Erro ao ler pasta local:", err)
			} else {
				fmt.Println("-> Arquivos locais:")
				for _, f := range files {
					if !f.IsDir() {
						fmt.Println("  -", f.Name())
					}
				}
			}

		case 3:
			fmt.Print("Digite o nome do arquivo para baixar: ")
			var nome string
			fmt.Scanln(&nome)

			req := Args{Nome: nome}
			var resp Args

			err := client.Call("BancoImagens.EnviaImagens", &req, &resp)
			if err != nil {
				fmt.Println("Erro ao receber do servidor:", err)
			} else {
				caminho := filepath.Join(pastaLocal, resp.Nome)
				os.WriteFile(caminho, resp.Arquivo, 0644)
				fmt.Println("-> Arquivo recebido e salvo com sucesso!")
			}

		case 4:
			fmt.Print("Digite o nome do arquivo local para enviar: ")
			var nome string
			fmt.Scanln(&nome)

			caminho := filepath.Join(pastaLocal, nome)
			conteudo, err := os.ReadFile(caminho)
			if err != nil {
				fmt.Println("Erro ao encontrar/ler o arquivo local:", err)
				continue
			}

			req := Args{Nome: nome, Arquivo: conteudo}
			var resp Vazio

			err = client.Call("BancoImagens.RecebeImagens", &req, &resp)
			if err != nil {
				fmt.Println("Erro ao enviar para o servidor:", err)
			} else {
				fmt.Println("-> Arquivo enviado com sucesso!")
			}

		case 5:
			parar = true
			fmt.Println("Saindo...")

		default:
			fmt.Println("-> Opção inválida. Tente novamente.")
		}

	}
}
