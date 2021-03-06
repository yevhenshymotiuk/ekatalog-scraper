package scraper

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/yevhenshymotiuk/ekatalog-scraper/items"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Test Page</title>
</head>
<body>
<a href="/k298.htm" class="path_lnk">Ноутбуки</a>
<div id="top-page-title">
<b class="ib">Apple MacBook Pro 13 (2020)</b>
</div>
<table>
<tbody>
<tr class="conf-tr">
<td class="conf-td c21"><span title="Серия процессора">Core i5&nbsp;</span></td>
<td class="conf-td c21"><span title="Модель процессора">8257U&nbsp;</span></td>
<td class="conf-td c21"><span title="Объем оперативной памяти">8&nbsp;ГБ</span></td>
<td class="conf-td c21"><span title="Модель видеокарты">Iris Plus Graphics 645&nbsp;</span></td>
<td class="conf-td c21"><span title="Тип накопителя">SSD&nbsp;</span></td>
<td class="conf-td c21"><span title="Емкость накопителя">256&nbsp;ГБ</span></td>
<td class="conf-td conf-price" align="right"><span class="price-int"><span>36&nbsp;949&nbsp;</span>.. <span>43&nbsp;176&nbsp;</span>грн.</span></a></td>
</tr>
</tbody>
</table>
</body>
</html>`))
	})

	return httptest.NewServer(mux)
}

func TestScrapeLaptops(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	products, err := ScrapeProducts([]string{ts.URL + "/html"})
	if err != nil {
		t.Error("Failed to scrape products")
	}

	wantProducts := []items.Product{
		{
			Name: "Apple MacBook Pro 13 (2020)",
			Modifications: []items.ModificationType{
				items.Laptop{
					CPU: items.CPU{
						Series: "Core i5",
						Model:  "8257U",
					},
					RAM: items.RAM{
						Capacity: 8,
					},
					GPU: items.GPU{
						Model: "Iris Plus Graphics 645",
					},
					Drive: items.Drive{
						Type:     "SSD",
						Capacity: 256,
					},
					Price: items.Price{
						Min: 36949,
						Max: 43176,
					},
				},
			},
		},
	}

	for i, got := range products {
		want := wantProducts[i]

		if !reflect.DeepEqual(got, want) {
			t.Errorf("products are not equal: got: %+v, want: %+v", got, want)
		}
	}
}
