package main
import (
    "fmt"
    "log"
    "crypto/tls"
    "time"
    "runtime"
    "flag"
    "os"
//    "strconv"
    "github.com/garyburd/redigo/redis"
)

func loopconnection(name string,i int) {
    keep_alive := "GET / HTTP/1.1\r\n Connection: keep-alive\r\n"
    message := []byte(keep_alive)
    ip := "10.14.6.35"
    key :=  [32]byte{}
    copy(key[:],ip)

    log.SetFlags(log.Lshortfile)
    conf := &tls.Config{
        InsecureSkipVerify: true,
        SessionTicketsDisabled: false,
        ServerName: name,
        SessionTicketKey: key,
    }
//   sessioncache := tls.NewLRUClientSessionCache(1)
//   clientsessionstate := make([]tls.ClientSessionState, 1)
//   sessionkey := name+strconv.Itoa(i)
//   fmt.Println(sessionkey)
//   sessioncache.Put(sessionkey, &clientsessionstate[i])
//    fmt.Println(sessioncache.Get(sessionkey))
//for {
    conn, err := tls.Dial("tcp", name, conf)
    defer conn.Close()
    if err != nil {
        fmt.Println("error",i)
        log.Println(err)
        ConRedis("failed")
        return
    }else{
        ConRedis("num")
    //待機時間をつくり同時アクティブセッション数を増やす
    //    time.Sleep(30000 * time.Millisecond)
    //    fmt.Println("success",i)       
    }
    fmt.Println(message)
    //送信データ
    n, err := conn.Write(message)
    if err != nil {
        log.Println(n, err)
        return
    }
//}
/*
    buf := make([]byte, 1000)
    n, err = conn.Read(buf)
    if err != nil {
        println("failed recieve")
        log.Println(n, err)
        return
    }
    println("------------------")
    println(string(buf[:n]))
*/

}

func connection(name string,i int) {
    log.SetFlags(log.Lshortfile)
    conf := &tls.Config{
        InsecureSkipVerify: true,
    }
    //セッション開始
    conn, err := tls.Dial("tcp", name, conf)
    if err != nil {
        fmt.Println("error",i)
        log.Println(err)
        ConRedis("failed")
        return
    }else{
        ConRedis("num")
    //待機時間をつくり同時アクティブセッション数を増やす
        time.Sleep(30000 * time.Millisecond)
    //    fmt.Println("success",i)       
    }
    conn.Close()
}

func loop(num int){
    fmt.Println("loop")
    t := time.Now()
    for i :=0; i<num; i++{
       loopconnection("10.14.6.35:443",i)
    }
    f := time.Now()
    fmt.Println(f.Sub(t))
    log.Println(runtime.NumGoroutine())
    time.Sleep(5000 * time.Millisecond)
//    log.Println(runtime.NumGoroutine())
//    time.Sleep(22000 * time.Millisecond)
}

func ConRedis(key string){
    c, err := redis.Dial("tcp","127.0.0.1:6379")
    if err != nil {
          fmt.Println(err)
          c.Close()
          os.Exit(1)
    }

    val, err := redis.Int(c.Do("GET", key))
    if err != nil {
       fmt.Println(err)
       os.Exit(1)
    }
    c.Do("SET", key, val+1)
    c.Close()
}

func main(){
  num := flag.Int("n", 1, "flag 1")
  flag.Parse()

  for j :=0; j<1000; j++{
    loop(*num)
  }
}

