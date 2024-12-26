create extension if not exists pgcrypto;

CREATE ROLE readaccess;
GRANT CONNECT ON DATABASE studios TO readaccess;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readaccess;
CREATE USER studios_ro WITH PASSWORD 'studios';
GRANT readaccess TO studios_ro;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readaccess;

create table if not exists "user"
(
    id serial primary key not null ,
    login text not null ,
    password text not null ,
    role text not null ,
    first_name text not null ,
    second_name text not null ,
    third_name text not null
);

create table if not exists studio
(
    id serial primary key not null ,
    name text not null
);

create table if not exists Room
(
    id serial primary key not null,
    name text not null ,
    studio_id int references studio(id) on delete cascade not null,
    start_hour int not null ,
    end_hour int not null
);

create table if not exists equipment
(
    id serial primary key not null ,
    name text not null ,
    type int not null ,
    studio_id int references studio(id) on delete cascade not null
);

create table if not exists producer
(
    id serial primary key not null ,
    name text not null ,
    studio_id int references studio(id) on delete cascade not null ,
    start_hour int not null,
    end_hour int not null
);

create table if not exists instrumentalist
(
    id serial primary key not null ,
    name text not null ,
    studio_id int references studio(id) on delete cascade not null ,
    start_hour int not null ,
    end_hour int not null
);

create table if not exists reserve
(
    id serial primary key not null ,
    user_id int not null,
    room_id int references room(id) on delete cascade not null,
    producer_id int not null,
    instrumentalist_id int not null,
    start_time timestamp not null ,
    end_time timestamp not null
);

create table if not exists reserved_equipments
(
    id serial primary key not null ,
    reserve_id int references reserve(id) on delete cascade not null ,
    equipment_id int references equipment(id) on delete cascade not null
);

-- insert into "user"(login, password, role, first_name, second_name, third_name) values ('admin', 'admin', 'admin', 'admin','admin', 'admin')


create or replace procedure update_reserved_equipment(reserve_id integer, equipment_id integer) as '
        insert into reserved_equipments(reserve_id, equipment_id) values (reserve_id, equipment_id);
    ' language sql;

create or replace procedure delete_reserved_equipment(reserveId integer, equipmentId integer) as '
    delete from reserved_equipments where reserved_equipments.reserve_id = reserveId and reserved_equipments.equipment_id = equipmentId;
' language sql;

create or replace function is_intersect(reserveStartTime timestamp, 
					reserveEndTime timestamp,
					choosenStartTime timestamp,
					choosenEndTime timestamp) returns boolean language plpgsql as $$
    declare
        ans boolean;
    begin
    ans = false;
    if ((choosenStartTime >= reserveStartTime and choosenStartTime < reserveEndTime) or
		(choosenEndTime <= reserveEndTime and choosenEndTime > reserveStartTime) or
		(choosenStartTime <= reserveStartTime and choosenEndTime >= reserveEndTime)) then
		ans = true;
		end if;
    return ans;
	end;
    $$;

create or replace function is_reserve(userId integer,
                                      roomId integer,
                                      producerId integer,
                                      instrumentalistId integer,
                                      startTime timestamp,
                                      endTime timestamp) returns integer language plpgsql as $$
    declare
        count integer;
    begin
        select count(*) into count
        from reserve
        where (user_id = userId or
               room_id = roomId or
               producer_id = producerId or
               instrumentalist_id = instrumentalistId)
          and is_intersect(reserve.start_time,
                           reserve.end_time,
                           startTime,
                           endTime);

        return count;
    end;
    $$;

-- studio fill

insert into studio(name) values ('first_studio');
insert into studio(name) values ('second_studio');
insert into studio(name) values ('third_studio');

-- room fill

insert into room(name, studio_id, start_hour, end_hour) values ('red', 1, 9, 21);
insert into room(name, studio_id, start_hour, end_hour) values ('blue', 1, 9, 21);
insert into room(name, studio_id, start_hour, end_hour) values ('green', 2, 9, 21);

-- producer fill

insert into producer(name, studio_id, start_hour, end_hour) values ('nik', 1, 9, 21);
insert into producer(name, studio_id, start_hour, end_hour) values ('ivan', 1, 9, 21);
insert into producer(name, studio_id, start_hour, end_hour) values ('masha', 2, 9, 21);

-- instrumentalist fill

insert into instrumentalist(name, studio_id, start_hour, end_hour) values ('slash', 1, 9, 21);
insert into instrumentalist(name, studio_id, start_hour, end_hour) values ('page', 1, 9, 21);
insert into instrumentalist(name, studio_id, start_hour, end_hour) values ('jack', 1, 9, 21);

--equipment fill

insert into equipment(name, studio_id, type) values ('les paul', 1, 2);
insert into equipment(name, studio_id, type) values ('telecaster', 1, 2);
insert into equipment(name, studio_id, type) values ('strat', 1, 2);

insert into "user"(id, login, password, role, first_name, second_name, third_name) values
           ('1', 'admin', '$2a$10$53A6AVUmQFr3nx3taVYnjOKQpHh2JTeBxEaIX1NYEh2I9nnepfVPC', 'admin', 'admin', 'admin', 'admin');
