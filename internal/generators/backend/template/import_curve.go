package template

// ImportCurve ...
const ImportCurve = `

{{ define "import_fr" }}

{{ if eq .Curve "BLS377"}}
	"github.com/consensys/gurvy/bls377/fr"
{{ else if eq .Curve "BLS381"}}
	"github.com/consensys/gurvy/bls381/fr"
{{ else if eq .Curve "BN256"}}
	"github.com/consensys/gurvy/bn256/fr"
{{ else if eq .Curve "BW761"}}
	"github.com/consensys/gurvy/bw761/fr"
{{end}}

{{end}}

{{ define "import_curve" }}
{{if eq .Curve "BLS377"}}
	curve "github.com/consensys/gurvy/bls377"
	"github.com/consensys/gurvy/bls377/fr"
{{else if eq .Curve "BLS381"}}
	curve "github.com/consensys/gurvy/bls381"
	"github.com/consensys/gurvy/bls381/fr"
{{else if eq .Curve "BN256"}}
	curve "github.com/consensys/gurvy/bn256"	
	"github.com/consensys/gurvy/bn256/fr"
{{else if eq .Curve "BW761"}}
	curve "github.com/consensys/gurvy/bw761"	
	"github.com/consensys/gurvy/bw761/fr"
{{end}}

{{end}}

{{ define "import_backend" }}
{{if eq .Curve "BLS377"}}
	{{toLower .Curve}}backend "github.com/philsippl/gnark/internal/backend/bls377"
{{else if eq .Curve "BLS381"}}
	{{toLower .Curve}}backend "github.com/philsippl/gnark/internal/backend/bls381"
{{else if eq .Curve "BN256"}}
	{{toLower .Curve}}backend "github.com/philsippl/gnark/internal/backend/bn256"
{{ else if eq .Curve "BW761"}}
	{{toLower .Curve}}backend "github.com/philsippl/gnark/internal/backend/bw761"
{{end}}

{{end}}

{{ define "import_fft" }}
{{if eq .Curve "BLS377"}}
	"github.com/philsippl/gnark/internal/backend/bls377/fft"
{{else if eq .Curve "BLS381"}}
	"github.com/philsippl/gnark/internal/backend/bls381/fft"
{{else if eq .Curve "BN256"}}
	"github.com/philsippl/gnark/internal/backend/bn256/fft"
{{ else if eq .Curve "BW761"}}
	"github.com/philsippl/gnark/internal/backend/bw761/fft"
{{end}}
{{end}}

`
