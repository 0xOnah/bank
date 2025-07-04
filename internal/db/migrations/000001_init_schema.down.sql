-- Drop foreign keys first
ALTER TABLE "transfers" DROP CONSTRAINT "transfers_to_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT "transfers_from_account_id_fkey";
ALTER TABLE "entries" DROP CONSTRAINT "entries_account_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "transfers_from_account_id_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_idx";
DROP INDEX IF EXISTS "entries_account_id_idx";
DROP INDEX IF EXISTS "accounts_owner_idx";

-- Drop tables
DROP TABLE IF EXISTS "transfers";
DROP TABLE IF EXISTS "entries";
DROP TABLE IF EXISTS "accounts";
