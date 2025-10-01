# GET /health
curl -i http://localhost:8080/health

# POST /tasks (OK)
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk"}'

# POST /tasks (422 - слишком короткий title)
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"ok"}'

# GET /tasks (список)
curl -i http://localhost:8080/tasks

# GET /tasks?q=молоко (фильтр)
curl -i "http://localhost:8080/tasks?q=молоко"

# GET /tasks/{id}
curl -i http://localhost:8080/tasks/1

# PATCH /tasks/{id} -> done=true
curl -i -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"done":true}'

# DELETE /tasks/{id}
curl -i -X DELETE http://localhost:8080/tasks/1

# Ошибки:
curl -i http://localhost:8080/tasks/abc
