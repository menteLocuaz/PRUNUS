ALTER TABLE usuario
    DROP COLUMN IF EXISTS nombre_ticket,
    DROP COLUMN IF EXISTS usu_tarjeta_nfc,
    DROP COLUMN IF EXISTS usu_telefono;
