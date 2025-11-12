# Retroskb

**Retroskb** es una aplicación web fullstack desarrollada con **Go (GoFiber)** en el backend y **React.js + TailwindCSS** en el frontend.  
Su objetivo es ofrecer un sistema simple y moderno para **gestionar mangas**, con autenticación segura mediante **JWT** y operaciones **CRUD** completas.

---

## 🚀 Descripción general

Retroskb está construida bajo los principios de **arquitectura limpia**, separando las responsabilidades en capas bien definidas para mantener un código mantenible y escalable.  

El backend expone una **API REST** en Go que se comunica con una base de datos **MongoDB**, mientras que el frontend (hecho con React + Tailwind) se sirve directamente desde el mismo servidor en modo producción.

En modo desarrollo, el frontend se ejecuta con **Vite** y consume la API del backend por medio de las variables configuradas en `.env`.

---


## ⚙️ Tecnologías utilizadas

### 🖥️ Backend
- **Go 1.22+**
- **GoFiber** como framework HTTP
- **MongoDB** como base de datos
- **JWT** para autenticación
- **Arquitectura limpia**
- **Validador interno** para entrada de datos

### 💡 Frontend
- **React.js**
- **TailwindCSS**
- **Vite** (entorno de desarrollo y build de producción)
- El código del front lo encuentras [aquí](https://github.com/FabricioAsat/retroskb-client)

---

## 🧩 Funcionamiento del entorno

- En **modo desarrollo**, el frontend se levanta con Vite (`npm run dev`) y el backend con Go (`go run cmd/server/main.go`), trabajando de forma separada.
- En **modo producción**, el backend sirve automáticamente los archivos del frontend desde `web/dist`.

Esto se controla mediante una variable en `.env` llamada, `APP_ENV`, que puede ser `dev` o `prod`.

```env
# Backend
PORT=4096
MONGO_URI=mongodb://localhost:27017
MONGO_DB=retroskb
JWT_SECRET=tu_secreto_super_seguro

# Frontend
APP_ENV=dev       # usa "prod" para servir el frontend

---

## 🧱 Diseño del backend

El backend sigue los principios de **Clean Architecture**, separando responsabilidades de la siguiente forma:

- **domain** → define entidades base (`User`, `Manga`) y sus interfaces.  
- **repository** → implementa la persistencia en **MongoDB**.  
- **service** → contiene la **lógica de negocio**.  
- **transport/http** → define **endpoints**, **middlewares** y **rutas** con **GoFiber**.  
- **utils** → utilidades compartidas (validadores, helpers).  

Esta estructura facilita mantener, probar y escalar el proyecto.

---

## 🔒 Autenticación

El sistema utiliza **JWT** para el manejo de sesiones:

- Los usuarios se autentican mediante `/auth/login`.  
- El token JWT se devuelve al cliente y se envía en cada request autenticada.  
- Middlewares en `middleware.go` protegen las rutas privadas.  

---

## 📚 CRUD de mangas

La API permite **crear, listar, actualizar y eliminar mangas**.  
Estas operaciones están gestionadas en `manga_handler.go` y `manga_service.go`,  
con persistencia en `mongo_manga.go`.

---

## ⚡ Ejecución rápida

### 1️⃣ Backend
```bash
cd cmd/server
go run main.go

### 2️⃣ Frontend
```bash
cd cmd/server
go run main.go

```bash
npm run build
# El backend servirá automáticamente el contenido de web/dist


---

## 🧰 Buenas prácticas aplicadas

- Arquitectura limpia para **escalabilidad y mantenibilidad**.  
- Separación clara entre **lógica**, **transporte** y **persistencia**.  
- Uso de **GoFiber** por su rendimiento y sintaxis ligera.  
- **Variables de entorno** para diferenciar entornos `dev` / `prod`.  
- Autenticación **segura con JWT**.  
- Frontend **moderno, rápido y responsivo** con **React + TailwindCSS**.  


---

## 👨‍💻 Autor

**Fabricio Asat**  
💻 Proyecto personal — desarrollado con Go, Fiber, MongoDB, React y TailwindCSS.  
📧 [fabricioasat00@gmail.com]  
🔗 [LinkedIn](https://www.linkedin.com/in/fabricio-daniel-asat-780127237/)

---

## 📄 Licencia

Este proyecto se distribuye bajo la licencia **MIT**.
