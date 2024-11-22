BEGIN;

ALTER TABLE reconciliation_jobs ADD COLUMN error_information VARCHAR;

END;
