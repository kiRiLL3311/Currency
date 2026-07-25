-- DECIMAL(10,6) only allows 4 digits before the decimal point.
-- Several FX rates vs USD exceed that (IRR, VND, IDR, LBP, …)
-- and abort the whole sync with "numeric field overflow".
ALTER TABLE rates
    ALTER COLUMN rate TYPE DECIMAL(20, 8);

-- Required for SaveRate's ON CONFLICT (base_currency, target_currency).
CREATE UNIQUE INDEX IF NOT EXISTS rates_base_target_uidx
    ON rates (base_currency, target_currency);
