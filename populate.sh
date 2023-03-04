#echo $RANDOM
## 准备数据
for shard in localhost:8080 localhost:8081; do
   echo $shard
    for i in {1..1000}; do
   echo curl "http//:$shard/set?key=key-$RANDOM&value=value-$RANDOM"
  done
done