insert into infokeeper.credit_card
(id, number, owner_id, expires_at, cvv_code, pin_code)
values ($1, $2, $3, $4, $5, $6);