-- CREACION DE CUENTA AL REGISTRAR UN NUEVO USUARIO
create or replace function createAccount()
returns trigger as $$
begin
      insert into accounts (alias, balance) values (new.alias, default);
      return new;
end;
$$ language plpgsql;

create or replace trigger tg_createAccount
after insert on users 
for each row execute function createAccount();


-- BORRAR CUENTA AL BORRAR USUARIO
create or replace function deleteAccount()
returns trigger as $$
begin
      delete from accounts where alias = old.alias; 
      return old;
end;
$$ language plpgsql;

create or replace trigger tg_deleteAccount
before delete on users
for each row execute function deleteAccount();
