-- Habilita la extensión para la generación de UUID si aún no está habilitada
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Crea la tabla users, si no existe
CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);
