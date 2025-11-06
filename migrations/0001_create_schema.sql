-- 创建或修改UINT256域：
-- 这个域基于NUMERIC类型，但有限制：值必须大于等于0，小于2^256，并且不能有小数部分（即整数）。
-- 如果还没有创建这个域，就创建它。
-- 如果已经存在，则删除原有的检查约束（uint256_check）并添加新的检查约束（注意：这里假设原有的约束叫uint256_check，但实际上在创建域时如果没有指定约束名，系统会自动生成一个，这里可能存在问题，因为这里硬编码了约束名。不过，如果之前也是用同样的脚本创建的，那么约束名就是uint256_check）。

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'uint256') THEN
            CREATE DOMAIN UINT256 AS NUMERIC
                CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
        ELSE
            ALTER DOMAIN UINT256 DROP CONSTRAINT uint256_check;
            ALTER DOMAIN UINT256 ADD
                CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
        END IF;
    END $$;

-- 创建block_headers表：
-- 存储区块头信息。
-- 包含字段：hash（主键）、parent_hash（非空且唯一）、number（区块高度:非空且唯一，使用UINT256类型）、timestamp（非空且唯一，且大于0）、rlp_bytes（区块的RLP编码原始数据 非空）。
-- 在timestamp和number字段上创建索引。
CREATE TABLE IF NOT EXISTS block_headers (
                                             hash        VARCHAR   PRIMARY KEY,
                                             parent_hash VARCHAR NOT NULL UNIQUE,
                                             number      UINT256 NOT NULL UNIQUE,
                                             timestamp   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0),
                                             rlp_bytes   VARCHAR NOT NULL
);
CREATE INDEX IF NOT EXISTS block_headers_timestamp ON block_headers(timestamp);
CREATE INDEX IF NOT EXISTS block_headers_number ON block_headers(number);

-- event_blocks表：
CREATE TABLE IF NOT EXISTS event_blocks(
                                           guid        VARCHAR PRIMARY KEY,
                                           hash        VARCHAR NOT NULL UNIQUE,
                                           parent_hash VARCHAR NOT NULL UNIQUE,
                                           number      UINT256 NOT NULL UNIQUE,
                                           timestamp   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS event_blocks_timestamp ON event_blocks(timestamp);
CREATE INDEX IF NOT EXISTS event_blocks_number ON event_blocks(number);

-- 创建contract_events表：
-- 存储合约事件
-- 包含字段：guid（主键）、block_hash（外键，引用block_headers的hash，级联删除）、contract_address、transaction_hash、log_index、event_signature、timestamp（大于0）、rlp_bytes。
-- 在timestamp、block_hash、event_signature、contract_address字段上创建索引。
CREATE TABLE IF NOT EXISTS contract_events (
                                               guid             VARCHAR PRIMARY KEY,
                                               block_hash       VARCHAR NOT NULL REFERENCES block_headers(hash) ON DELETE CASCADE,
                                               contract_address VARCHAR NOT NULL,
                                               transaction_hash VARCHAR NOT NULL,
                                               log_index        INTEGER NOT NULL,
                                               event_signature  VARCHAR NOT NULL,
                                               timestamp        INTEGER NOT NULL CHECK (timestamp > 0),
                                               rlp_bytes        VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS contract_events_timestamp ON contract_events(timestamp);
CREATE INDEX IF NOT EXISTS contract_events_block_hash ON contract_events(block_hash);
CREATE INDEX IF NOT EXISTS contract_events_event_signature ON contract_events(event_signature);
CREATE INDEX IF NOT EXISTS contract_events_contract_address ON contract_events(contract_address);


-- 创建deposit_tokens表：
CREATE TABLE IF NOT EXISTS deposit_tokens (
                                              guid                          VARCHAR PRIMARY KEY,
                                              block_number                  UINT256 NOT NULL,
                                              token_address                 VARCHAR NOT NULL,
                                              sender                        VARCHAR NOT NULL,
                                              amount                        UINT256,
                                              timestamp                     INTEGER NOT NULL CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS deposit_tokens_sender ON deposit_tokens(sender);
CREATE INDEX IF NOT EXISTS deposit_tokens_timestamp ON deposit_tokens(timestamp);

CREATE TABLE IF NOT EXISTS withdraw_tokens (
                                               guid                          VARCHAR PRIMARY KEY,
                                               block_number                  UINT256 NOT NULL,
                                               token_address                 VARCHAR NOT NULL,
                                               sender                        VARCHAR NOT NULL,
                                               receiver                      VARCHAR NOT NULL,
                                               amount                        UINT256,
                                               timestamp                     INTEGER NOT NULL CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS withdraw_tokens_token_address ON withdraw_tokens(token_address);
CREATE INDEX IF NOT EXISTS withdraw_tokens_timestamp ON withdraw_tokens(timestamp);


CREATE TABLE IF NOT EXISTS grant_reward_tokens (
                                                   guid                          VARCHAR PRIMARY KEY,
                                                   block_number                  UINT256 NOT NULL,
                                                   token_address                 VARCHAR NOT NULL,
                                                   granter                       VARCHAR NOT NULL,
                                                   amount                        UINT256,
                                                   timestamp                     INTEGER NOT NULL CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS grant_reward_tokens_token_address ON grant_reward_tokens(token_address);
CREATE INDEX IF NOT EXISTS grant_reward_tokens_timestamp ON grant_reward_tokens(timestamp);


CREATE TABLE IF NOT EXISTS withdraw_manager_update (
                                                       guid                          VARCHAR PRIMARY KEY,
                                                       block_number                  UINT256 NOT NULL,
                                                       withdraw_manager              VARCHAR NOT NULL,
                                                       timestamp                     INTEGER NOT NULL CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS withdraw_manager_update_withdraw_manager ON withdraw_manager_update(withdraw_manager);
CREATE INDEX IF NOT EXISTS withdraw_manager_update_timestamp ON withdraw_manager_update(timestamp);