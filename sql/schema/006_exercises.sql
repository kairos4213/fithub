-- +goose Up
CREATE TABLE exercises (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    primary_muscle_group TEXT,
    secondary_muscle_group TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose statementbegin
INSERT INTO exercises (
    id,
    name,
    description,
    primary_muscle_group,
    secondary_muscle_group,
    created_at,
    updated_at
) VALUES
(
    gen_random_uuid(),
    'bench press',
    'barbell press for chest, shoulders, triceps. lie on bench, grip bar '
    || 'slightly wider than shoulders, lower to chest, press back up.',
    'chest',
    'triceps',
    now(),
    now()
),
(
    gen_random_uuid(),
    'back squat',
    'barbell squat for quads and glutes. place bar on upper back, stand '
    || 'feet shoulder width, lower hips until thighs are parallel, drive up.',
    'quadriceps',
    'glutes',
    now(),
    now()
),
(
    gen_random_uuid(),
    'deadlift',
    'compound lift for hamstrings, glutes, lower back. stand hip-width, '
    || 'grip bar on floor, keep back flat, extend hips to lift bar.',
    'hamstrings',
    'lower back',
    now(),
    now()
),
(
    gen_random_uuid(),
    'pull-up',
    'bodyweight pulling for lats, biceps. grip bar with palms forward, hang '
    || 'fully extended, pull chest to bar, lower under control.',
    'lats',
    'biceps',
    now(),
    now()
),
(
    gen_random_uuid(),
    'overhead press',
    'press for shoulders, triceps. hold bar at shoulders, brace core, press '
    || 'bar overhead until arms lock, lower slowly to start.',
    'shoulders',
    'triceps',
    now(),
    now()
),
(
    gen_random_uuid(),
    'bicep curl',
    'isolation for biceps. hold dumbbells at sides, palms forward, curl '
    || 'weights up while keeping elbows still, lower slowly.',
    'biceps',
    null,
    now(),
    now()
),
(
    gen_random_uuid(),
    'tricep dip',
    'dip for triceps and chest. use parallel bars, lower body until arms '
    || 'bend 90 degrees, press back to start position.',
    'triceps',
    'chest',
    now(),
    now()
),
(
    gen_random_uuid(),
    'lunge',
    'leg exercise for quads and glutes. step forward, lower back knee '
    || 'toward floor, front thigh parallel, push back up to stand.',
    'quadriceps',
    'glutes',
    now(),
    now()
),
(
    gen_random_uuid(),
    'plank',
    'core stability. hold forearms on ground, body straight from head to '
    || 'heels, engage core, avoid hips sagging or rising.',
    'abdominals',
    'lower back',
    now(),
    now()
),
(
    gen_random_uuid(),
    'mountain climber',
    'dynamic cardio for core, hip flexors. start in plank, drive knees '
    || 'alternately toward chest quickly, maintain steady rhythm.',
    'abdominals',
    'hip flexors',
    now(),
    now()
);
-- +goose statementend

-- +goose Down
DROP TABLE exercises;
