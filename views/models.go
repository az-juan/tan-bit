package views

import (
	"net/http"
	"net/url"
	"strings"

	db "tan-bit.com/tan-bit/db/sqlc"
)

type FiltroCategoria struct {
	Key   db.Categoria
	Label string
	Count int
}

func CategoriaLabel(c db.Categoria) string {
	s := string(c)
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

type VenderData struct {
	Titulo      string
	Descripcion string
	Condicion   db.Condicion
	Categoria   db.Categoria
	Ubicacion   string
	Precio      string
	Stock       string
	TipoEntrega string
	CostoEnvio  string
	ZonaEntrega string
	RutaImagen  string
}

func ParseVenderData(r *http.Request) VenderData {
	return VenderData{
		Titulo:      r.FormValue("titulo"),
		Descripcion: r.FormValue("descripcion"),
		Condicion:   db.Condicion(r.FormValue("condicion")),
		Categoria:   db.Categoria(r.FormValue("categoria")),
		Ubicacion:   r.FormValue("ubicacion"),
		Precio:      r.FormValue("precio"),
		Stock:       r.FormValue("stock"),
		TipoEntrega: r.FormValue("tipo_entrega"),
		CostoEnvio:  r.FormValue("costo_envio"),
		ZonaEntrega: r.FormValue("zona_entrega"),
		RutaImagen:  r.FormValue("ruta_imagen"),
	}
}

func (d VenderData) Query() string {
	v := url.Values{}
	v.Set("titulo", d.Titulo)
	v.Set("descripcion", d.Descripcion)
	v.Set("condicion", string(d.Condicion))
	v.Set("categoria", string(d.Categoria))
	v.Set("ubicacion", d.Ubicacion)
	v.Set("precio", d.Precio)
	v.Set("tipo_entrega", d.TipoEntrega)
	v.Set("costo_envio", d.CostoEnvio)
	v.Set("zona_entrega", d.ZonaEntrega)
	v.Set("ruta_imagen", d.RutaImagen)
	return v.Encode()
}
func CategoriasDB() []db.Categoria { // revisar
	return []db.Categoria{
		db.CategoriaNotebook, db.CategoriaPantalla, db.CategoriaElectronica, db.CategoriaPeriferico,
		db.CategoriaAccesorio, db.CategoriaAudio, db.CategoriaVideo, db.CategoriaGaming, db.CategoriaAnalogico, db.CategoriaComunicacion,
	}
}
func CondicionesDB() []db.Condicion { //revisar
	return []db.Condicion{"nuevo", "muy bueno", "usado", "reacondicionado", "a reparar"}
}

func BuildFiltros(arts []db.GetArticulosRow) []FiltroCategoria { // revisar
	counts := map[db.Categoria]int{}
	for _, a := range arts {
		counts[a.Categoria]++
	}
	// order by db enum order (matches schema)
	order := []db.Categoria{"notebook", "pantalla", "electronica", "periferico", "accesorio", "audio", "video", "gaming", "analogico", "comunicacion"}
	var out []FiltroCategoria
	total := len(arts)
	out = append(out, FiltroCategoria{Key: "", Label: "TODOS", Count: total})
	for _, c := range order {
		if n, ok := counts[c]; ok && n > 0 {
			out = append(out, FiltroCategoria{Key: c, Label: CategoriaLabel(c), Count: n})
		}
	}
	return out

}

func BuildFiltrosFav(arts []db.GetFavoritosByUsuarioRow) []FiltroCategoria { // revisar
	counts := map[db.Categoria]int{}
	for _, a := range arts {
		counts[a.Categoria]++
	}
	// order by db enum order (matches schema)
	order := []db.Categoria{"notebook", "pantalla", "electronica", "periferico", "accesorio", "audio", "video", "gaming", "analogico", "comunicacion"}
	var out []FiltroCategoria
	total := len(arts)
	out = append(out, FiltroCategoria{Key: "", Label: "TODOS", Count: total})
	for _, c := range order {
		if n, ok := counts[c]; ok && n > 0 {
			out = append(out, FiltroCategoria{Key: c, Label: CategoriaLabel(c), Count: n})
		}
	}
	return out
}
