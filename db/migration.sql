CREATE TABLE competition (
    cmp_id INT IDENTITY (1,1),
    cmp_name VARCHAR,
    cmp_last_registration_date DATE,
    cmp_start_date DATE,
    cmp_prize_pool INT,
    cmp_desc TEXT,
    cmp_image_src TEXT,
    PRIMARY KEY (cmp_id)
)