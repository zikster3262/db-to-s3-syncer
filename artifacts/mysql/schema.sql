CREATE TABLE request (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  uuid text,
  time text
);

CREATE TABLE files (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  request_id INT,
  time text,
  FOREIGN KEY (request_id) REFERENCES request(id)
);