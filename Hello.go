package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const nMonitoramento = 5
const nTempo = 5

func main() {
	exibeIntroducao()
	for {
		exibeMenu()
		comando := lerComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
		fmt.Println("=============================")
	}
}

func exibeIntroducao() {
	nome := "Robsu"
	fmt.Println("ola !!", nome)
}

func exibeMenu() {
	fmt.Println("<=============================>")
	fmt.Println(" 1 - Iniciar monitoramento")
	fmt.Println(" 2 - Exibir Logs")
	fmt.Println(" 0 - Sair do Programa")
	fmt.Println("<=============================>")
}

func lerComando() int {
	var comando int
	fmt.Scanf("%d", &comando)
	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando Monitoramento")
	sites := lerSitesArquivo()
	for i := 0; i < nMonitoramento; i++ {
		for _, site := range sites {
			testaSite(site)
		}
		fmt.Println("=============================")
		time.Sleep(nTempo * time.Second)
	}
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		trataErr(err)
	}

	fmt.Println(string(arquivo))
}

func testaSite(site string) {
	response, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	switch response.StatusCode {
	case 200:
		fmt.Println("A url : ", site, " foi carregada com sucesso")
		registraLog(site, response.StatusCode, true)
	case 404:
		fmt.Println("Acesso negado a Url : ", site)
		registraLog(site, response.StatusCode, false)
	case 500:
		fmt.Println("Ocorreu um erro ao carregar a Url : ", site)
		registraLog(site, response.StatusCode, false)
	default:
		fmt.Println("Recebido status : ", response.StatusCode, " para o site : ", site)
		registraLog(site, response.StatusCode, false)
	}
}

func lerSitesArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		trataErr(err)
	}
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				trataErr(err)
			}
		}
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
	}
	arquivo.Close()
	return sites
}

func trataErr(err error) {
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		log.Fatal(err)
		os.Exit(-1)
	}
}

func registraLog(site string, typeStatus int, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		trataErr(err)
	}
	arquivo.WriteString(time.Now().Format("02/01/2006") + "- Site : " + site + " Status : " + strconv.FormatBool(status) + " Responta = " + strconv.Itoa(typeStatus) + "\n")

	arquivo.Close()
}
