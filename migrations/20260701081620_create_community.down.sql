-- 4.2.4 Откат сообщества
DROP INDEX IF EXISTS idx_follows_followed;
DROP INDEX IF EXISTS idx_follows_follower;
DROP INDEX IF EXISTS idx_comments_post;
DROP INDEX IF EXISTS idx_posts_trip_date;
DROP INDEX IF EXISTS idx_posts_type;
DROP INDEX IF EXISTS idx_posts_route;
DROP INDEX IF EXISTS idx_posts_user;

DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;

DROP TYPE IF EXISTS post_type;
