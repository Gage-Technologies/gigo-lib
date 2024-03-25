-- MySQL v5.7

ALTER TABLE bytes ADD COLUMN files_easy JSON;
ALTER TABLE bytes ADD COLUMN files_medium JSON;
ALTER TABLE bytes ADD COLUMN files_hard JSON;
UPDATE bytes
SET
    files_easy = JSON_ARRAY(
                JSON_OBJECT(
                    'content', outline_content_easy,
                    'file_name',
                        CASE
                            WHEN lang = 5 THEN 'main.py'
                            WHEN lang = 6 THEN 'main.go'
                            ELSE 'default.txt'
                        END
                )
            ),
    files_medium = JSON_ARRAY(
            JSON_OBJECT(
                    'content', outline_content_medium,
                    'file_name',
                    CASE
                        WHEN lang = 5 THEN 'main.py'
                        WHEN lang = 6 THEN 'main.go'
                        ELSE 'default.txt'
                        END
            )
                 ),
    files_hard = JSON_ARRAY(
            JSON_OBJECT(
                    'content', outline_content_hard,
                    'file_name',
                    CASE
                        WHEN lang = 5 THEN 'main.py'
                        WHEN lang = 6 THEN 'main.go'
                        ELSE 'default.txt'
                        END
            )
                 );


ALTER TABLE byte_attempts ADD COLUMN files_easy JSON;
ALTER TABLE byte_attempts ADD COLUMN files_medium JSON;
ALTER TABLE byte_attempts ADD COLUMN files_hard JSON;
UPDATE byte_attempts ba
SET
    files_easy = (
        SELECT JSON_ARRAY(
                       JSON_OBJECT(
                               'content', ba.content_easy,
                               'file_name', CASE
                                                WHEN b.lang = 5 THEN 'main.py'
                                                WHEN b.lang = 6 THEN 'main.go'
                                                ELSE 'default.txt'
                                   END
                       )
               )
        FROM bytes b
        WHERE b._id = ba.byte_id
    ),
    files_medium = (
        SELECT JSON_ARRAY(
                       JSON_OBJECT(
                               'content', ba.content_medium,
                               'file_name', CASE
                                                WHEN b.lang = 5 THEN 'main.py'
                                                WHEN b.lang = 6 THEN 'main.go'
                                                ELSE 'default.txt'
                                   END
                       )
               )
        FROM bytes b
        WHERE b._id = ba.byte_id
    ),
    files_hard = (
        SELECT JSON_ARRAY(
                       JSON_OBJECT(
                               'content', ba.content_hard,
                               'file_name', CASE
                                                WHEN b.lang = 5 THEN 'main.py'
                                                WHEN b.lang = 6 THEN 'main.go'
                                                ELSE 'default.txt'
                                   END
                       )
               )
        FROM bytes b
        WHERE b._id = ba.byte_id
    )
WHERE EXISTS (
    SELECT 1
    FROM bytes b
    WHERE b._id = ba.byte_id
);

