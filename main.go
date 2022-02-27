package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// kullanacağımız yapının struct ı
type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

// veri tabanı olmadığı için verileri geçici olarak dizizde tutmak için dizi oluşturuyoruz
type allEvents []event

// ilk başlangıçta diziye bir eleman ekleme
var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func main() {

	// routur oluşturuyoruz bu sayede gelen isteklere cevap veriyoruz
	router := mux.NewRouter().StrictSlash(true)

	// gelen isteklerin nasıl olaağını(endpointleri belirliyoruz)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	//eğer hata olursa logla diyerek port açıyoruz ve isteklerimizi bu port üzerinden yapıyoruz
	log.Fatal(http.ListenAndServe(":8080", router))

}

//oluşturma işlemleri
func createEvent(w http.ResponseWriter, r *http.Request) {
	// diziye ekleyeceğimiz eventi oluşturuyoruz gelen verileri buna parse edicez
	var newEvent event

	//gelen verileri okuyup reqBody e atıyoz
	reqBody, err := ioutil.ReadAll(r.Body)

	// hata varse diye bakıyoruz
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	// gelen body i oluşturduğumuz newEvent a parse ettik
	json.Unmarshal(reqBody, &newEvent)

	// veri tabanı yerine diziye ekledik
	events = append(events, newEvent)

	// status oluşturduk 200 ok gibi
	w.WriteHeader(http.StatusCreated)
	// ve eklenen veriyi json formatında geri yansıttık
	json.NewEncoder(w).Encode(newEvent)
}

// verilen id ye uygun olanı getirme
func getOneEvent(w http.ResponseWriter, r *http.Request) {

	// grilen id ye uygun olan varsa göstereceğiz
	eventID := mux.Vars(r)["id"]

	// diziyi kontrol ediyoruz girilen id de eleman var mı diye
	for _, singleEvent := range events {
		// girilen id dizideki id ile eşit mi diye kontrol ediyoruz
		if singleEvent.ID == eventID {
			//varsa onu yine json formatında dönüyoruz
			json.NewEncoder(w).Encode(singleEvent)
		} else {
			// yoksa da mesaj bırakıyoruz
			json.NewEncoder(w).Encode("Böyle bir veri yok")
		}
	}
}

//hepsinin getirme
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	//tüm verileri json formatında dön diyoruz
	json.NewEncoder(w).Encode(events)
}

//güncelleme işlemleri
func updateEvent(w http.ResponseWriter, r *http.Request) {
	// güncellme istediğimiz verinin id sini alıyoruz
	eventID := mux.Vars(r)["id"]

	//güncellenen veriyi tutacağımız değişken
	var updatedEvent event

	// girilen bdy in doğrumu diye kontrol yapıyoruz
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	// girilen değerleri oluşturduğumuz değişkene atıyoruz
	json.Unmarshal(reqBody, &updatedEvent)

	// dizi içinde dönerek değişmesini istediğimiz veriyi buluyoruz id leri karşılaştırarak
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			// buluduğumuzda ise eski değerleri girilen değerle değiştiriyoruz
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)

			//en son ise güncel veriyi geri dönüyoruz evet tabiki json formatında
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

//silme işlemleri
func deleteEvent(w http.ResponseWriter, r *http.Request) {

	// silmek istediğimiz verinin id sini çekiyoruz
	eventID := mux.Vars(r)["id"]

	//dizi kontrlü id si aldığımzı id ye eşit olan veriyi bulmak için
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {

			// silme işlemi
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		} else {
			json.NewEncoder(w).Encode("Böyle bir veri bulunamadı")
		}

	}
}
