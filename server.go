package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

func main() {
	stripe.Key = "your_stripe_key"
	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)
	http.HandleFunc("/health", handleHealth)

	log.Println("listening on port 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProductId string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address1  string `json:"address_1"`
		Address2  string `json:"address_2"`
		City      string `json:"city"`
		State     string `json:"state"`
		Country   string `json:"country"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := &stripe.PaymentIntentParams{
		Amount:                  stripe.Int64(calculateProductAmount(req.ProductId)),
		Currency:                stripe.String(string(stripe.CurrencyJPY)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{Enabled: stripe.Bool(true)},
	}

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response struct {
		ClientSecret string `json:"clientSecret"`
	}

	response.ClientSecret = paymentIntent.ClientSecret

	var byt bytes.Buffer
	err = json.NewEncoder(&byt).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, &byt)
	if err != nil {
		fmt.Println(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := []byte("Server is up and running")
	_, err := w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}

func calculateProductAmount(productId string) int64 {
	switch productId {
	case "Forever Pants":
		return 26000
	case "Forever Shirt":
		return 15000
	case "Forever Shorts":
		return 30000
	}
	return 0
}
