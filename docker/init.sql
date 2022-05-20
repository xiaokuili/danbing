

CREATE TABLE danbing (
        id SERIAL UNIQUE NOT NULL,
        code VARCHAR(10) NOT NULL, -- not unique
        article TEXT,
        name TEXT NOT NULL -- not unique
);


insert into danbing (
    code, article, name
)
select
    left(md5(i::text), 10),
    md5(random()::text),
    md5(random()::text)
from generate_series(1, 100) s(i)