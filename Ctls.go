package main
import (
    "fmt"
    "log"
    "crypto/tls"
    "time"
    "runtime"
    "flag"
//    "strconv"
)

func loopconnection(name string,i int,loop bool,success ,failed  chan int) {
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
      for loop {
            conn, err := tls.Dial("tcp", name, conf)
            defer conn.Close()
            if err != nil {
                fmt.Println("error",i)
                log.Println(err)
                tmp := <-failed
                failed <- tmp+1
                return
            }else{
                tmp := <-success
                success <-tmp+1
            //待機時間をつくり同時アクティブセッション数を増やす
            //    time.Sleep(30000 * time.Millisecond)
            }
            //送信データ
            n, err := conn.Write(message)
            if err != nil {
                log.Println(n, err)
                return
            }
      }
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

func recieving(success, failed chan int){
    success <-0
    failed  <-0
    for {
        select{//ひたすら受け取る
        case s := <-success:
            if s%2 == 1{
               println(s)
            }
        case f := <-failed:
            if f%2 == 1{
               println(f)
            }
        }
    }
}

func main(){
  num := flag.Int("n", 1, "process num")
  loop := flag.Bool("l", false, "loop")
  flag.Parse()
  success := make(chan int)
  failed  := make(chan int)
  go recieving(success,failed)

  for j :=0; j<1000; j++{
    t := time.Now()
    for i :=0; i<*num; i++{//通信プロセスの作成
       go loopconnection("10.14.6.35:443",i,*loop,success,failed)
    }
    f := time.Now()
    fmt.Println(f.Sub(t))//すべての通信プロセスの作成にかかった時間
    log.Println(runtime.NumGoroutine())
    time.Sleep(2000 * time.Millisecond)//mainプロセスを待機させておく
  }
}


