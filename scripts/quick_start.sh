ps aux | grep "sdb" | grep 'config' | awk '{print "kill -9 " $2}' | sh -x
nohup go run cmd/sdb/main.go -config ./configs/master.yml > logs/master.log 2>&1 &
sleep 5;
nohup go run cmd/sdb/main.go -config ./configs/slave1.yml > logs/slave1.log 2>&1 &
sleep 5;
nohup go run cmd/sdb/main.go -config ./configs/slave2.yml > logs/slave2.log 2>&1 &