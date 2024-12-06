del *db
del *db.lock
cd ./log
del *.csv
cd..
start "011" go run main.go -S 4 -f 1 -s S0 -n N1
start "012" go run main.go -S 4 -f 1 -s S0 -n N2
start "013" go run main.go -S 4 -f 1 -s S0 -n N3
start "015" go run main.go -S 4 -f 1 -s S1 -n N1
start "016" go run main.go -S 4 -f 1 -s S1 -n N2
start "017" go run main.go -S 4 -f 1 -s S1 -n N3
start "019" go run main.go -S 4 -f 1 -s S2 -n N1
start "020" go run main.go -S 4 -f 1 -s S2 -n N2
start "021" go run main.go -S 4 -f 1 -s S2 -n N3
start "023" go run main.go -S 4 -f 1 -s S3 -n N1
start "024" go run main.go -S 4 -f 1 -s S3 -n N2
start "025" go run main.go -S 4 -f 1 -s S3 -n N3

start "026" go run main.go -S 4 -f 1 -s SC -n N0
start "027" go run main.go -S 4 -f 1 -s SC -n N1
start "028" go run main.go -S 4 -f 1 -s SC -n N2
start "029" go run main.go -S 4 -f 1 -s SC -n N3

start "000"  go run main.go -S 4 -f 1 -s S0 -n N0 -t len3.csv
start "001"  go run main.go -S 4 -f 1 -s S1 -n N0 -t len3.csv
start "002"  go run main.go -S 4 -f 1 -s S2 -n N0 -t len3.csv
start "003"  go run main.go -S 4 -f 1 -s S3 -n N0 -t len3.csv

start "0client" go run main.go -S 4 -f 1 -c -t len3.csv