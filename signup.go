package main
import(
	"database/sql"
	_"database/sql/driver/mysql"
	"net/http"
	"log"
	"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	//"net/url"
	"html/template"
	"path"
	"os"
	//"strings"
)

var db *sql.DB 
var err error

func main(){
	db, err = sql.Open("mysql","root:system@tcp(127.0.0.1:3306)/godatabase")                //database connection
    
    if err != nil {
        panic(err.Error()) 
        defer db.Close() 
    }
    
    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

    http.HandleFunc("/", LoginHandler)
    http.HandleFunc("/signup.html", SignupHandler)
    //http.HandleFunc("/activate-account", ActivationHandler)
    http.HandleFunc("/profile.html", ProfileHandler)
    http.HandleFunc("/profile-save", ProfileSaveHandler)
    http.HandleFunc("/logout", LogoutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}


func SignupHandler(w http.ResponseWriter, r *http.Request){
	if r.Method !="POST"{
		http.ServeFile(w,r,"signup.html")
		return
	}
//insertion of signup form values to the database table
	Fname:= r.FormValue("fname")
	Lname:=  r.FormValue("lname")
	Gender:=   r.FormValue("gender")
	Dob:= r.FormValue("birthdate")
	Email:=    r.FormValue("email_id")
	Password:= r.FormValue("password")

	//password encryption
	hash,_:= bcrypt.GenerateFromPassword([]byte(Password), 13)
	fmt.Println("asd:" +  string(hash) )
	
	_, err1 := db.Exec("INSERT INTO inactive_user(fname,email_id)VALUES(?,?)",Fname,Email)
	if err1 != nil{
		log.Fatal(err1)
	}
	
	_, err2 := db.Exec("INSERT INTO user(fname,lname,gender,dob,email_id,password,reg_date,status)VALUES(?,?,?,?,?,?,?,?)",Fname,Lname,Gender,Dob,Email,string(hash),time.Now(),"inactive")
	if err2 != nil{
		log.Fatal(err2)
	}		
	
	if err1 != nil {
		//
		fmt.Println("<script> alert('User already exist !'); window.location.href='/signup.html';</script>")
		fmt.Fprintf(w,"<b>User Already Exist !! </b>")
		//http.Redirect(w,r,"/signup.html",301)
	}
	
	w.Write([]byte("<script> alert('Signup successfull ! Please Activate Your Account via the link send to your registered Email-Id'); window.location.href='/login';</script>"))
	//Encryption of Registration_id of the user to send as a token in activation link
	var id string
	err :=db.QueryRow("select reg_id from user where email_id = ?", Email).Scan(&id)
	if err != nil {
			log.Fatal(err)
			return
	}
	fmt.Println(id)
	code,_:= bcrypt.GenerateFromPassword([]byte(id), 15)
	uid := string(code)
	fmt.Println(uid)
	e:= [] string{Email}
	send_mail(e,uid)
	
	http.Redirect(w,r,"/login.html",301)
	
}


	

func send_mail(to [] string, uid string){
	auth:=smtp.PlainAuth("","godummytest@gmail.com","gomailtest","smtp.gmail.com")
	if err != nil{
			log.Fatal(err)
			fmt.Println(err)
	}
	
	//Code to send the activation link to the user via email
	activation_link := "http://localhost:8080/activate-account?registration-id="+uid;
	msg:=[]byte("Registration sucessfull !!\r\n"+
				"Please click on the link to activate your account\r\n"+activation_link)
	err3:=smtp.SendMail("smtp.gmail.com:587",auth,"godummytest@gmail.com",to,msg)
	if err3!=nil{
			log.Fatal(err3)
			fmt.Println("smtp error: %s", err3)
	return
	}
}



/*func ActivationHandler(w http.ResponseWriter, r *http.Request){

	activationurl := w.URL.RequestURI()
	urllen := url.Parse(activationurl)
	mail,_ :=url.ParseQuery(urllen.RawQuery)
	codedmail:= mail["registration-id"][0]
	err1 := bcrypt.CompareHashAndPassword([]byte(pass), []byte(Password))


	stmt, err := db.Prepare("update user set status=? where email_id=?")
        if err != nil{
				log.Fatal(err)
				fmt.Println(err)
		}

  	Estring := strings.Join(to, "")
  	fmt.Println("js" , Estring)
	_, err = stmt.Exec("active", Estring)
	if err != nil{
			log.Fatal(err)
			fmt.Println(err)
	}
}*/



func LoginHandler(w http.ResponseWriter, r *http.Request){

if r.Method !="POST"{
		http.ServeFile(w,r,"login.html")
		return
}

Email:= r.FormValue("email_id") 
Password:= r.FormValue("password")

//Validating the credentials of the user
var pass string
var user_status string
row, err := db.Query("select password, status from user where email_id ='"+Email+"'")
fmt.Println("afetr query")
if err != nil {
	log.Fatal(err)
}
defer row.Close()
for row.Next() {
	 
    err := row.Scan(&pass, &user_status);
	if err != nil {
			log.Fatal(err)
		}
	log.Println(pass)
	log.Println(user_status)

}

err1 := bcrypt.CompareHashAndPassword([]byte(pass), []byte(Password))
if err != nil {
	fmt.Println("Password not matched")
	log.Fatal(err1)
	fmt.Fprintf(w,"<script> alert('Incorrect Email-Id or Password');</script>")//;window.location.href='/signup'
	http.ServeFile(w,r,"/login.html")
	//http.Redirect(w,r,"/login.html",301)
		
}else if  (user_status == "active"){

		fmt.Println("User_status: active and password matched")


		cookie_exp :=time.Now().Add(365*24*time.Hour)
		user_cookie := http.Cookie{ Name: "loggedin", Value: Email, Expires: cookie_exp }
		http.SetCookie(w, &user_cookie)
		fmt.Println("cookievalue", user_cookie.Value)

		http.Redirect(w,r,"/profile.html",301)

	}else{

		fmt.Fprintf(w,"<script> alert('Please activate your account');</script>")
		http.ServeFile(w,r,"/login.html")
		//http.Redirect(w,r,"/login",301)
	}
}


func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    
	if r.Method!= "POST"{
			cookie, err := r.Cookie("loggedin")
			fmt.Println("cokiiiiii",cookie.Value)
			if err != nil {
				// User is not logged in and don't have any cookie
				fmt.Println("err", err)
				http.Redirect(w,r,"/login.html",301)
				return
			}

			fmt.Println("cookievalue", cookie)
			lp := path.Join("", "profile.html")


    // Return a 404 if the template doesn't exist
    	info, err := os.Stat(lp)
    	if err != nil {
        	if os.IsNotExist(err) {
            		http.NotFound(w, r)
            	return
        	}
   		  }

   
    	if info.IsDir() {
        		http.NotFound(w, r)
        		return
   		 }

   		templates, err := template.ParseFiles(lp)
    	if err != nil {
       			fmt.Println(err)
        		http.Error(w, "500 Internal Server Error", 500)
        return
    	}

	//Code to display the profile information of the user	
		type UserProfile struct{
	 			FirstName string
	 			LastName  string
	 			Dob string
	 			Gender string
	 			Email_Id string
		}
		data := UserProfile{} 
		
		rows, err := db.Query("select fname,lname,dob,gender,email_id from user where reg_id = ?",cookie.Value)//userck.Value)//user_ck.Value)
		if err != nil {
				log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
				err := rows.Scan(&data.FirstName, &data.LastName, &data.Dob, &data.Gender, &data.Email_Id);
				if err != nil {
						log.Fatal(err)
				}
				log.Println(data.FirstName)
				log.Println(data.Gender)
		}
		w.Header().Set("Content-Type", "text/html")
    	templates.ExecuteTemplate(w, "profile.html",data)
    }
	
}

func ProfileSaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("yes working!!")

		Pfname:= r.FormValue("fn")
		Plname:=  r.FormValue("ln")
		Pdob:= r.FormValue("dob")

		fmt.Println("FNAME:",Pfname)
		fmt.Println("LNAME:",Plname)
		fmt.Println("DATE:",Pdob)

		_, err= db.Exec("update user set fname='"+Pfname+"', lname='"+Plname+"', dob='"+Pdob+"' where reg_id='1'")
	     
       if err != nil{
		log.Fatal("update:error",err)
		fmt.Println(err)
		return
		}
		http.Redirect(w,r,"/profile.html", 301)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		cookieexpire := http.Cookie{Name: "loggedin", Value:" ", Expires: time.Now()}
		http.SetCookie(w, &cookieexpire)
		fmt.Println("cookievalue", cookieexpire.Value)

		http.Redirect(w,r,"/login.html",301)
	}
}
