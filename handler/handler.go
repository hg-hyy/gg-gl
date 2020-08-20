package handler

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}

//Recover middleware
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				//打印错误堆栈信息
				log.Printf("panic: %v\n", r)
				debug.PrintStack()
				c.JSON(http.StatusOK, gin.H{
					"code": "1",
					"msg":  errorToString(r),
					"data": nil,
				})
				//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
				c.Abort()
			}
		}()
		//加载完 defer recover，继续后续接口调用
		c.Next()
	}
}

//Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

// WhiteList middleware
func WhiteList() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 定义ip白名单
		whiteList := []string{
			"127.0.0.1",
		}

		ip := c.ClientIP()

		flag := false

		for _, host := range whiteList {
			if ip == host {
				flag = true
				break
			}
		}

		if !flag {
			c.String(http.StatusNetworkAuthenticationRequired, "your ip is not trusted: %s", ip)
			c.Abort()
		}

	}
}

// Persion struct
type Persion struct {
	Name string `json:"name" form:"username"`
	Age  int    `json:"age" form:"age"`
}

// Study ...
func (p Persion) Study() string {
	fmt.Println("i am study how to make loud with fhh ")
	s := fmt.Sprintf("艺名:%s 年龄:%d", p.Name, p.Age)
	return s
}

// Make ...
func (p *Persion) Make(name string, age int) {
	fmt.Println("i am make loud with fsh")
	p.Name = name
	p.Age = age
}

// Employee ...
type Employee struct {
	ID   string
	Name string
	Age  int
}

// UpdateAge ...
func (e *Employee) UpdateAge(newVal int) {
	e.Age = newVal
}

// UpdateAge1 ...
func (e Employee) UpdateAge1(newVal int) {
	e.Age = newVal
}

// GetAge ...
func (e Employee) GetAge() int {
	return e.Age
}

// TestMethod ...
func TestMethod() {
	e := Employee{"1", "fhh", 10}
	e.UpdateAge(99)
	fmt.Println(e.GetAge())
}

// TestMethod1 ...
func TestMethod1() {
	e := Employee{"2", "fsh", 20}
	e.UpdateAge1(10)
	fmt.Println(e.GetAge())
}

// MyStruct ...
type MyStruct struct {
	N int
}

// Printreflect ...
func Printreflect() {

	n := MyStruct{1}

	// get
	immutable := reflect.ValueOf(n)
	val := immutable.FieldByName("N").Int()
	fmt.Printf("N=%d\n", val) // prints 1

	// set
	mutable := reflect.ValueOf(&n).Elem()
	mutable.FieldByName("N").SetInt(7)
	fmt.Printf("N=%d\n", n.N) // prints 7
}

// Testreflect ...
func Testreflect(i interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:", err)
		}

	}()
	if v.Elem().Kind() == reflect.Int {

		v.Elem().SetInt(999)

	} else if v.Elem().Kind() == reflect.String {
		v.Elem().SetString("hello,golang")
	}
	switch v.Kind() {
	case reflect.Int:
		fmt.Println("Int 类型")
	case reflect.Float32:
		fmt.Println("Float32 类型")
	case reflect.Float64:
		fmt.Println("Float64 类型")
	case reflect.String:
		fmt.Println("String 类型")
	case reflect.Array:
		fmt.Println("Array 类型")
	case reflect.Slice:
		fmt.Println("Slice 类型")
	case reflect.Map:
		fmt.Println("Map 类型")
	case reflect.Ptr:
		fmt.Println("ptr 类型")
	case reflect.Struct:
		fmt.Println("Struct 类型")
	default:
		fmt.Println("未找到匹配的类型")
	}
	// 判断是不是结构体
	if t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Struct {
		fmt.Println("not a struct")
		return
	}
	// 这里有疑问，如果不这样写有错误
	t = t.Elem()

	field0 := t.Field(0)
	field1, ok := t.FieldByName("Age")
	if ok {
		fmt.Println(field0, field1.Name, field1.Tag.Get("json"), field1.Type)
	}
	// 遍历结构体字段
	num := t.NumField()
	for i := 0; i < num; i++ {
		fmt.Println(t.Field(i).Name, t.Field(i).Tag.Get("json"))

	}

	// mt0 := t.Method(0)
	// fmt.Println(mt0)
	mt, ok := t.MethodByName("Make")
	if ok {
		fmt.Println(mt)
		ss := v.MethodByName("Study").Call(nil)
		fmt.Println(ss)

	}
	var param []reflect.Value
	param = append(param, reflect.ValueOf("fsh"))
	param = append(param, reflect.ValueOf(33))

	v.MethodByName("Make").Call(param)

	TestMethod()
	TestMethod1()
	Printreflect()
	goroution()
	testchan(ch)
}

