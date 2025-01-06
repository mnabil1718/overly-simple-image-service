 CREATE TABLE IF NOT EXISTS images (
 id bigserial PRIMARY KEY,
 name text UNIQUE NOT NULL,
 alt text NOT NULL,
 file_name text NOT NULL,
 size integer NOT NULL,
 width integer NOT NULL,
 height integer NOT NULL,
 mime_type text NOT NULL,
 is_temp boolean NOT NULL DEFAULT true,
 created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
 updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
 version integer NOT NULL DEFAULT 1
 );