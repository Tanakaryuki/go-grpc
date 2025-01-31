-- ユーザーデータベース
CREATE DATABASE userdb;
CREATE USER app_user WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE userdb TO app_user;

-- プロダクトデータベース
CREATE DATABASE productdb;
CREATE USER productuser WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE productdb TO productuser;

-- オーダーデータベース
CREATE DATABASE orderdb;
CREATE USER orderuser WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE orderdb TO orderuser;

-- インベントリーデータベース
CREATE DATABASE inventorydb;
CREATE USER inventoryuser WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE inventorydb TO inventoryuser;
