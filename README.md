# golang-excercise
simple golang excercise using rabbitmq, redis, and channels 


For its use just run:

go run main.go ./exercise.json

It assumes that a local rabbitmq and redis is configured

Successful execution will result:

Successfully sent 
successfully received 
Successfully set values
Successfully get values
Successfully validated data

| Data      | Throughput KB | Thoughput MB |
|-----------|---------------|--------------|
|Average ul | 5997188       | 5997         |
|Highest ul | 56789454      | 56789        |
|Average dl | 28399651872   | 28399651     |
|Highest dl | 567890234320  | 567890234    |
End with: Ok
