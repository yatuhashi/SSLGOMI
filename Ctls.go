package main
import (
    "fmt"
    "log"
    "crypto/tls"
    "time"
    "runtime"
    "flag"
)

func loopconnection(name string,i int,loop bool,result chan int) {
    log.SetFlags(log.Lshortfile)
    conf := &tls.Config{
        InsecureSkipVerify: true,
        SessionTicketsDisabled: false,
        ServerName: name,
    }
    for j :=0; j<1 || loop; j++{
            conn, err := tls.Dial("tcp", name, conf)
            if err != nil {
                fmt.Println("error",i)
                log.Println(err)
                result <-1
                continue
            }else{
                result <-0
            //待機時間をつくり同時アクティブセッション数を増やす
            //    time.Sleep(30000 * time.Millisecond)
            }
            conn.Close()
    }
}

func recieving(result chan int,loop bool,num int){
    success := 0
    failed  := 0
    all := 0.0
    s := time.Now()

    for loop {
        r := <-result
        all = all+1
        if r == 0 {
            success = success+ 1
        }
        if int(all)%1000==0 {
            fmt.Println("all:",all,"Suc:",success,"q/s:",1000/time.Now().Sub(s).Seconds())
            s = time.Now()
        }
    }
    for i :=0; i<num; i++{
        r := <-result
        if r == 0 {
            success = success+ 1
        }else{
            failed = failed + 1
        }
    }
    fmt.Println("Success :",success, "Failed : ",failed)
}

func main(){
  fmt.Println("Please wait 5s....")
  ip := flag.String("i", "127.0.0.1:443", "dst server ip address(default 127.0.0.1:443)")
  num := flag.Int("n", 1, "process num")
  loop := flag.Bool("l", false, "loop flag")
  flag.Parse()
  result := make(chan int)

  go recieving(result,*loop,*num)

  t := time.Now()
  for i :=0; i<*num; i++{//通信プロセスの作成
     go loopconnection(*ip,i,*loop,result)
  }
  f := time.Now()
  //通信プロセスの作成時間
  fmt.Println("Create process time: ",f.Sub(t))
  //現在動いてるプロセス数
  fmt.Println("Current Process num: ",runtime.NumGoroutine())
  //mainプロセスを待機させておく
  time.Sleep(3 * time.Second)
  for runtime.NumGoroutine() > 2{
     time.Sleep(1 * time.Second)
     fmt.Println("Current Process num: ",runtime.NumGoroutine())
  }
}


