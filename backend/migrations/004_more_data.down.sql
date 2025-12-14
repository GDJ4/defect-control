DELETE FROM defect_comments WHERE id IN (
  'eeeeeeee-cccc-cccc-cccc-cccccccccc01',
  'eeeeeeee-cccc-cccc-cccc-cccccccccc02',
  'eeeeeeee-cccc-cccc-cccc-cccccccccc03'
);

DELETE FROM defects WHERE id IN (
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee05',
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee06',
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee07',
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeee08'
);

DELETE FROM project_members WHERE project_id IN (
  'cccccccc-cccc-cccc-cccc-ccccccccccaa',
  'dddddddd-dddd-dddd-dddd-ddddddddddee'
) AND user_id IN (
  '11111111-1111-1111-1111-111111111111',
  '22222222-2222-2222-2222-222222222222'
);

DELETE FROM projects WHERE id IN (
  'cccccccc-cccc-cccc-cccc-ccccccccccaa',
  'dddddddd-dddd-dddd-dddd-ddddddddddee'
);
