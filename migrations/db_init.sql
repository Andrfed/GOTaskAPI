CREATE TABLE IF NOT EXISTS public.tasks
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    title text COLLATE pg_catalog."default" NOT NULL,
    text text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT tasks_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER to postgres;
create or replace FUNCTION addTask(task_title text, task_text text)
returns table("id" integer, "title" text, "text" text)
language plpgsql
as $$
declare
begin
	return query
	insert into tasks ("title", "text") values (task_title, task_text) 
	returning tasks."id", tasks."title", tasks."text";
end;
$$;
create or replace FUNCTION getAllTasks()
returns table("id" integer, "title" text, "text" text)
language plpgsql
as $$
declare
begin
	return query
	select tasks."id", tasks."title", tasks."text" from tasks;
end;
$$;
create or replace FUNCTION getTask(id1 integer)
returns table("id" integer, "title" text, "text" text)
language plpgsql
as $$
declare
begin
	return query
	select tasks."id", tasks."title", tasks."text" from tasks
	where tasks."id" = id1;
end;
$$;
create or replace FUNCTION updateTask(id1 integer, title1 text, text1 text)
returns table("id" integer, "title" text, "text" text)
language plpgsql
as $$
declare
begin
	return query
	update tasks set
	"title" = title1, 
	"text" = text1
	where tasks."id" = id1
	returning tasks."id", tasks."title", tasks."text";
end;
$$;
create or replace procedure deleteTask(id1 integer)
language plpgsql
as $$
declare
begin
	delete from tasks
	where tasks."id" = id1;
end;
$$;