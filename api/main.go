package handler

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scrap() (float64, error) {
	url := "https://alsuper.com/producto/pulpa-de-res-picada-357825"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code error: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, err
	}

	var priceText string
	doc.Find(".as-font-24").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "$") {
			priceText = text
			return false // break loop
		}
		return true
	})

	// Extract number from price string, e.g. "$139.90" → 139.90
	re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)
	match := re.FindString(priceText)
	if match == "" {
		return 0, fmt.Errorf("no price found")
	}

	price, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

var content embed.FS
var tmpl *template.Template

//var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func Handler(w http.ResponseWriter, r *http.Request) {
	// Parse templates from the embedded filesystem
	var err error
	tmpl, err = template.ParseFS(content, "../templates/index.html")
	if err != nil {
		// Log the error or panic; in a serverless function, panicking here
		// means the function won't start correctly.
		panic(err)
	}

	http.Handle("/public/sounds/", http.StripPrefix("/public/sounds/", http.FileServer(http.Dir("sounds"))))

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
	price, _ := Scrap()

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