var wg sync.WaitGroup

func printhello(num int) {

	for i := 0; i < 5; i++ {

		fmt.Printf("协程ID:%v ------ 输出：%v\n", num, i)
		time.Sleep(time.Second * 1)
	}

	wg.Done()
}

func goroution() {

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go printhello(i)

	}
	wg.Wait()
}

var ch = make(chan int, 10)

func writechan(ch chan int) {

	for i := 0; i < 10; i++ {

		ch <- i
		time.Sleep(time.Second * 1)
	}
	wg.Done()

}
func readchan(ch chan int) {
	// for i := 0; i < 10; i++ {

	// 	it := <-ch1
	// 	fmt.Println(it)
	// }

	for v := range ch {
		fmt.Println(v)

	}
	wg.Done()
}

func testchan(ch chan int) {
	wg.Add(1)
	go writechan(ch)
	wg.Add(1)
	go readchan(ch)
	wg.Wait()
}

func putnum(intchan chan int) {
	for i := 2; i < 100; i++ {
		intchan <- i
	}
	wg.Done()
}

func isprime(intchan, primechan chan int, flagchan chan bool) {
	for num := range intchan {
		flag := true
		for i := 2; i < num; i++ {
			if num%i == 0 {
				flag = false
				break
			}
		}
		if flag {
			primechan <- num
		}
	}

	flagchan <- true
	wg.Done()
}

func printprime(primechan chan int) {
	for val := range primechan {
		fmt.Println(val)
	}
	wg.Done()
}

func isdone(flagchan chan bool, primechan chan int) {
	for i := 0; i < 8; i++ {
		<-flagchan
	}
	close(primechan)
	wg.Done()
}

// Chantestprime ...
func Chantestprime() {
	start := time.Now().Unix()
	intchan := make(chan int, 1000)
	primechan := make(chan int, 1000)
	flagchan := make(chan bool, 16)
	wg.Add(1)
	go putnum(intchan)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go isprime(intchan, primechan, flagchan)
	}
	wg.Add(1)
	go printprime(primechan)
	wg.Add(1)
	go isdone(flagchan, primechan)

	end := time.Now().Unix()
	fmt.Println(end-start, "---s")
	wg.Wait()
}

var count int = 0

var mutex sync.Mutex

var m = make(map[int]int, 0)

func factorial(num int) {
	mutex.Lock()
	var sum = 1
	for i := 1; i <= num; i++ {
		sum *= i
	}
	m[num] = sum
	fmt.Printf("key=%v value=%v\n", num, sum)
	time.Sleep(time.Millisecond * 10)
	mutex.Unlock()
	wg.Done()
}

// Testlock ...
func Testlock() {
	for i := 1; i <= 40; i++ {
		wg.Add(1)
		go factorial(i)

	}
	wg.Wait()
}

var rwmutex sync.RWMutex

func read() {
	rwmutex.RLock()

	fmt.Println("===this is read===")
	time.Sleep(time.Second * 2)
	rwmutex.RUnlock()

	wg.Done()

}
func write() {
	mutex.Lock()
	fmt.Println("---this is write---")
	time.Sleep(time.Second * 2)
	mutex.Unlock()
	wg.Done()
}

//Readandwrite 读写互斥，读并行，写串行
func Readandwrite() {

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go read()
	}
	wg.Wait()
}
