rm -rf ./*db
rm -rf ./*db.lock
rm -rf ./log/*.csv
go run main.go -S 4 -f 1 -c -t len140.csv &
go run main.go -S 4 -f 1 -s S0 -n N1 &
go run main.go -S 4 -f 1 -s S0 -n N2 &
go run main.go -S 4 -f 1 -s S0 -n N3 &
go run main.go -S 4 -f 1 -s S1 -n N1 &
go run main.go -S 4 -f 1 -s S1 -n N2 &
go run main.go -S 4 -f 1 -s S1 -n N3 &
go run main.go -S 4 -f 1 -s S2 -n N1 &
go run main.go -S 4 -f 1 -s S2 -n N2 &
go run main.go -S 4 -f 1 -s S2 -n N3 &
go run main.go -S 4 -f 1 -s S3 -n N1 &
go run main.go -S 4 -f 1 -s S3 -n N2 &
go run main.go -S 4 -f 1 -s S3 -n N3 &
go run main.go -S 4 -f 1 -s SC -n N1 &
go run main.go -S 4 -f 1 -s SC -n N2 &
go run main.go -S 4 -f 1 -s SC -n N3 &
go run main.go -S 4 -f 1 -s S0 -n N0 -t len140.csv &
go run main.go -S 4 -f 1 -s S1 -n N0 -t len140.csv &
go run main.go -S 4 -f 1 -s S2 -n N0 -t len140.csv &
go run main.go -S 4 -f 1 -s S3 -n N0 -t len140.csv &
go run main.go -S 4 -f 1 -s SC -n N0 &