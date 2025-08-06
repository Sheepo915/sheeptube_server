CREATE TABLE
  channels (
    id BIGSERIAL PRIMARY KEY,
    channel_id UUID NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "description" TEXT,
    pic TEXT,
    banner TEXT,
    owned_by TEXT NOT NULL,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  users (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  videos (
    id BIGSERIAL PRIMARY KEY,
    video_id UUID NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "description" TEXT,
    source TEXT NOT NULL,
    poster TEXT NOT NULL,
    posted_by BIGINT NOT NULL, -- FK to channels or users table
    video_metadata BIGINT NOT NULL, -- FK to metadata table
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      CONSTRAINT fk_channel FOREIGN KEY (posted_by) REFERENCES channels (id) ON DELETE CASCADE
  );

CREATE TABLE
  video_metadata (
    id BIGSERIAL PRIMARY KEY,
    video_metadata_id UUID NOT NULL UNIQUE,
    likes BIGINT DEFAULT 0,
    views BIGINT DEFAULT 0,
    shares BIGINT DEFAULT 0,
    comments BIGINT DEFAULT 0,
    run_id UUID UNIQUE,
    encoding_status TEXT DEFAULT 'pending',
    encoding_progress INT DEFAULT 0,
    error_message TEXT,
    posted_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  video_resolutions (
    id BIGSERIAL PRIMARY KEY,
    video_metadata_id BIGSERIAL NOT NULL REFERENCES video_metadata (id) ON DELETE CASCADE,
    resolution TEXT NOT NULL,
    "status" TEXT DEFAULT 'pending',
    bitrate TEXT,
    "size" BIGINT,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  categories (
    id BIGSERIAL PRIMARY KEY,
    category_id UUID NOT NULL UNIQUE,
    category TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  video_categories (
    video_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      PRIMARY KEY (video_id, category_id),
      CONSTRAINT fk_video FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
      CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
  );

CREATE TABLE
  subtitles (
    id BIGSERIAL PRIMARY KEY,
    subtitle_id UUID NOT NULL UNIQUE,
    language TEXT NOT NULL,
    file_path TEXT NOT NULL,
    belong_to BIGINT NOT NULL,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      CONSTRAINT fk_video FOREIGN KEY (belong_to) REFERENCES videos (id) ON DELETE CASCADE
  );

CREATE TABLE
  tags (
    id BIGSERIAL PRIMARY KEY,
    tag_id UUID NOT NULL UNIQUE,
    "name" TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  video_tags (
    video_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (video_id, tag_id),
    CONSTRAINT fk_video FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
    CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
  );

CREATE TABLE
  comments (
    id BIGSERIAL PRIMARY KEY,
    comment_id UUID NOT NULL UNIQUE,
    video_id BIGINT NOT NULL,
    comment TEXT NOT NULL,
    commented_by BIGINT NOT NULL,
    likes BIGINT DEFAULT 0,
    parent_comment BIGINT,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
      FOREIGN KEY (parent_comment) REFERENCES comments (id) ON DELETE CASCADE
  );

CREATE TABLE
  actors (
    id BIGSERIAL PRIMARY KEY,
    actor_id UUID NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    bio TEXT,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW ()
  );

CREATE TABLE
  video_actors (
    video_id BIGINT NOT NULL,
    actor_id BIGINT NOT NULL,
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT NOW (),
      PRIMARY KEY (video_id, actor_id),
      CONSTRAINT fk_video FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
      CONSTRAINT fk_actor FOREIGN KEY (actor_id) REFERENCES actors (id) ON DELETE CASCADE
  );