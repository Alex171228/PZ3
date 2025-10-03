# Практическое занятие №3  
**Дисциплина:** Технологии индустриального программирования   
## Цели работы
-	Освоить базовую работу со стандартной библиотекой net/http без сторонних фреймворков.
-	Научиться поднимать HTTP-сервер, настраивать маршрутизацию через http.ServeMux.
-	Научиться обрабатывать параметры запроса (query, path), тело запроса (JSON/form-data) и формировать корректные ответы (код статуса, заголовки, JSON).
-	Научиться базовому логированию запросов и обработке ошибок.
## Структура проекта
<img width="413" height="517" alt="изображение" src="https://github.com/user-attachments/assets/749444b0-d23a-4037-b03c-61c2718809b9" />

## Запуск проекта
### стандартный запуск
go run .\cmd\server

### запуск на другом порту
$env:PORT="9090"; go run .\cmd\server

### или через build.ps1
.\build.ps1 -Port 9090
## Запросы и результат их выполнения
### 1. Запрос /health
   <img width="800" height="61" alt="изображение" src="https://github.com/user-attachments/assets/9b4b0362-c80c-4c6d-9cc3-ae22e825fa5d" /> 
### 2. Создание задачи
   <img width="1108" height="739" alt="изображение" src="https://github.com/user-attachments/assets/48f63da1-6369-4db3-ade4-b4339df3c2a3" /> 
### 3. Список задач
   <img width="777" height="55" alt="изображение" src="https://github.com/user-attachments/assets/e19048f5-9130-40a5-b0a0-992c05f93026" /> 
### 4. Фильтр
   <img width="910" height="58" alt="изображение" src="https://github.com/user-attachments/assets/c1cabeda-1d1c-4758-8a8d-15e245e18813" /> 
### 5. Получение по id
   <img width="812" height="53" alt="изображение" src="https://github.com/user-attachments/assets/1f492bab-2736-4ca4-8f6a-6ff056c2a644" /> 
### 6. Обновление (PATCH)
   <img width="1139" height="759" alt="изображение" src="https://github.com/user-attachments/assets/76892ef0-880b-4518-8cd7-e19f5918c00c" /> 
### 7. Удаление
   <img width="1124" height="529" alt="изображение" src="https://github.com/user-attachments/assets/b6d49613-08af-4fab-9f4c-515d221175cd" /> 
   <img width="789" height="61" alt="изображение" src="https://github.com/user-attachments/assets/10cc7fc8-ff07-4418-a5bc-a28d71875e11" /> 
   










