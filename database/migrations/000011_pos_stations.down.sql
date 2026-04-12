-- Revertir: Eliminar función y tabla de estaciones
DROP FUNCTION IF EXISTS estaciones_ia_estacion(INTEGER, UUID, VARCHAR, VARCHAR, VARCHAR, UUID, VARCHAR, UUID);
DROP TABLE IF EXISTS estaciones_pos CASCADE;
