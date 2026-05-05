# Proyecto 2 BDD1
---

## Tecnologías utilizadas
- **Frontend:** React
- **Backend:** GO
- **Base de datos:** PostgreSQL

---

## Requisitos previos

Antes de correr el proyecto, asegurate de tener instalado:

- [Docker](https://www.docker.com/get-started) 
- [Docker Compose](https://docs.docker.com/compose/install/) 

---

## Configuración inicial

1. **Clonar el repositorio:**

```bash
git clone https://github.com/Sofilayerdi/Proyecto2_BDD1.git
cd Proyecto2_BDD1
```

2. **Crear el archivo de variables de entorno:**

En la raíz del proyecto, crear un archivo llamado `.env` con el siguiente contenido:

```env
DB_USER=proy2
DB_PASSWORD=secret
DB_NAME=proyecto2
DB_HOST=db
DB_PORT=5432
```

## Correr el proyecto

Desde la raíz del proyecto, ejecutar:

```bash
docker compose up --build
```

Esto levantará tres servicios:

| Servicio   | Puerto local |
|------------|-------------|
| `db`       | `5435`      |
| `backend`  | `8000`      |
| `frontend` | `5173`      |

Una vez que todos los contenedores estén corriendo, abrir navegador en:

```
http://localhost:5173
```

---

## Detener el proyecto

```bash
docker compose down
```

Para detener **y eliminar los volúmenes** (esto borrará los datos de la base de datos):

```bash
docker compose down -v
```