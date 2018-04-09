package main

import (
    "fmt"
    "net/http"
    "strings"
    "log"
    "io/ioutil"
    "encoding/json"
    "html/template"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

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

func gethello(w http.ResponseWriter, r *http.Request)  {
    
       if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
        ingredients := r.FormValue("ingredients")
        receipe := r.FormValue("receipe")
        link := r.FormValue("link")
        image:=r.FormValue("image")
        fmt.Fprintf(w,ingredients+receipe+link+image)

         db, err := sql.Open("mysql",
    "root:Welcome1@tcp(127.0.0.1:3306)/receipe")
    if err != nil {
        log.Fatal(err)
      }
      stmt, err := db.Prepare("INSERT INTO favourite(ingredients,receipe,image,link) VALUES(?,?,?,?)")
       if err != nil {
       log.Fatal(err)
      }
      res, err := stmt.Exec(ingredients,receipe,image,link)
      if err != nil {
       log.Fatal(err)
      }

      lastId, err := res.LastInsertId()
if err != nil {
  log.Fatal(err)
}
rowCnt, err := res.RowsAffected()
if err != nil {
  log.Fatal(err)
}
log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
  

        defer db.Close()




}

func getFavourites(w http.ResponseWriter, r *http.Request) {
  var s Result
  var receipeList  []Receipe
  var receipeObj Receipe
  db, err := sql.Open("mysql","root:Welcome1@tcp(127.0.0.1:3306)/receipe")
  if err != nil {
    log.Fatal(err)
  }
  rows, err := db.Query("select ingredients,receipe,image, link from favourite")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()
  for rows.Next() {
      err := rows.Scan(&receipeObj.Ingredients, &receipeObj.Title,&receipeObj.Thumbnail,&receipeObj.Href)
      if err != nil {
        log.Fatal(err)
      }
      receipeList=append(receipeList,receipeObj)
    }
  err = rows.Err()
  if err != nil {
      log.Fatal(err)
  }
  s.Results=receipeList
  tmpl := template.Must(template.ParseFiles("favourite.html"))
  tmpl.Execute(w, s)
  defer db.Close()

}
func receipeFinder(w http.ResponseWriter, r *http.Request) {

    if r.Method=="GET" {
       // fmt.Fprintf(w,"Receipe...!!!!!!!")
       var url="http://www.recipepuppy.com/api/?"
       var ingredients= r.URL.Query().Get("i")
       var rec=r.URL.Query().Get("q")
      // fmt.Fprintf( w,r.URL.Query().Get("i"))
       var finalURL= url+"i="+ingredients+"&q="+rec
         res, err := http.Get(finalURL)
        if err != nil {
        log.Fatal(err)
        }else{
           data, _ := ioutil.ReadAll(res.Body)
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
    http.HandleFunc("/addFavourite",gethello)
    http.HandleFunc("/getFavourites",getFavourites)
    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    
}