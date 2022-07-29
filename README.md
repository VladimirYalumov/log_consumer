Consumer and Producer for logs in rest projects
=============================================
This project contains a consumer what read queues in rabbit and put them to the mongo db.

In the package logger you can find the pusher what put message to the rabbit

<h3>How to start?</h3>
<h5>spoiler: easy)</h5>

<pre>
    // download log_consumer
    git clone https://github.com/VladimirYalumov/log_consumer.git

    // download and run environment
    git clone https://github.com/VladimirYalumov/docker_pgdb_mongodb_redis.git
    cd docker_pgdb_mongodb_redis
    docker-compose up -d

    // go to project
    cd log_consumer
    go mod vendor
    go run main.go
</pre>

<h3>How to test?</h3>
<h5>spoiler: even easier than start</h5>

<pre>
    // go to project test client
    cd log_consumer/test_client
    go run main.go
    curl --location --request GET '127.0.0.1:8089/numbers_sum' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "first_num": 1,
        "second_num": 2
    }'
</pre>
And watch the result in the same mongo client?

Configs
in the block config (if you use my environment) you can't change anything
In the block mongo you must add your own logical mongo databases
And add all routes you have in project
Example:
<pre>
mongo:
    auth:
        signin: signin
        signup: signup
        signout: signout
    event:
        add: add_event
        update: update_event
        delete: delete_event
</pre>

P. S. Consumer create the collection with current time of incoming log.
Example:
<pre>
    update_event_202207
</pre>