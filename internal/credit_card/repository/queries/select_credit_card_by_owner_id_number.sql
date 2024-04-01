select
    id,
    number,
    owner_id,
    owner_name,
    expires_at,
    cvv_code,
    pin_code,
    created_at,
    metadata
from infokeeper.credit_card
where owner_id = $1 and number = $2