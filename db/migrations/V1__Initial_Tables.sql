CREATE TABLE users
(
  id            INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
  username      VARCHAR(255) NOT NULL UNIQUE,
  email         VARCHAR(255) NOT NULL UNIQUE,
  password_hash TEXT         NOT NULL
) ENGINE = INNODB;

CREATE TABLE cookbooks
(
  id      INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
  user_id INT         NOT NULL,
  name    VARCHAR(75) NOT NULL,

  FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
) ENGINE = INNODB;

CREATE TABLE sections
(
  id          INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
  cookbook_id INT         NOT NULL,
  name        VARCHAR(75) NOT NULL,

  FOREIGN KEY (cookbook_id)
    REFERENCES cookbooks (id)
    ON DELETE CASCADE
) ENGINE = INNODB;

CREATE TABLE recipes
(
  id          INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
  creator     INT          NOT NULL,
  name        VARCHAR(400) NOT NULL,
  description TEXT         NOT NULL,
  servings    INT          NOT NULL,
  prep_time   VARCHAR(50),
  cook_time   VARCHAR(50),
  cool_time   VARCHAR(50),
  total_time  VARCHAR(50),
  source      VARCHAR(400),

  FOREIGN KEY (creator)
    REFERENCES users (id)
    ON DELETE CASCADE
) ENGINE = INNODB;

CREATE TABLE ingredients
(
  id   INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL
) ENGINE = INNODB;

CREATE TABLE measurements
(
  id   INT         NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL
) ENGINE = INNODB;

CREATE TABLE recipe_ingredients
(
  recipe_id      INT NOT NULL,
  ingredient_id  INT NOT NULL,
  ingredient_no  INT NOT NULL,
  amount         VARCHAR(50),
  measurement_id INT,
  preparation    varchar(255),

  PRIMARY KEY (recipe_id, ingredient_id, ingredient_no),

  FOREIGN KEY (recipe_id)
    REFERENCES recipes (id)
    ON DELETE CASCADE,
  FOREIGN KEY (ingredient_id)
    REFERENCES ingredients (id)
    ON DELETE CASCADE,
  FOREIGN KEY (measurement_id)
    REFERENCES measurements (id)
    ON DELETE CASCADE
) ENGINE = INNODB;

CREATE TABLE recipe_steps
(
  recipe_id    INT  NOT NULL,
  step_no      INT  NOT NULL,
  instructions text NOT NULL,

  PRIMARY KEY (recipe_id, step_no),

  FOREIGN KEY (recipe_id)
    REFERENCES recipes (id)
    ON DELETE CASCADE
) ENGINE = INNODB;

CREATE TABLE recipe_locations
(
  id          INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  recipe_id   INT NOT NULL,
  cookbook_id INT,
  section_id  INT,

  FOREIGN KEY (recipe_id)
    REFERENCES recipes (id)
    ON DELETE CASCADE,
  FOREIGN KEY (cookbook_id)
    REFERENCES recipes (id)
    ON DELETE CASCADE,
  FOREIGN KEY (section_id)
    REFERENCES recipes (id)
    ON DELETE CASCADE
) ENGINE = INNODB;

ALTER TABLE recipe_locations
  ADD CONSTRAINT CK_nulltest
    CHECK (cookbook_id IS NOT NULL OR section_id IS NOT NULL);
