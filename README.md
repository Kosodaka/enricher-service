# enricher-service
Сервис, который получает по api ФИО в виде JSON объектов и из открытых api дополняет ответ наиболее вероятными возрастом, полом и национальностью и сохраняет данные в БД.

to make build or run use git bash or linux console

-to run locale:
```
 make run-windows
 or 
 make run-linux
```

- to run service in docker-container:
```
 make doker-up
```
-to make migrations
```
make migrate
```
- to down migrations
```
 make migrate-down
```
-to run tests 
```
 make test 
```