package main
import (
    "fmt"
    "log"
    "crypto/tls"
    "time"
    "runtime"
    "flag"
    "io/ioutil"
    "net/http"
//    "io"
//    "golang.org/x/net/http2"
//    "os"
//    "strconv"
)

func HttpsClient(ecdh bool) *http.Transport{
    var tr = &http.Transport{
            TLSHandshakeTimeout: 1 * time.Second,
            MaxIdleConnsPerHost: 65000,
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
    }
    if ecdh {
           tr = &http.Transport{
                   TLSHandshakeTimeout: 1 * time.Second,
                   MaxIdleConnsPerHost: 65000,
                   TLSClientConfig: &tls.Config{
                       InsecureSkipVerify: true,
                       CipherSuites: []uint16 {
                            tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
                       },
                   },
           }
    }

    return tr
}

func connect(URL string,i int,loop ,show ,ecdh bool,result chan int) {
    tr := HttpsClient(ecdh)
    hc := &http.Client{Transport: tr}
    for j :=0; j<1 || loop; j++ {//デフォ１回、または無限ループする
            req, _ := http.NewRequest("GET", URL, nil)
            resp, err := hc.Do(req)
            if err != nil {
                //失敗した時
                log.Println(err)
                result <-1
            }else{
                //成功した時
                result <-0
                //待機時間をつくり同時アクティブセッション数を増やす
                //time.Sleep(30000 * time.Millisecond)
                body, _ := ioutil.ReadAll(resp.Body)
                //io.Copy(ioutil.Discard, resp.Body)
                defer resp.Body.Close()
                if(show){//フラグが有効な時だけ受信内容を表示
                    fmt.Printf("%s\n", resp.Header)
                    fmt.Printf("%s", body)
                }
            }
            tr.CloseIdleConnections()
    }

}

func recieving(result chan int,loop bool,num int){
    //通信結果を受け取るためのプロセス
    success := 0
    failed  := 0
    all := 0.0
    s := time.Now()
    for loop {//無限ループする場合
        r := <-result
        all = all+1
        if r == 0 {
            success = success+ 1
        }
        if int(all)%10000==0 {//結果を１万回ごとに表示
            fmt.Println("all:",all,"Suc:",success,"q/s:",10000/time.Now().Sub(s).Seconds())
            s = time.Now()
        }
    }
    for i :=0; i<num; i++{//指定プロセス数だけの時
        r := <-result
        if r == 0 {
            success = success+ 1
        }else{
            failed = failed + 1
        }
    }
    fmt.Println("RESULT   Success : ",success," Failed : ",failed)
}

func main(){
  fmt.Println("Please wait....")
  URL := flag.String("i", "https://127.0.0.1/", "dst server URL")
  num := flag.Int("n", 1, "process num")
  loop := flag.Bool("l", false, "loop flag")
  ecdh := flag.Bool("e", false, "ecdh flag")
  show := flag.Bool("s", false, "show body flag")
  flag.Parse()
  result := make(chan int)
  //通信結果を受け取るためのプロセス
  go recieving(result,*loop,*num)

  t := time.Now()
  for i:=0; i<*num;i++ {//通信プロセスの作成
     go connect(*URL,i,*loop,*show,*ecdh,result)
  }
  f := time.Now()
  //通信プロセスの作成時間
  fmt.Println("Create process time: ",f.Sub(t))
  //現在動いてるプロセス数
  fmt.Println("Current Process num: ",runtime.NumGoroutine())
  //mainプロセスを待機させておく
  time.Sleep(3 * time.Second)
  for runtime.NumGoroutine() > 4 || *loop{
     time.Sleep(1 * time.Second)
  }
  fmt.Println("Current Process num: ",runtime.NumGoroutine())
}


