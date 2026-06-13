-- Migration: create tables used by system_service controllers
-- Generated for controllers: branches, countries, currencies, permissions,
-- role_permissions, roles, status
-- Includes FK dependency tables referenced by those models: users, actions

BEGIN;

-- 1) Base lookup tables ------------------------------------------------------

CREATE TABLE IF NOT EXISTS roles (
  role_id BIGSERIAL PRIMARY KEY,
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
  action_id BIGSERIAL PRIMARY KEY,
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
  currency_id BIGSERIAL PRIMARY KEY,
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
  country_id BIGSERIAL PRIMARY KEY,
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

-- 3) Users (dependency for branches.branch_manager) --------------------------

CREATE TABLE IF NOT EXISTS users (
  user_id BIGSERIAL PRIMARY KEY,
  user_details_id BIGINT,
  image_path VARCHAR(200),
  user_type INTEGER,
  full_name VARCHAR(255) NOT NULL,
  username VARCHAR(40),
  password VARCHAR(255) NOT NULL,
  email VARCHAR(255),
  phone_number VARCHAR(255),
  gender VARCHAR(10) NOT NULL,
  dob TIMESTAMP NOT NULL,
  address VARCHAR(255),
  id_type VARCHAR(5),
  id_number VARCHAR(100),
  marital_status VARCHAR(20),
  active INTEGER,
  role BIGINT,
  is_verified BOOLEAN DEFAULT FALSE,
  date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP,
  created_by INTEGER,
  modified_by INTEGER,
  CONSTRAINT fk_users_role
    FOREIGN KEY (role)
    REFERENCES roles(role_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- 4) Branches ----------------------------------------------------------------

CREATE TABLE IF NOT EXISTS branches (
  branch_id BIGSERIAL PRIMARY KEY,
  branch VARCHAR(80) NOT NULL,
  country_id BIGINT NOT NULL,
  location TEXT NOT NULL,
  phone_number VARCHAR(255) NOT NULL,
  active INTEGER DEFAULT 1,
  date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  date_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER DEFAULT 0,
  modified_by INTEGER DEFAULT 0,
  branch_manager BIGINT,
  CONSTRAINT uq_branches_branch UNIQUE (branch),
  CONSTRAINT fk_branches_country
    FOREIGN KEY (country_id)
    REFERENCES countries(country_id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
  CONSTRAINT fk_branches_branch_manager
    FOREIGN KEY (branch_manager)
    REFERENCES users(user_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_branches_country_id ON branches(country_id);
CREATE INDEX IF NOT EXISTS idx_branches_branch_manager ON branches(branch_manager);

-- 5) Permissions -------------------------------------------------------------

CREATE TABLE IF NOT EXISTS permissions (
  permission_id BIGSERIAL PRIMARY KEY,
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

-- 6) Role permissions --------------------------------------------------------

CREATE TABLE IF NOT EXISTS role_permissions (
  role_permission_id BIGSERIAL PRIMARY KEY,
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

-- 7) Status ------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS status (
  status_id BIGSERIAL PRIMARY KEY,
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
