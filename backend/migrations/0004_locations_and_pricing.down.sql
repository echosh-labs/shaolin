-- Revert locations and pricing changes
DROP TABLE IF EXISTS token_packages;
-- SQLite does not easily support dropping a column in older migrations, but we drop the tables
DROP TABLE IF EXISTS locations;
