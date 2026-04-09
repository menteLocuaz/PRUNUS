# Cómo Configurar y Probar el Sistema de Autenticación

Esta guía describe los pasos para inicializar el sistema de seguridad de Prunus, configurar los módulos de navegación y validar los tres métodos de acceso disponibles (Email, Username y PIN).

---

## 1. Configuración del Entorno

Asegúrate de que tu archivo `.env` contenga las claves de seguridad necesarias:

```env
# Configuración JWT
JWT_SECRET=tu_clave_secreta_de_32_caracteres
JWT_EXPIRATION_HOURS=24
```

---

## 2. Preparación de la Base de Datos

### Opción A: Inicialización Automática (Recomendado)
El sistema incluye comandos que automatizan la creación de tablas, módulos y permisos.

1.  **Ejecutar Migraciones:** Crea el esquema.
    ```bash
    go run cmd/main.go migrate
    ```
2.  **Sembrado (Seed):** Carga rutas, iconos y permisos para el rol "Administrador".
    ```bash
    go run cmd/main.go seed
    ```
3.  **Registro de Usuario:** Crea el administrador inicial (Password: `Admin123`).
    ```bash
    go run cmd/main.go register usuario \
      --sucursal "22222222-2222-4222-a222-222222222222" \
      --rol "b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e" \
      --email "admin@prunus.com" \
      --username "admin" \
      --nombre "Admin Master" \
      --dni "12345678" \
      --password "Admin123" \
      --status "3a99d245-b34f-48a5-ac08-a5a010c5822f"
    ```

---

### Opción B: Inicialización Manual vía SQL (Directo en DB)
Si prefieres usar un cliente como DBeaver o psql, ejecuta este script para montar la base de datos completa:

```sql
-- 1. Asegurar Datos Base (Empresa, Sucursal, Rol)
INSERT INTO empresa (id_empresa, nombre, rut, id_status)
VALUES ('11111111-1111-4111-a111-111111111111', 'Empresa Demo', '12345678-9', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT DO NOTHING;

INSERT INTO sucursal (id_sucursal, id_empresa, nombre_sucursal, id_status)
VALUES ('22222222-2222-4222-a222-222222222222', '11111111-1111-4111-a111-111111111111', 'Sucursal Central', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT DO NOTHING;

INSERT INTO rol (id_rol, nombre_rol, id_sucursal, id_status)
VALUES ('b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e', 'Administrador', '22222222-2222-4222-a222-222222222222', '7f7b0e11-1234-4a21-9591-316279f06742')
ON CONFLICT DO NOTHING;

-- 2. Cargar Módulos y Rutas
INSERT INTO modulo (mdl_id, mdl_descripcion, id_status, is_active, abreviatura, nivel, orden, ruta, icono)
VALUES 
    (1, 'Configuración Empresa', '7f7b0e11-1234-4a21-9591-316279f06742', true, 'EMP',  0, 1, '/config/empresa', 'settings'),
    (2, 'Gestión Sucursales',    '7f7b0e11-1234-4a21-9591-316279f06742', true, 'SUC',  0, 2, '/config/sucursales', 'store'),
    (3, 'Usuarios y Roles',      '7f7b0e11-1234-4a21-9591-316279f06742', true, 'USR',  0, 3, '/config/usuarios', 'users'),
    (4, 'Catálogo Productos',    '7f7b0e11-1234-4a21-9591-316279f06742', true, 'PROD', 0, 4, '/productos', 'package'),
    (5, 'Ventas y POS',          '7f7b0e11-1234-4a21-9591-316279f06742', true, 'VENT', 0, 5, '/ventas', 'shopping-cart'),
    (8, 'Control de Caja',       '7f7b0e11-1234-4a21-9591-316279f06742', true, 'CAJA', 0, 6, '/caja', 'monitor')
ON CONFLICT (mdl_id) DO UPDATE SET ruta = EXCLUDED.ruta, icono = EXCLUDED.icono;

-- 3. Asignar Permisos Totales al Administrador
INSERT INTO permiso_rol (id_rol, id_modulo, can_read, can_write, can_update, can_delete)
SELECT 'b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e', id_modulo, true, true, true, true
FROM modulo
ON CONFLICT DO NOTHING;

-- 4. Crear Usuario Administrador (Password: Admin123, PIN: 1234)
INSERT INTO usuario (id_usuario, id_sucursal, id_rol, email, username, usu_nombre, usu_dni, password, usu_pin_pos, id_status)
VALUES (
  'c5d6e7f8-a1b2-4c3d-9e4f-5a6b7c8d9e0f',
  '22222222-2222-4222-a222-222222222222',
  'b4c5d6e7-f8a1-4b2c-9d3e-4f5a6b7c8d9e',
  'admin@prunus.com',
  'admin',
  'Admin Master',
  '12345678',
  '$2a$10$U.sUS/qwAXlDPrJZ9wAaLe78DmRtcnWVY39wFp85YLiL0iIVPVkkK', -- Hash de Admin123
  '1234',
  '3a99d245-b34f-48a5-ac08-a5a010c5822f' -- Estatus Activo
) ON CONFLICT DO NOTHING;
```

---

## 3. Pruebas de Login (3 Métodos)

El endpoint único es `POST /api/v1/login`.

*   **Email:** `{"email": "admin@prunus.com", "password": "Admin123"}`
*   **Username:** `{"username": "admin", "password": "Admin123"}`
*   **PIN:** `{"pin": "1234"}`

---

## 4. Validación de Permisos (Frontend)

Tras el login, el array de `permisos` entregará las rutas habilitadas. Ejemplo:
`"permisos": ["/config/empresa", "/productos", "/ventas"]`

Utiliza este array en tu Frontend para renderizar dinámicamente el Sidebar.
