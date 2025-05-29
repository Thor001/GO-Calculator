package main

import (
	"html/template"
	"main/app"
	"net/http"
	"strconv"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
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
	priceres, _ := app.Scrap("https://alsuper.com/producto/pulpa-de-res-picada-357825")
	pricetoc, _ := app.Scrap("https://alsuper.com/producto/tocineta-413218")
	pricejam, _ := app.Scrap("https://alsuper.com/producto/jamon-de-pierna-5130264")
	pricesal, _ := app.Scrap("https://alsuper.com/producto/salchicha-para-asar-238828")   //800gr
	pricecho, _ := app.Scrap("https://alsuper.com/producto/chorizo-319544")               //100gr
	priceceb, _ := app.Scrap("https://alsuper.com/producto/cebolla-blanca-924")           //250gr cu
	pricecer, _ := app.Scrap("https://alsuper.com/producto/cerveza-six-pack-lata-323328") //six

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(
		"<p> Total Aproximado:" + format(priceres*c/1000+pricetoc*d/1000+pricejam*e/1000+pricesal*f/800+pricecho*g/100+priceceb*h/1000+pricecer*i/6) + "</p>" +

			"<p> Res picada  = " + format(c) + " gr " + "| Precio: " + format(priceres) + "/kg | Total: $" + format(priceres*c/1000) + "</p>" +
			"<p> Tocino picado  = " + format(d) + " gr" + "| Precio: " + format(pricetoc) + "/kg | Total: $" + format(pricetoc*d/1000) + "</p>" +
			"<p> Jamon en cuadros  = " + format(e) + " gr" + "| Precio: " + format(pricejam) + "/kg | Total: $" + format(pricejam*e/1000) + "</p>" +
			"<p> Salchicha para asar  = " + format(f) + " gr" + "| Precio: " + format(pricesal) + "/800gr | Total: $" + format(pricesal*f/800) + "</p>" +
			"<p> Chorizo   = " + format(g) + "| Precio: " + format(pricecho) + "/100gr| Total: $" + format(pricecho*g/100) + "</p>" +
			"<p> Cebolla  = " + format(h) + " gr" + " o Cebolla grande: " + format(ha) + "| Precio: " + format(priceceb) + "/kg | Total: $" + format(priceceb*h/1000) + "</p>" +
			"<p> Cerveza  = " + format(i) + " Latas" + "| Precio: " + format(pricecer) + "/Six | Total: $" + format(pricecer*i/6) + "</p>" +
			"<p> Pure de tomate / clamato / jugo V8 al gusto </p>",
	))
}

func format(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
