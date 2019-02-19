# GoTrader bot

Projeto para criar um robo de trade para a bitmex

[![Go Report Card](https://goreportcard.com/badge/github.com/thiago-scherrer/gotrader)](https://goreportcard.com/report/github.com/thiago-scherrer/gotrader)

![gopher](assets/gopher.png)


## requerimentos

- docker >18.09.1
- bitmex api
- gopkg.in/yaml.v2 (go get gopkg.in/yaml.v2)

## Como funciona

Este robo está em construção ainda. Seu objetivo é automatizar uma regra criada pelo trader. Não é uma máquina de fazer dinheiro...

## Build

Copie o arquivo de exemplo de configuraçao, que está dentro de configs. Altere os dados necessários.
Entre na pasta ./internal/gotrader e execute o teste:

```bash
go test -args config ../../config.yml
```

Em seguida, faça o build:

```bash
go build -o gotrader 
```

## Refs

- [golang-standards](https://github.com/golang-standards/project-layout)
- [bitmex api](https://www.bitmex.com/api/explorer/)
- [goreportcard](https://goreportcard.com/)
- [gopherize](https://gopherize.me)
