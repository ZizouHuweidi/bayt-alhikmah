-- Create databases for each service
CREATE DATABASE maktaba;
CREATE DATABASE murshid;
CREATE DATABASE ory_kratos;
CREATE DATABASE ory_hydra;

-- Create service users with passwords
CREATE USER maktaba WITH PASSWORD 'maktaba';
CREATE USER murshid WITH PASSWORD 'murshid';

-- Grant privileges to service users
GRANT ALL PRIVILEGES ON DATABASE maktaba TO maktaba;
GRANT ALL PRIVILEGES ON DATABASE murshid TO murshid;

-- Connect to maktaba and grant schema permissions
\c maktaba
GRANT ALL ON SCHEMA public TO maktaba;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO maktaba;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO maktaba;

-- Connect to murshid and grant schema permissions
\c murshid
GRANT ALL ON SCHEMA public TO murshid;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO murshid;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO murshid;
