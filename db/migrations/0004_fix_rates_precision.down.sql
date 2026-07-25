DROP INDEX IF EXISTS rates_base_target_uidx;

ALTER TABLE rates
    ALTER COLUMN rate TYPE DECIMAL(10, 6);
