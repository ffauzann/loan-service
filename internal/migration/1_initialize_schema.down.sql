-- 1. Child tables of loan
DROP TABLE IF EXISTS loan_disbursement;
DROP TABLE IF EXISTS loan_investment;
DROP TABLE IF EXISTS loan_approval;

-- 2. Loan table (depends on user)
DROP TABLE IF EXISTS loan;

-- 3. Drop enum type for loan_state
DROP TYPE IF EXISTS loan_state;

-- 4. User table (depends on role)
DROP TABLE IF EXISTS "user";

-- 5. Drop enum type for user_status
DROP TYPE IF EXISTS user_status;

-- 6. Role table
DROP TABLE IF EXISTS "role";

-- 7. Drop pgcrypto extension if not needed anymore
DROP EXTENSION IF EXISTS pgcrypto;