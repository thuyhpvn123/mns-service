package errors
import(
	"errors"
	"net/http"

)
var(
	ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))
)