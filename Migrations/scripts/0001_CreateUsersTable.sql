CREATE TABLE "Users" (
    "Id" UUID PRIMARY KEY,
    "Email" TEXT NOT NULL UNIQUE,
    "PasswordHash" TEXT NULL,
    "GoogleId" TEXT NULL,
    "FullName" TEXT NULL,
    "CreatedAt" TIMESTAMP WITH TIME ZONE NOT NULL
);
