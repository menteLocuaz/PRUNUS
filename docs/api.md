# API - Referencia de endpoints

Base URL: `http://localhost:9090`

> Todos los endpoints que modifican datos requieren `Content-Type: application/json`

---

## Categoría

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/categoria` | Listar todas |
| POST | `/categoria` | Crear |
| GET | `/categoria/{id}` | Obtener por ID |
| PUT | `/categoria/{id}` | Actualizar |
| DELETE | `/categoria/{id}` | Eliminar |

**POST /categoria**
```json
{
  "nombre": "Electrónica",
  "id_sucursal": 1
}
```

**PUT /categoria/{id}**
```json
{
  "nombre": "Electrónica actualizada",
  "id_sucursal": 1
}
```

---

## Cliente

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/cliente` | Listar todos |
| POST | `/cliente` | Crear |
| GET | `/cliente/{id}` | Obtener por ID |
| PUT | `/cliente/{id}` | Actualizar |
| DELETE | `/cliente/{id}` | Eliminar |

**POST /cliente**
```json
{
  "empresa_cliente": "Empresa ABC S.A.",
  "nombre": "Juan Pérez",
  "ruc": "20123456789",
  "direccion": "Av. Lima 123",
  "telefono": "987654321",
  "email": "juan@empresa.com",
  "estado": 1
}
```

**PUT /cliente/{id}**
```json
{
  "empresa_cliente": "Empresa ABC S.A.",
  "nombre": "Juan Pérez",
  "ruc": "20123456789",
  "direccion": "Av. Lima 456",
  "telefono": "987654321",
  "email": "juan@empresa.com",
  "estado": 1
}
```

---

## Medida (Unidad)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/medida` | Listar todas |
| POST | `/medida` | Crear |
| GET | `/medida/{id}` | Obtener por ID |
| PUT | `/medida/{id}` | Actualizar |
| DELETE | `/medida/{id}` | Eliminar |

**POST /medida**
```json
{
  "nombre": "Kilogramo",
  "id_sucursal": 1
}
```

**PUT /medida/{id}**
```json
{
  "nombre": "Kilogramo",
  "id_sucursal": 1
}
```

---

## Moneda

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/moneda` | Listar todas |
| POST | `/moneda` | Crear |
| GET | `/moneda/{id}` | Obtener por ID |
| PUT | `/moneda/{id}` | Actualizar |
| DELETE | `/moneda/{id}` | Eliminar |

**POST /moneda**
```json
{
  "nombre": "Sol peruano",
  "id_sucursal": 1,
  "estado": 1
}
```

**PUT /moneda/{id}**
```json
{
  "nombre": "Sol peruano",
  "id_sucursal": 1,
  "estado": 1
}
```

---

## Producto

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/producto` | Listar todos |
| POST | `/producto` | Crear |
| GET | `/producto/{id}` | Obtener por ID |
| PUT | `/producto/{id}` | Actualizar |
| DELETE | `/producto/{id}` | Eliminar |

**POST /producto**
```json
{
  "nombre": "Laptop HP",
  "descripcion": "Laptop HP 15 pulgadas",
  "precio_compra": 2500.00,
  "precio_venta": 3200.00,
  "stock": 10,
  "fecha_vencimiento": "2027-01-01T00:00:00Z",
  "imagen": "https://ejemplo.com/imagen.jpg",
  "estado": 1,
  "id_sucursal": 1,
  "id_categoria": 1,
  "id_moneda": 1,
  "id_unidad": 1
}
```

**PUT /producto/{id}**
```json
{
  "nombre": "Laptop HP",
  "descripcion": "Laptop HP 15 pulgadas actualizada",
  "precio_compra": 2500.00,
  "precio_venta": 3400.00,
  "stock": 8,
  "fecha_vencimiento": "2027-01-01T00:00:00Z",
  "imagen": "https://ejemplo.com/imagen.jpg",
  "estado": 1,
  "id_sucursal": 1,
  "id_categoria": 1,
  "id_moneda": 1,
  "id_unidad": 1
}
```

---

## Proveedor

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/proveedor` | Listar todos |
| POST | `/proveedor` | Crear |
| GET | `/proveedor/{id}` | Obtener por ID |
| PUT | `/proveedor/{id}` | Actualizar |
| DELETE | `/proveedor/{id}` | Eliminar |

**POST /proveedor**
```json
{
  "nombre": "Distribuidora Norte S.A.C.",
  "ruc": "20456789012",
  "telefono": "01-4567890",
  "direccion": "Av. Industrial 789",
  "email": "contacto@distribuidora.com",
  "estado": 1,
  "id_sucursal": 1,
  "id_empresa": 1
}
```

**PUT /proveedor/{id}**
```json
{
  "nombre": "Distribuidora Norte S.A.C.",
  "ruc": "20456789012",
  "telefono": "01-4567890",
  "direccion": "Av. Industrial 999",
  "email": "contacto@distribuidora.com",
  "estado": 1,
  "id_sucursal": 1,
  "id_empresa": 1
}
```

---

## Notas

- `estado`: `1` = activo, `0` = inactivo
- `fecha_vencimiento`: formato ISO 8601 `"YYYY-MM-DDT00:00:00Z"`. Omitir el campo si no aplica
- Los campos `id_categoria`, `id_moneda`, `id_unidad`, `id_sucursal`, `id_empresa` deben existir previamente en la base de datos
- Los campos `id_*`, `created_at`, `updated_at`, `deleted_at` son generados por el servidor, no se envían en POST/PUT
