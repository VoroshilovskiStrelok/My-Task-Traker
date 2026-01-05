# Makefile для Task Tracker CLI

.PHONY: build test clean run help

## build: Собрать исполняемый файл
build:
	go build -o task-cli main.go

## test: Запустить все тесты с покрытием
test:
	go test ./... -v -cover

## clean: Удалить бинарник и файл задач (для сброса)
clean:
	rm -f task-cli task-cli.exe tasks.json

## run: Собрать и сразу запустить список задач
run: build
	./task-cli list

## help: Показать эту справку
help:
	@echo "Доступные команды:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'