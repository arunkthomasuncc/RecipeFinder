package main

import (
    "fmt"
    "net/http"
    "strings"
    "log"
    "io/ioutil"
    "encoding/json"
    "html/template"

)
type Receipe struct{

     Title string `json:"title"`
     Href string `json:"href"`
     Ingredients string `json:"ingredients"`
     Thumbnail string `json:"thumbnail"`

} 
type Result struct {
    MainTitle string `json:"title"`
    Version float32 `json:"version"`
    Href string `json:"href"`
    Results []Receipe `json:"results"`
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    
    var s Result
    tmpl := template.Must(template.ParseFiles("result.html"))
    tmpl.Execute(w, s)
    r.ParseForm()  // parse arguments, you have to call this by yourself
    fmt.Println(r.Form)  // print form information in server side
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
   // fmt.Fprintf(w, "Welcome to Receipe Finder") // send data to client side
}
func getReceipies(body []byte) (*Result, error) {
    var s = new(Result)
    err := json.Unmarshal(body, &s)
    if(err != nil){
        fmt.Println("whoops:", err)
    }
    fmt.Println(s)
    return s, err
}
func receipeFinder(w http.ResponseWriter, r *http.Request) {

    if r.Method=="GET" {
       // fmt.Fprintf(w,"Receipe...!!!!!!!")
       var url="http://www.recipepuppy.com/api/?"
       var ingredients= r.URL.Query().Get("i")
       var rec=r.URL.Query().Get("q")
      // fmt.Fprintf( w,r.URL.Query().Get("i"))
       var finalURL= url+"i="+ingredients+"&q="+rec
      // fmt.Fprintf(w,finalURL)

       // fmt.Print( r.Form["q"])
        //fmt.Fprintf(w, "ingredients are"+ ingredients[0])
        //fmt.Fprintf(w, "receipe is"+ rec[0])
      // res, err := http.NewRequest(http.MethodGet, finalURL, nil)
         res, err := http.Get(finalURL)
        if err != nil {
        log.Fatal(err)
        }else{
           data, _ := ioutil.ReadAll(res.Body)
          // fmt.Fprintf(w,string(data))
           //fmt.Println("before calling",data)
         //  s, err:=getReceipies([]byte (data))
           var s Result
           json.Unmarshal(data, &s)
           fmt.Println("title is"+s.MainTitle)        
           tmpl := template.Must(template.ParseFiles("result.html"))
            err1 :=tmpl.Execute(w, s)
             if err1 != nil {
                log.Fatal("Execute: ", err1)
                return
        }
           
        }


             
    }
}



func main() {
    http.HandleFunc("/", sayhelloName) // set router
    http.HandleFunc("/Receipe",receipeFinder)

    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}