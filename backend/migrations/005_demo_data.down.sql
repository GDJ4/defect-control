DELETE FROM defect_history
WHERE id IN (
  'aaaa0000-b111-c222-d333-eeeeffff6601',
  'aaaa0000-b111-c222-d333-eeeeffff6602',
  'aaaa0000-b111-c222-d333-eeeeffff6603',
  'aaaa0000-b111-c222-d333-eeeeffff6604'
);

DELETE FROM defect_comments
WHERE id IN (
  'aaaa9999-aaaa-bbbb-cccc-ddddeeeeff00',
  'bbbb9999-aaaa-bbbb-cccc-ddddeeeeff00',
  'cccc9999-aaaa-bbbb-cccc-ddddeeeeff00',
  'dddd9999-aaaa-bbbb-cccc-ddddeeeeff00',
  'eeee9999-aaaa-bbbb-cccc-ddddeeeeff00',
  'ffff9999-aaaa-bbbb-cccc-ddddeeeeff00'
);

DELETE FROM defects
WHERE id IN (
  'aaaa0000-1111-2222-3333-444455556666',
  'bbbb0000-1111-2222-3333-444455556666',
  'cccc0000-1111-2222-3333-444455556666',
  'dddd0000-1111-2222-3333-444455556666',
  'eeee0000-1111-2222-3333-444455556666',
  'ffff0000-1111-2222-3333-444455556666'
);

DELETE FROM project_members
WHERE project_id IN (
  'aaaa1111-2222-3333-4444-555566667777',
  'bbbb1111-2222-3333-4444-555566667777',
  'cccc1111-2222-3333-4444-555566667777'
);

DELETE FROM projects
WHERE id IN (
  'aaaa1111-2222-3333-4444-555566667777',
  'bbbb1111-2222-3333-4444-555566667777',
  'cccc1111-2222-3333-4444-555566667777'
);

DELETE FROM users
WHERE id IN (
  '44444444-4444-4444-4444-444444444444',
  '55555555-5555-5555-5555-555555555555',
  '66666666-6666-6666-6666-666666666666'
);
