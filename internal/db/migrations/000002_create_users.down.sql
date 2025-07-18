ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP INDEX IF EXISTS "accounts_owner_currency_idx";
DROP TABLE IF EXISTS "users";
