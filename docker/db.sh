#!/bin/bash


search() {
        docker exec -it docker_pg_1 /bin/bash -c  'psql -U postgres -h 127.0.0.1 -c "select * from danbing limit 1;"'

} 

searchES() {
    curl -X GET "localhost:9200/danbing/_search?pretty" -H 'Content-Type: application/json' -d'
    {
    "query": {
        "match_all": { }
    },
    "size": 30
    }
    '
}


search
