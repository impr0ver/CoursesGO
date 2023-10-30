package calculator

import (
	//"fmt"

	"fmt"
	"math/rand"
	"sync"
	"time"
	"encoding/json"

	//"golang.org/x/text/number"
)

/*type Task struct {
	Id	int `json:"id"`
	Text string `json:"text"`
	Tags []string `json:"tags"`
	Due time.Time `json:"due"`
}*/

type Number struct {
	sync.Mutex
	First  int
	Second int
	Result int
}

func New() *Number {
	num := &Number{}

	return num
}

// Create first number
func (num *Number) CreateFirstNumber() int {
	num.Lock()
	defer num.Unlock()

	//rand.Seed(20) rand.Seed is deprecated
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	num.First = r1.Intn(1000)
	return num.First
}

func (num *Number) CreateSecondNumber() int {
	num.Lock()
	defer num.Unlock()

	//rand.Seed(20) rand.Seed is deprecated
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	num.Second = r1.Intn(1000)

	return num.Second
}

func (num *Number) SumNumbers(number *Number) int {
	num.Lock()
	defer num.Unlock()

	number.Result = number.First + number.Second

	return number.Result
}

func (num *Number) SubNumbers(number *Number) int {
	num.Lock()
	defer num.Unlock()

	number.Result =  number.First - number.Second

	return number.Result
}

func (num *Number) MulNumbers(number *Number) int {
	num.Lock()
	defer num.Unlock()

	number.Result =  number.First * number.Second

	return number.Result
}

func (num *Number) DivNumbers(number *Number) (int, error) {
	num.Lock()
	defer num.Unlock()
	var div float32
	if number.Second != 0 {
		div = float32(number.First) / float32(number.Second)
	} else {
		return 0, fmt.Errorf("err: division by %d", num.Second)
	}
	return int(div), nil

}

func (num *Number) GetJSONFormat(number *Number) ([]byte, error){
	num.Lock()
	defer num.Unlock()

    bytes, err := json.Marshal(number)
    if err != nil {
        fmt.Printf("Error: %s", err)
        return nil, err
    }
    //fmt.Println(string(b))
	return bytes, nil
}
