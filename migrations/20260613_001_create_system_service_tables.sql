-- Migration: create tables for system_service ORM models
-- Models included: actions, countries, currencies, permissions,
-- role_permissions, roles, status

BEGIN;

-- 1) Base lookup tables ------------------------------------------------------

CREATE TABLE IF NOT EXISTS roles (
  role_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  role VARCHAR(100) NOT NULL,
  description VARCHAR(500) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_roles_role ON roles(role);

CREATE TABLE IF NOT EXISTS actions (
  action_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  action VARCHAR(50) NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_actions_action ON actions(action);

CREATE TABLE IF NOT EXISTS currencies (
  currency_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  symbol VARCHAR(20) NOT NULL,
  currency VARCHAR(50) NOT NULL,
  active INTEGER NOT NULL DEFAULT 1,
  date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_currencies_currency ON currencies(currency);

-- 2) Countries ---------------------------------------------------------------

CREATE TABLE IF NOT EXISTS countries (
  country_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  country VARCHAR(255) NOT NULL,
  description VARCHAR(500) NOT NULL DEFAULT '',
  country_code VARCHAR(20) NOT NULL,
  default_currency BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  CONSTRAINT fk_countries_default_currency
    FOREIGN KEY (default_currency)
    REFERENCES currencies(currency_id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_countries_country_code ON countries(country_code);
CREATE INDEX IF NOT EXISTS idx_countries_country ON countries(country);

-- 3) Permissions -------------------------------------------------------------

CREATE TABLE IF NOT EXISTS permissions (
  permission_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  permission VARCHAR(100) NOT NULL,
  permission_code VARCHAR(10) NOT NULL,
  permission_description VARCHAR(500) NOT NULL DEFAULT '',
  date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_permissions_permission_code ON permissions(permission_code);
CREATE INDEX IF NOT EXISTS idx_permissions_permission ON permissions(permission);

-- 4) Role permissions --------------------------------------------------------

CREATE TABLE IF NOT EXISTS role_permissions (
  role_permission_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  action_id BIGINT NOT NULL,
  date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER NOT NULL DEFAULT 0,
  modified_by INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 1,
  CONSTRAINT fk_role_permissions_role
    FOREIGN KEY (role_id)
    REFERENCES roles(role_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  CONSTRAINT fk_role_permissions_permission
    FOREIGN KEY (permission_id)
    REFERENCES permissions(permission_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  CONSTRAINT fk_role_permissions_action
    FOREIGN KEY (action_id)
    REFERENCES actions(action_id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  CONSTRAINT uq_role_permissions_triplet UNIQUE (role_id, permission_id, action_id)
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_action_id ON role_permissions(action_id);

-- 5) Status ------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS status (
  status_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  status_code VARCHAR(128) NOT NULL,
  status VARCHAR(128) NOT NULL,
  active INTEGER,
  created_by INTEGER,
  modified_by INTEGER,
  date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_status_status_code ON status(status_code);
CREATE INDEX IF NOT EXISTS idx_status_status ON status(status);

COMMIT;
