SET TIME ZONE 'UTC';

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

CREATE TABLE IF NOT EXISTS providers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(50) NOT NULL,
  logo text NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  operation varchar(50) NOT NULL,
  provider_id UUID NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS workflows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  name VARCHAR(50) NOT NULL,
  trigger_id UUID,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS workflow_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  text VARCHAR(255),
  type VARCHAR(50),
  params JSONB,
  workflow_id UUID NOT NULL REFERENCES workflows(id),
  action_id uuid NOT NULL REFERENCES actions(id),
  next_action_id UUID,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT now(),
  version INTEGER NOT NULL DEFAULT 1
);

ALTER TABLE workflows ADD CONSTRAINT fk_workflow_workflow_actions
  FOREIGN KEY (trigger_id) REFERENCES workflow_actions(id);

ALTER TABLE workflow_actions ADD CONSTRAINT fk_workflow_actions_workflow
  FOREIGN KEY (next_action_id) REFERENCES workflow_actions(id);

-- Insert users
INSERT INTO users (id) VALUES
('550e8400-e29b-41d4-a716-446655440000'),
('550e8400-e29b-41d4-a716-446655440006');

-- Insert providers
INSERT INTO providers (id, name, logo) VALUES
('c4f9b885-2df5-4b1b-9fa4-81f87f824da8', 'System', 'image.png');

-- Insert actions
INSERT INTO actions (id, operation, provider_id, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Create', 'c4f9b885-2df5-4b1b-9fa4-81f87f824da8', now(), now()),
('550e8400-e29b-41d4-a716-446655440004', 'Update', 'c4f9b885-2df5-4b1b-9fa4-81f87f824da8', now(), now()),
('550e8400-e29b-41d4-a716-446655440007', 'Review', 'c4f9b885-2df5-4b1b-9fa4-81f87f824da8', now(), now()),
('550e8400-e29b-41d4-a716-446655440008', 'Approve', 'c4f9b885-2df5-4b1b-9fa4-81f87f824da8', now(), now());

-- Insert workflows
INSERT INTO workflows (id, user_id, name, trigger_id, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 'User Onboarding', NULL, now(), now()),
('550e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440006', 'Document Approval', NULL, now(), now());

-- Insert workflow actions
INSERT INTO workflow_actions (id, text, type, params, workflow_id, action_id, next_action_id, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440003', 'Begin Onboarding', 'Init', '{}', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440005', now(), now()),
('550e8400-e29b-41d4-a716-446655440005', 'Complete Onboarding', 'Finish', '{}', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004', NULL, now(), now()),
('550e8400-e29b-41d4-a716-446655440010', 'Submit Document', 'Init', '{}', '550e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440011', now(), now()),
('550e8400-e29b-41d4-a716-446655440011', 'Approve Document', 'Finish', '{}', '550e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440008', NULL, now(), now());

-- Update workflows for the trigger ID
UPDATE workflows SET trigger_id = '550e8400-e29b-41d4-a716-446655440003' WHERE id = '550e8400-e29b-41d4-a716-446655440002';
UPDATE workflows SET trigger_id = '550e8400-e29b-41d4-a716-446655440010' WHERE id = '550e8400-e29b-41d4-a716-446655440009';
