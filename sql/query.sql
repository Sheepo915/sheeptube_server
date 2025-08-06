/* Videos */
-- name: GetAllVideos :many
SELECT
  v.video_id,
  v.name,
  v.poster,
  v.source,
  vm.likes,
  c.channel_id,
  c.name,
  c.pic,
  v.created_at
FROM
  videos v
  INNER JOIN video_metadata vm ON v.video_metadata = vm.id
  INNER JOIN channels c ON v.posted_by = c.id
LIMIT
  $1
OFFSET
  $2;

-- name: GetVideoByID :one
SELECT
  v.name,
  v.poster,
  v.source,
  vm.likes,
  vm.views,
  vm.shares,
  c.channel_id,
  c.name AS channel_name,
  c.pic AS channel_pic,
  v.created_at,
  COALESCE(
    array_agg (DISTINCT cat.category) FILTER (
      WHERE
        cat.category IS NOT NULL
    ),
    '{}'
  ) AS categories,
  COALESCE(
    array_agg (DISTINCT t.name) FILTER (
      WHERE
        t.name IS NOT NULL
    ),
    '{}'
  ) AS tags,
  COALESCE(
    array_agg (DISTINCT a.name) FILTER (
      WHERE
        a.name IS NOT NULL
    ),
    '{}'
  ) AS actors
FROM
  videos v
  INNER JOIN video_metadata vm ON v.video_metadata = vm.id
  INNER JOIN channels c ON v.posted_by = c.id
  LEFT JOIN video_categories vc ON vc.video_id = v.id
  LEFT JOIN categories cat ON cat.id = vc.category_id
  LEFT JOIN video_tags vt ON vt.video_id = v.id
  LEFT JOIN tags t ON t.id = vt.tag_id
  LEFT JOIN video_actors va ON va.video_id = v.id
  LEFT JOIN actors a ON a.id = va.actor_id
WHERE
  v.video_id = $1
GROUP BY
  v.id,
  vm.likes,
  vm.views,
  vm.shares,
  c.channel_id,
  c.name,
  c.pic,
  v.created_at;

-- name: NewVideo :exec
WITH metadata AS (
  INSERT INTO video_metadata (video_metadata_id) 
  VALUES (gen_random_uuid())
  RETURNING id
)
INSERT INTO videos (
  video_id,
  "name",
  "description",
  source,
  poster,
  posted_by,
  video_metadata
) VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4,
  $5,
  (SELECT id FROM metadata)
);

/* Channel */
/* User */
/* Model */