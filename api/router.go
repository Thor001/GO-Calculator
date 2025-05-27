package main

import (
	"html/template"
	"main/app"
	"net/http"
	"strconv"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func Handler(w http.ResponseWriter, r *http.Request) {
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/sounds/", http.StripPrefix("/sounds/", http.FileServer(http.Dir("sounds"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculate", calculateHandler)
	http.ListenAndServe(":8080", nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.FormValue("inputA")
	bStr := r.FormValue("inputB")

	a, errA := strconv.ParseFloat(aStr, 64)
	b, errB := strconv.ParseFloat(bStr, 64)
	if errA != nil || errB != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	c := (b * 0.55) * a
	d := (b * 0.075) * a
	e := (b * 0.175) * a
	f := (b * 0.125) * a
	g := (b * 0.075) * a
	h := (b * 0.08) * a
	ha := h / 200
	i := (b * a) / 1000
	price, _ := app.Scrap()

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(
		"<p> Res picada  = " + format(c) + " gr " + "| Precio: " + format(price) + "/kg | Total: $" + format(price*c/1000) + "</p>" +
			"<p> Tocino picado  = " + format(d) + " gr" + "</p>" +
			"<p> Jamon en cuadros  = " + format(e) + " gr" + "</p>" +
			"<p> Salchicha para asar  = " + format(f) + " gr" + "</p>" +
			"<p> Chorizo   = " + format(g) + "</p>" +
			"<p> Cebolla  = " + format(h) + " gr" + " o Cebolla grande: " + format(ha) + "</p>" +
			"<p> Cerveza  = " + format(i) + " Latas</p>" +
			"<p> Pure de tomate / clamato / jugo V8 al gusto </p>",
	))
}

func format(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
