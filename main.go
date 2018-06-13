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

	k := 5                        //Максимальное число потоков.
	min := runtime.NumGoroutine() //При запуске программы уже есть некоторые goroutine (Обычно 1)
	var wg sync.WaitGroup         //Группа ожидания выполнения
	totalcount := 0               //Всего вхождение "Go"

	//Чтение входящих данных
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')       //Чтение строки
		text = strings.Replace(text, "\n", "", -1) //Удаление перехода из строки

		for runtime.NumGoroutine()-min >= k { // Ждём пока будет можно запустить новый поток runtime.NumGoroutine() < k
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

			defer wg.Done() //Завершение ожидания одного goroutine
		}(text)
		//Выход из цикла при прочтении всего файла
		if err == io.EOF {
			break
		}
	}

	//Ожидание выполнения всех goroutine
	wg.Wait()

	//Итоговый вывод
	fmt.Println("Total: ", totalcount)

}
