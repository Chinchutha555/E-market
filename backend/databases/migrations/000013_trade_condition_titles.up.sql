ALTER TABLE market_trade_conditions ADD COLUMN title VARCHAR(255) NOT NULL DEFAULT '';
UPDATE market_trade_conditions SET title = CASE trade_type
 WHEN 'trial' THEN 'ทดลองขาย'
 WHEN 'temporary' THEN 'ทำการค้าแบบชั่วคราว'
 WHEN 'contract' THEN 'ทำการค้าแบบมีสัญญา'
 ELSE trade_type END WHERE title = '';