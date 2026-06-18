UPDATE tables
SET qr_code_url = 'http://localhost:3000/menu?chain_id=' || b.chain_id || '&branch_id=' || tables.branch_id || '&table_id=' || tables.id
FROM branches b
WHERE tables.branch_id = b.id
  AND (tables.qr_code_url IS NULL OR tables.qr_code_url = '');
