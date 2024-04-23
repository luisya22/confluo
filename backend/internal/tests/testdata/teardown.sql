ALTER TABLE workflows DROP CONSTRAINT IF EXISTS fk_workflow_workflow_actions;
ALTER TABLE workflow_actions DROP CONSTRAINT IF EXISTS fk_workflow_actions_workflow;

DROP TABLE IF EXISTS workflow_actions;
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS providers;
DROP TABLE IF EXISTS users;
