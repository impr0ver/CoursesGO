/*
## Задача № 1
Написать API для указанных маршрутов(endpoints)
"/info"   // Информация об API
"/first"  // Случайное число
"/second" // Случайное число
"/add"    // Сумма двух случайных чисел
"/sub"    // Разность
"/mul"    // Произведение
"/div"    // Деление

*результат вернуть в виде JSON

"math/rand"
number := rand.Intn(100)
! не забудьте про Seed() // SEED IS DEPRECATED!!!


GET http://127.0.0.1:1234/first

GET http://127.0.0.1:1234/second

GET http://127.0.0.1:1234/add
GET http://127.0.0.1:1234/sub
GET http://127.0.0.1:1234/mul
GET http://127.0.0.1:1234/div
GET http://127.0.0.1:1234/info
*/

package main

import (
	//"encoding/json"

	calculator "calcOnRestApi/internal"
	"fmt"
	"log"
	"net/http"
)

type calcServer struct {
	store *calculator.Number
}

func NewCalcServer() *calcServer {
	store := calculator.New()
	return &calcServer{store: store}
}

func (cs *calcServer) calcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Path {
		case "/first":
			first := cs.store.CreateFirstNumber()
			fmt.Fprintf(w, "First number: %d", first)
		case "/second":
			second := cs.store.CreateSecondNumber()
			fmt.Fprintf(w, "Second number: %d", second)
		case "/add":
			sum := cs.store.SumNumbers(cs.store)
			b, _ := cs.store.GetJSONFormat(cs.store)
			fmt.Fprintf(w, "Sum number: %d\r\n%s", sum, string(b))
		case "/sub":
			sub := cs.store.SubNumbers(cs.store)
			b, _ := cs.store.GetJSONFormat(cs.store)
			fmt.Fprintf(w, "Sub number: %d\r\n%s", sub, string(b))
		case "/mul":
			mul := cs.store.MulNumbers(cs.store)
			b, _ := cs.store.GetJSONFormat(cs.store)
			fmt.Fprintf(w, "Mul number: %d\r\n%s", mul, string(b))
		case "/div":
			div, err := cs.store.DivNumbers(cs.store)
			b, _ := cs.store.GetJSONFormat(cs.store)
			if err != nil {
				fmt.Fprintf(w, "err: divison by zero")
			} else {
				fmt.Fprintf(w, "Div number: %d\r\n%s", div, string(b))
			}
		case "/info":
			fmt.Fprintf(w, "Help page: calculator on REST API\r\n\r\n\"/first\" - generate first random number\n\"/second\" - generate second random number\n\"/add, /sub, /mul, /div,\" - sum, difference, multiple andf division random numbers")
		default:
			fmt.Fprintf(w, r.URL.Path+" is not supported")
		}
	} else {
		fmt.Fprintf(w, "Unsupported HTTP method!\r\n")
	}

}

func main() {
	mux := http.NewServeMux()
	server := NewCalcServer()
	// Added routing for "/task/" URL
	mux.HandleFunc("/", server.calcHandler)

	log.Fatal(http.ListenAndServe(":1234", mux))
}
