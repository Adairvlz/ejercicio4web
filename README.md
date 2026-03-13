# 🏀 NBA Teams API

API REST desarrollada en **Go** usando únicamente la librería estándar `net/http`. Permite gestionar equipos de la NBA con persistencia real en un archivo JSON.

---

## 📁 Estructura del repositorio

```
ejercicio4web/
├── main.go
├── Dockerfile
├── docker-compose.yml
├── README.md
└── data/
    └── teams.json
```

---

## 🚀 Cómo correr el servidor

### Local
```bash
go run main.go
```

### Docker
```bash
docker-compose up --build
```

El servidor corre en el puerto **24596**.

---

## 📋 Modelo de datos

Cada equipo tiene la siguiente estructura:

```json
{
  "id":            1,
  "name":          "Boston Celtics",
  "city":          "Boston",
  "championships": 17,
  "pet":           "Lucky the Leprechaun",
  "arena":         "TD Garden"
}
```

| Campo           | Tipo    | Requerido | Descripción                  |
|-----------------|---------|-----------|------------------------------|
| `id`            | integer | Auto      | Generado automáticamente     |
| `name`          | string  | ✅        | Nombre del equipo            |
| `city`          | string  | ✅        | Ciudad del equipo            |
| `championships` | integer | ✅        | Número de campeonatos        |
| `pet`           | string  | ✅        | Mascota del equipo           |
| `arena`         | string  | ✅        | Arena donde juega el equipo  |

---

## 📡 Endpoints

### Health check de mi jugador favorito de Basketball
```
GET /api/teams/jayson
```
```json
{ "message": "Tatum" }
```

---

### GET /api/teams — Obtener todos los equipos
```bash
curl http://localhost:24596/api/teams
```
**Respuesta 200:**
```json
[
  {
    "id": 1,
    "name": "Boston Celtics",
    "city": "Boston",
    "championships": 17,
    "pet": "Lucky the Leprechaun",
    "arena": "TD Garden"
  }
]
```

---

### GET /api/teams?id=1 — Filtrar por query parameter
```bash
curl http://localhost:24596/api/teams?id=1
```
**Respuesta 200:**
```json
{
  "id": 1,
  "name": "Boston Celtics",
  "city": "Boston",
  "championships": 17,
  "pet": "Lucky the Leprechaun",
  "arena": "TD Garden"
}
```

---

### GET /api/teams/{id} — Obtener por path parameter
```bash
curl http://localhost:24596/api/teams/1
```
**Respuesta 200:**
```json
{
  "id": 1,
  "name": "Boston Celtics",
  "city": "Boston",
  "championships": 17,
  "pet": "Lucky the Leprechaun",
  "arena": "TD Garden"
}
```

---

### POST /api/teams — Crear un equipo
```bash
curl -X POST http://localhost:24596/api/teams \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Miami Heat",
    "city": "Miami",
    "championships": 3,
    "pet": "Burnie",
    "arena": "Kaseya Center"
  }'
```
**Respuesta 201:**
```json
{
  "id": 11,
  "name": "Miami Heat",
  "city": "Miami",
  "championships": 3,
  "pet": "Burnie",
  "arena": "Kaseya Center"
}
```

---

### PUT /api/teams/{id} — Reemplazar equipo completo
```bash
curl -X PUT http://localhost:24596/api/teams/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Boston Celtics",
    "city": "Boston",
    "championships": 18,
    "pet": "Lucky the Leprechaun",
    "arena": "TD Garden"
  }'
```
**Respuesta 200:** el equipo actualizado.

---

### PATCH /api/teams/{id} — Actualización parcial
```bash
curl -X PATCH http://localhost:24596/api/teams/1 \
  -H "Content-Type: application/json" \
  -d '{ "championships": 18 }'
```
**Respuesta 200:** el equipo con el campo actualizado.

---

### DELETE /api/teams/{id} — Eliminar equipo
```bash
curl -X DELETE http://localhost:24596/api/teams/1
```
**Respuesta 200:**
```json
{ "message": "Team 'Boston Celtics' (id=1) deleted successfully" }
```

---

## ❌ Manejo de errores

Todos los errores devuelven JSON con esta estructura consistente:

```json
{
  "error": "descripción corta",
  "code": 404,
  "details": "mensaje detallado"
}
```

### Caso 1 — 404 Not Found: equipo inexistente
```bash
curl http://localhost:24596/api/teams/9999
```
```json
{
  "error": "Team not found",
  "code": 404,
  "details": "No team found with the specified ID"
}
```

### Caso 2 — 422 Unprocessable Entity: campos faltantes
```bash
curl -X POST http://localhost:24596/api/teams \
  -H "Content-Type: application/json" \
  -d '{ "name": "Test" }'
```
```json
{
  "error": "Validation failed",
  "code": 422,
  "details": "Missing or empty required fields: 'city', 'pet', 'arena'"
}
```

| Código | Causa                                    |
|--------|------------------------------------------|
| 400    | JSON inválido o parámetro mal formado    |
| 404    | Equipo no encontrado                     |
| 405    | Método HTTP no soportado en el endpoint  |
| 422    | Validación fallida (campos requeridos)   |
| 500    | Error interno al guardar el archivo      |
