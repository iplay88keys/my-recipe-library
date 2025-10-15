-- password_hash values are bcrypt hashes of:
-- user: "password"
-- user2: "password2"
INSERT INTO users (username, email, password_hash)
VALUES ("user", "email@example.com", "$2a$10$OrHtjl4.k1bvqakWLJPkb.SoPQfhs70NI1TvgkAncJb4DfUGdeTFa"),
       ("user2", "email2@example.com", "$2a$10$YQ2GjqrcCthWQTtI7Fke3uXhS9Nl0Bf.EA7.gwcLACChiF9Tuoj1m");

INSERT INTO cookbooks (user_id, name)
VALUES (1, "Favorites"),
       (2, "Drinks");

INSERT INTO sections (name, cookbook_id)
VALUES ("Sides", 1);

INSERT INTO recipes
(creator,
 name,
 description,
 servings,
 prep_time,
 cook_time,
 cool_time,
 total_time,
 source)
VALUES (1,
        "Root Beer Float",
        "Delicious drink for a hot summer day.",
        1, "5 m", NULL, NULL, "5 m", NULL),
       (2,
        "Nana's Beans",
        "Spruced up baked beans.",
        8, "10 m", "1-2 hrs", NULL, "1-2 hrs",
        NULL);

INSERT INTO ingredients (name)
VALUES ("Vanilla Ice Cream"),
       ("Root Beer"),
       ("Baked Beans"),
       ("Bacon"),
       ("Onion"),
       ("Ketchup"),
       ("Worcestershire Sauce"),
       ("Brown Sugar"),
       ("Yellow Mustard"),
       ("BBQ Sauce");

INSERT INTO measurements (name)
VALUES ("Scoop"),
       ("Cup"),
       ("Can"),
       ("Tablespoon");

INSERT INTO recipe_ingredients (recipe_id, ingredient_id, ingredient_no, amount, measurement_id, preparation)
VALUES (1, 1, 1, 1, 1, NULL),
       (1, 2, 2, NULL, NULL, NULL),
       (2, 3, 1, "2 28oz", NULL, NULL),
       (2, 4, 2, "3-4", NULL, "sliced"),
       (2, 5, 3, "2", 4, NULL),
       (2, 6, 4, "1", 4, NULL),
       (2, 7, 5, "2", 4, NULL),
       (2, 8, 6, "1", 4, NULL),
       (2, 9, 7, "2-3", 4, NULL);

INSERT INTO recipe_steps (recipe_id, step_no, instructions)
VALUES (1, 1, "Place ice cream in glass."),
       (1, 2, "Top with Root Beer."),
       (2, 1, "Combine all ingredients in a pot."),
       (2, 2, "Simmer 1-2 hours.");

INSERT INTO recipe_locations (recipe_id, cookbook_id, section_id)
VALUES (1, 2, NULL),
       (2, NULL, 1);
