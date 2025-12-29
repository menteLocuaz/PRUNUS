CREATE TABLE IF NOT EXISTS rol (
    id_rol       SERIAL PRIMARY KEY,
    nombre_rol   VARCHAR(100) NOT NULL,
    id_sucursal  INTEGER NOT NULL,
    estado       INTEGER NOT NULL DEFAULT 1,

    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP NULL,

    CONSTRAINT fk_rol_sucursal
        FOREIGN KEY (id_sucursal)
        REFERENCES sucursal(id_sucursal)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);


Función
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

Trigger para rol
CREATE TRIGGER trg_rol_updated
BEFORE UPDATE ON rol
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

Indices recomendados
CREATE INDEX idx_rol_id_sucursal ON rol(id_sucursal);
CREATE INDEX idx_rol_estado ON rol(estado);
CREATE INDEX idx_rol_deleted_at ON rol(deleted_at);