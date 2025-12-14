DROP INDEX IF EXISTS idx_defects_assignee;
DROP INDEX IF EXISTS idx_defects_due_date;
DROP INDEX IF EXISTS idx_defects_priority;
DROP INDEX IF EXISTS idx_defects_status;
DROP INDEX IF EXISTS idx_defects_project;
DROP INDEX IF EXISTS idx_projects_name;

DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS defect_history;
DROP TABLE IF EXISTS defect_comments;
DROP TABLE IF EXISTS defect_attachments;
DROP TABLE IF EXISTS defects;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS defect_severity;
DROP TYPE IF EXISTS defect_priority;
DROP TYPE IF EXISTS defect_status;
DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS "pgcrypto";
