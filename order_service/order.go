package main

import (
    "net/http"
    "html/template"
)

type Order struct {
    Address string
    Phone   string
}

var orders []Order

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/order", orderHandler)
    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        http.ServeFile(w, r, "frontend/order.html") // Zaktualizowana ścieżka
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        address := r.FormValue("address")
        phone := r.FormValue("phone")
        if address == "" || phone == "" {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }
        orders = append(orders, Order{Address: address, Phone: phone})
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        http.Error(w, "Not found", http.StatusNotFound)
    }
}