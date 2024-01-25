# enricher-service
Сервис, который получает по api ФИО в виде JSON объектов и из открытых api дополняет ответ наиболее вероятными возрастом, полом и национальностью и сохраняет данные в БД.



to start the  program do:
-- docker-compose up -d
-- goose -dir ./migrations postgres "postgres://admin:qwerty@localhost:5432/human?sslmode=disable" up

to down the program
--docker-compose down
--goose -dir ./migrations postgres "postgres://admin:qwerty@localhost:5432/human?sslmode=disable" up  down