ALTER TABLE users
DROP CHECK users_chk_1;

ALTER TABLE users
ADD CONSTRAINT users_chk_1
CHECK (role IN ('admin', 'manager', 'employee'));