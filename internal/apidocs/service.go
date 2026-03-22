package apidocs

import "net/http"

type Service struct {
	openAPI []byte
	scalar  []byte
}

func NewService() *Service {
	return &Service{
		openAPI: MustMarshalIndentedDocument(),
		scalar:  []byte(scalarHTML),
	}
}

func (s *Service) OpenAPIJSON() []byte {
	return s.openAPI
}

func (s *Service) ScalarHTML() []byte {
	return s.scalar
}

func (s *Service) ServeOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(s.openAPI)
}

func (s *Service) ServeScalar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(s.scalar)
}

const scalarHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Student Service App API Docs</title>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/openapi.json"
      data-configuration='{"theme":"purple","layout":"modern","darkMode":false}'
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`
