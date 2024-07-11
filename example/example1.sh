#initial build
go build -ldflags '-X main.BuildID=1' -o test_app
echo "BUILT APP (1)"

echo "RUNNING APP"
./test_app &
APP_ID=$!

sleep 3
curl localhost:5001
curl localhost:5001

#request during an update
curl localhost:5001?d=5s &
curl localhost:5001?d=10s &
curl localhost:5001?d=15s &

go build -ldflags '-X main.BuildID=2' -o test_app_next
echo "BUILT APP (2)"

sleep 3
curl localhost:5001
curl localhost:5001
curl localhost:5001
curl localhost:5001

sleep 30
kill $APP_ID
rm test_app* 2> /dev/null
