-- Columnas adicionales de usuario referenciadas en el store pero ausentes en la tabla inicial
ALTER TABLE usuario
    ADD COLUMN IF NOT EXISTS usu_telefono    VARCHAR(20),
    ADD COLUMN IF NOT EXISTS usu_tarjeta_nfc VARCHAR(100),
    ADD COLUMN IF NOT EXISTS nombre_ticket   VARCHAR(100);
