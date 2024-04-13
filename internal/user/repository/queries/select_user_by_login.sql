select id, login, password, created_at
from infokeeper.user
where login = $1