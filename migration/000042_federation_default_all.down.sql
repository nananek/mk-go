-- Restore the legacy 'none' default. Existing rows are not touched — going
-- back to mk-go pre-#536 leaves the field unenforced regardless of value.
ALTER TABLE "meta" ALTER COLUMN "federation" SET DEFAULT 'none';
