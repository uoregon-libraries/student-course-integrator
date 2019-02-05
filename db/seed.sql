DELETE FROM courses;
ALTER TABLE courses AUTO_INCREMENT = 0;
INSERT INTO courses (id, canvas_id, description) VALUES
  (1, '201704.01071','DSGN 404 (Summer 2018) Internship'),
  (2, '201704.01074','DSGN 604 (Summer 2018) Internship'),
  (3, '201704.01075','DSGN 604 (Summer 2018) Internship'),
  (4, '201704.02076','AAAP 503 (Summer 2018) Thesis'),
  (5, '201704.02077','AAAP 601 (Summer 2018) Research'),
  (6, '201704.02078','AAAP 604 (Summer 2018) Internship'),
  (7, '201704.02079','AAAP 605 (Summer 2018) Reading'),
  (8, '201704.02080','AAAP 606 (Summer 2018) Special Problems'),
  (9, '201704.02081','AAAP 611 (Summer 2018) Terminal Project');

DELETE FROM faculty_courses;
ALTER TABLE faculty_courses AUTO_INCREMENT = 0;
INSERT INTO faculty_courses (login, course_id) VALUES
  ('dsgnprof', 1),
  ('dsgnprof', 2),
  ('dsgnprof', 3),
  ('aaapprof', 4),
  ('aaapprof', 5),
  ('aaapprof', 6),
  ('aaapprof', 7),
  ('aaapprof', 8),
  ('aaapprof', 9),
  ('noidear', 3),
  ('noidear', 9);
