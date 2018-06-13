package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

func main() {
	//Максимальное число потоков.
	k := 5
	var urls []string

	//Чтение входящих данных
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')       //Чтение строки
		text = strings.Replace(text, "\n", "", -1) //Удаление перехода из строки
		urls = append(urls, text)                  //Добавление строки в массив с url

		//Выход из цикла при прочтении всего файла
		if err == io.EOF {
			break
		}
	}

	min := runtime.NumGoroutine() //При запуске программы уже есть некоторые goroutine
	var wg sync.WaitGroup         //Группа ожидания выполнения
	totalcount := 0               //Всего вхождение "Go"

	for _, url := range urls { //Перебор массива urls
		for runtime.NumGoroutine()-min >= k { // Ждём пока будет можно запустить новый поток
		}
		wg.Add(1) //Перед запуском потока добавляем в группу ожидания
		//Запуск потока. Можно вынести в отдельную функцию. но тут это не нужно.
		go func(url string) {
			//Стандартная функция https://golang.org/pkg/net/http/
			//Пример взят оттуда

			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			count := strings.Count(string(body), "Go")  //Кол-во вхождений подстроки Go в теле ответа
			totalcount += count                         //Общее количество вхождений подстроки Go
			fmt.Println("Count for ", url, ": ", count) //Вывод на экран информации по текущему url
			defer wg.Done()                             //Завершение ожидания одного goroutine
		}(url)
	}

	//Ожидание выполнения всех потоков
	wg.Wait()

	//Итоговый вывод
	fmt.Println("Total: ", totalcount)

}
