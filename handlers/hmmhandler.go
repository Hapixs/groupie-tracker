package handlers

import (
	"encoding/base64"
	"net/http"
)

func hmmHandler(w http.ResponseWriter, r *http.Request) {
	decoded, _ := base64.StdEncoding.DecodeString("PCFET0NUWVBFIGh0bWw+IDxodG1sPiA8aGVhZD4gPC9oZWFkPiA8Ym9keSBzdHlsZT1cImJhY2tncm91bmQtaW1hZ2U6IHVybChodHRwczovL2kueXRpbWcuY29tL3ZpL1FsRjdjeWRuVjJFL21heHJlc2RlZmF1bHQuanBnKTsgYmFja2dyb3VuZC1yZXBlYXQ6IG5vLXJlcGVhdDtiYWNrZ3JvdW5kLWF0dGFjaG1lbnQ6IGZpeGVkO2JhY2tncm91bmQtc2l6ZTogY292ZXI7XCI+IDwvYm9keT4gPC9odG1sPg==")
	w.Write([]byte(decoded))
}
