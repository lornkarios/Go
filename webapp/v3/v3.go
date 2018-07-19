package main

// Импортируем всё, что нам может понадобиться
import (
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	//	"os"
)

func handler(iWrt http.ResponseWriter, iReq *http.Request) {
	// Отправляет "Привет, мир" в ответ на запрос
	fmt.Fprintln(iWrt, "Привет, мир")
}

func main() {
	// При получени запроса к "/*", если не задано других обработчиков для данного
	// запроса, вызываем функцию "handler".
	http.HandleFunc("/", handler)

	// Ну а как-же без этого?)
	log.Println("Запускаемся. Слушаем порт 8080")

	// Сканируем запросы к порту 8080. При наличии таковых - отвечаем так, как
	// указано выше
	http.ListenAndServe(":8080", nil)
}
