

CREATE TABLE danbing (
        id SERIAL UNIQUE NOT NULL,
        code VARCHAR(10) NOT NULL, -- not unique
        article TEXT,
        uptime  TIMESTAMP,
        name TEXT NOT NULL -- not unique
);


insert into danbing (
    code, article, uptime, name
)
select
    left(md5(i::text), 10),
    md5(random()::text),
    now()::timestamp,
    md5(random()::text)
from generate_series(1, 10000000) s(i)