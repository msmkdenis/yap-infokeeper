insert into infokeeper.credit_card
(id, number, owner_id, owner_name, expires_at, cvv_code, pin_code, metadata)
values ($1, $2, $3, $4, $5, $6, $7, $8);