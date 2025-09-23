INSERT INTO "statuses" ( "statusId", "title", "alias" ) VALUES ( 1, 'Опубликован', 'enabled' );
INSERT INTO "statuses" ( "statusId", "title", "alias" ) VALUES ( 2, 'Не опубликован', 'disabled' );
INSERT INTO "statuses" ( "statusId", "title", "alias" ) VALUES ( 3, 'Удален', 'deleted' );

-- password is 12345
INSERT INTO "users" ( "login", "password", "statusId" ) VALUES ( 'admin', '$2y$14$4IpqlaJ2Rvfgs.wb8f6lPODVLb/Ygl6zw1ZCUKz5CuT6WB6CV44AG', 1 );

INSERT INTO "vfsFolders" ("parentFolderId", title, "isFavorite", "createdAt", "statusId") VALUES (null, 'root', false, now(), 1);


INSERT INTO tags ("tagId", name, "statusId")
VALUES (1, 'Mascots', 1),
       (2, 'DRAFT', 2),
       (3, 'DELETED', 3);

INSERT INTO categories ("categoryId", title, "statusId")
VALUES (1, 'Accidents', 1),
       (2, 'DRAFT', 2),
       (3, 'DELETED', 3),
       (4, 'Events', 1);

INSERT INTO news (title, "shortText", content, author, "categoryId", "tagIds", "publishedAt", "statusId")
VALUES (
           -- Published
           'Drunk cat occurred massive traffic jam in the LA',
           'Breaking news from Los Angeles: a stray cat, apparently intoxicated from spilled alcohol, caused a massive traffic jam yesterday at the busy intersection of 5th and Main.',
           'In an unprecedented incident yesterday, a stray tabby cat believed to be intoxicated by spilled alcohol caused a massive traffic jam on downtown Los Angeles streets. Witnesses reported seeing the feline zigzagging across lanes near the intersection of 5th and Main, prompting drivers to slow down and stop altogether. Authorities suspect the cat may have ingested discarded alcohol from nearby trash cans. Animal control was called to safely retrieve the feline, and traffic was gradually restored after the animal was secured. Experts warn that stray animals consuming alcohol can exhibit unpredictable behavior, posing risks to both themselves and motorists.',
           'Bob the Cat',
           1,
           ARRAY [1],
           '2025-09-15 00:00:00 UTC',
           1),
       (
           -- Drafted
           'Draft: Drunk cat plays very sad blues in the downtown bar',
           'Last night, a mysterious feline, dubbed "Johnny Purr," caused a stir at the downtown bar by climbing onto the stage and unleashing a soulful, yet profoundly sad blues performance.',
           'In a surprising turn of events, a stray tabby named Whiskers was spotted last night at the local downtown bar, seemingly intoxicated and passionately playing a worn-out harmonica. Eyewitnesses claim the feline appeared melancholy, strumming soulful blues that moved the entire crowd to tears. Authorities are investigating whether Whiskers was given alcohol or if it''s a bizarre new trend among street cats trying to break into the music scene. Fans are already calling for a live album, dubbing the cat "The Blues Purrformer."',
           NULL,
           4,
           ARRAY [1],
           '2025-09-16 00:00:00 UTC',
           2),
       (
           -- Scheduled
           -- Tests will fail after September 2030. Not enough faith in the project, don't u think?
           'Drunk cats build starship',
           'Drunk Cats Build Starship in Backyard Laboratory',
           'In an astonishing turn of events, a group of neighborhood cats, reportedly intoxicated from spilled milk and leftover fish, have reportedly constructed a makeshift starship in a backyard laboratory. Witnesses claim the feline engineers, dubbed the "Meow-ronauts," spent weeks assembling the vessel using household items and scrap metal. While experts remain skeptical, some believe this bizarre incident hints at a new frontier in animal intelligence, or perhaps just a very creative feline party gone awry. The local authorities are investigating, but for now, the starship remains a mysterious and whimsical fixture in the suburban yard.',
           'Bob the Cat',
           4,
           ARRAY [1],
           '2030-09-16 00:00:00 UTC',
           1),
       (
           -- Deleted
           'Home cats beauty competition event just ended in Bronx',
           'But who cares?',
           NULL,
           NULL,
           4,
           ARRAY [1],
           '2025-09-15 00:00:00 UTC',
           3);