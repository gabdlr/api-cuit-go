package cuit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gabdlr/api-cuit-go/utils"
)

type Address struct {
	Provincia         string `json:"provincia"`
	Localidad         string `json:"localidad"`
	Domicilio         string `json:"domicilio"`
	PisoDeptoOfi      string `json:"pisoDeptoOfi"`
	CodigoPostal      string `json:"codigoPostal"`
	EstadoDeDomicilio string `json:"estadoDeDomicilio"`
}

type Society struct {
	RazonSocial         string `json:"razonSocial"`
	Cuit                string `json:"cuit"`
	TipoSocietario      string `json:"tipoSocietario"`
	FechaDeContrato     string `json:"fechaDeContrato"`
	NumeroRegistroLocal string `json:"numeroRegistroLocal"`
}

type CuitInfo struct {
	Sociedad           Society `json:"sociedad"`
	DomicilioFiscal    Address `json:"domicilioFiscal"`
	DomicilioLegal     Address `json:"domicilioLegal"`
	FechaActualizacion string  `json:"fechaActualizacion"`
}

const htmlOfInterestStart = `<tbody`
const htmlOfInterestEnd = `</tbody`
const exitSignal = "No se encuentran resultados"

func Search(cuit string) ([]byte, error) {
	url := fmt.Sprintf("https://argentina.gob.ar/justicia/registro-nacional-sociedades?cuit=%s&razon=", utils.StandardizeCuit(cuit))

	res, err := http.Get(url)
	if err != nil {
		return []byte{0}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{0}, err
	}

	cuitInfo, err := parseResponse(string(body))
	if err != nil {
		return []byte{0}, err
	}

	cuitInfoJSON, err := json.Marshal(cuitInfo)
	if err != nil {
		return []byte{0}, err
	}
	return cuitInfoJSON, nil
}

func searchElements(s, startMarker, endMarker string) []string {
	elements := make([]string, 0)
	thereAreMoreElements := true
	for thereAreMoreElements {
		startElement := strings.Index(s, startMarker)
		endElement := strings.Index(s, endMarker)
		if startElement > -1 || endElement > -1 {
			elements = append(elements, s[startElement+len(startMarker):endElement])
			s = s[endElement+len(endMarker):]
		} else {
			thereAreMoreElements = false
		}
	}
	return elements
}

func searchParagraphElements(s string) []string {
	return searchElements(s, "<p>", "</p>")
}

func updateCuitInfo(cuitInfo *CuitInfo, keyElements []string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		sociedadElements := searchParagraphElements(keyElements[0])
		if len(sociedadElements) == 5 {
			cuitInfo.Sociedad.RazonSocial = sociedadElements[0]
			cuitInfo.Sociedad.Cuit = sociedadElements[1]
			cuitInfo.Sociedad.TipoSocietario = sociedadElements[2]
			cuitInfo.Sociedad.FechaDeContrato = sociedadElements[3]
			cuitInfo.Sociedad.NumeroRegistroLocal = sociedadElements[4]
		}
	}()
	go func() {
		defer wg.Done()
		domicilioFiscalElements := searchParagraphElements(keyElements[1])
		if len(domicilioFiscalElements) == 6 {
			cuitInfo.DomicilioFiscal.Provincia = domicilioFiscalElements[0]
			cuitInfo.DomicilioFiscal.Localidad = domicilioFiscalElements[1]
			cuitInfo.DomicilioFiscal.Domicilio = domicilioFiscalElements[2]
			cuitInfo.DomicilioFiscal.PisoDeptoOfi = domicilioFiscalElements[3]
			cuitInfo.DomicilioFiscal.CodigoPostal = domicilioFiscalElements[4]
			cuitInfo.DomicilioFiscal.EstadoDeDomicilio = domicilioFiscalElements[5]
		}
	}()
	go func() {
		defer wg.Done()
		domicilioLegalElements := searchParagraphElements(keyElements[2])
		if len(domicilioLegalElements) == 6 {
			cuitInfo.DomicilioLegal.Provincia = domicilioLegalElements[0]
			cuitInfo.DomicilioLegal.Localidad = domicilioLegalElements[1]
			cuitInfo.DomicilioLegal.Domicilio = domicilioLegalElements[2]
			cuitInfo.DomicilioLegal.PisoDeptoOfi = domicilioLegalElements[3]
			cuitInfo.DomicilioLegal.CodigoPostal = domicilioLegalElements[4]
			cuitInfo.DomicilioLegal.EstadoDeDomicilio = domicilioLegalElements[5]
		}
	}()
	go func() {
		defer wg.Done()
		fechaActualizacionElements := searchParagraphElements(keyElements[3])
		if len(fechaActualizacionElements) == 1 {
			cuitInfo.FechaActualizacion = fechaActualizacionElements[0]
		}
	}()
}

func parseResponse(html string) (cuitInfo CuitInfo, err error) {
	notFoundErr := "información no disponible"

	if strings.Contains(html, exitSignal) {
		err = errors.New(notFoundErr)
		return cuitInfo, err
	}

	startPosition := strings.Index(html, htmlOfInterestStart)
	endPosition := strings.Index(html, htmlOfInterestEnd)

	if startPosition != -1 && endPosition != -1 {
		info := html[startPosition:endPosition]
		startMarker := "<td"
		endMarker := "</td"
		keyElements := searchElements(info, startMarker, endMarker)

		if len(keyElements) == 4 {
			var wg sync.WaitGroup
			wg.Add(4)
			updateCuitInfo(&cuitInfo, keyElements, &wg)
			wg.Wait()
		}
	} else {
		err = errors.New(notFoundErr)
	}
	return cuitInfo, err
}
