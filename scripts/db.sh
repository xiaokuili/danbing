#!/bin/bash

createtable(){
    docker exec -it docker_pg_1 /bin/bash -c 'psql -U postgres -h 127.0.0.1 -c "CREATE TABLE danbing (
        user_id serial PRIMARY KEY,
        name VARCHAR ( 50 ) NOT NULL
    );"'
}


showtable(){
    docker exec -it docker_pg_1 /bin/bash -c  'psql -U postgres -h 127.0.0.1 -c "\dt"'
}


insert() {
       docker exec -it docker_pg_1 /bin/bash -c 'psql -U postgres -h 127.0.0.1 -c  "INSERT INTO danbing VALUES (1, 1);
        INSERT INTO danbing VALUES (2, 2);
        INSERT INTO danbing VALUES (3, 3);
        INSERT INTO danbing VALUES (4, 4);
        INSERT INTO danbing VALUES (5, 5);
        INSERT INTO danbing VALUES (6, 6);
        INSERT INTO danbing VALUES (7, 7);
        INSERT INTO danbing VALUES (8, 8);
        INSERT INTO danbing VALUES (9, 9);
        INSERT INTO danbing VALUES (10, 10);
        INSERT INTO danbing VALUES (11, 11);
        INSERT INTO danbing VALUES (12, 12);
        INSERT INTO danbing VALUES (13, 13);
        INSERT INTO danbing VALUES (14, 14);
        INSERT INTO danbing VALUES (15, 15);
        INSERT INTO danbing VALUES (16, 16);
        INSERT INTO danbing VALUES (17, 17);
        INSERT INTO danbing VALUES (18, 18);
        INSERT INTO danbing VALUES (19, 19);
        INSERT INTO danbing VALUES (20, 20);"' 
}


search() {
        docker exec -it docker_pg_1 /bin/bash -c  'psql -U postgres -h 127.0.0.1 -c "select * from danbing offset 0 limit 10;"'

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
