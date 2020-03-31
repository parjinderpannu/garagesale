# garagesale

## How to run go run ./cmd/sales-api/
export SALES_DB_DISABLE_TLS=true  

## Helpful links
### database
//awesome-go.com/#/database  
https://github.com/jmoiron/sqlx  
https://golang.org/pkg/database/sql/  
https://github.com/golang/go/wiki/SQLDrivers  
https://godoc.org/github.com/lib/pq  
import github.com/lib/pq # driver gets registed when you import the package. You have to import that package to run its func init() to register the driver. But if you are not referring the package then it won't compile (use underscore _).  
note go.sum #verify version of mod we have Are same when we originally created it.  
https://awesome-go.com/ --> "Database schema migration." --> "darwin"

### Config links
https://github.com/kelseyhightower/envconfig  
https://github.com/peterbourgon/ff

### routes
https://awesome-go.com/#routers  --> chi  

### pprof
https://golang.org/pkg/net/http/pprof/  
cmd # hey -c 10 -n 15000 http://localhost:8000/v1/products  
cmd # go tool pprof http://localhost:6060/debug/pprof/profile\?seconds\=8  
cmd # top  
cmd # top -cum  
cmd # web  

 


