-- +goose Up
create table workout_templates (
    id uuid primary key,
    template_name text not null,
    description text not null,
    exercise_set_reps jsonb not null,
    duration_minutes integer not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose statementbegin
insert into workout_templates (
  id,
  template_name,
  description,
  exercise_set_reps,
  duration_minutes,
  created_at,
  updated_at) values
(
  -- push workouts (chest, shoulders, triceps)
  gen_random_uuid(),
'classic push day',
'traditional push workout focusing on chest, shoulders, and triceps. perfect for building upper body pressing strength.',
'{"chest": [{"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "shoulders": [{"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "triceps": [{"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}]}',
60,
now(),
now()),

(gen_random_uuid(), 
'chest & triceps blast', 
'high-volume chest and tricep workout for muscle growth. focuses on both compound and isolation movements.', 
'{"chest": [{"sets": 4, "reps_per_set": 6}, {"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}], "triceps": [{"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}]}', 
75, 
now(), 
now()),

(gen_random_uuid(), 
'shoulder focus', 
'comprehensive shoulder workout hitting all three deltoid heads plus traps for complete shoulder development.', 
'{"shoulders": [{"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 15}], "traps": [{"sets": 4, "reps_per_set": 12}]}', 
50, 
now(), 
now()),

-- pull workouts (back, biceps)
(gen_random_uuid(), 
'back & biceps power', 
'classic back and bicep workout combining heavy compounds with targeted arm work for a complete pull day.', 
'{"back": [{"sets": 4, "reps_per_set": 5}, {"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "biceps": [{"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}]}', 
70, 
now(), 
now()),

(gen_random_uuid(), 
'back thickness builder', 
'row-focused back workout designed to build thickness and density in the lats, rhomboids, and mid-back.', 
'{"back": [{"sets": 4, "reps_per_set": 6}, {"sets": 4, "reps_per_set": 8}, {"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 12}], "traps": [{"sets": 3, "reps_per_set": 15}]}', 
65, 
now(), 
now()),

(gen_random_uuid(), 
'pull-up focused back', 
'vertical pulling emphasis for building wide lats and overall back width. great for improving pull-up strength.', 
'{"back": [{"sets": 5, "reps_per_set": 5}, {"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "biceps": [{"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}]}', 
55, 
now(), 
now()),

-- leg workouts
(gen_random_uuid(), 
'leg day - quad focus', 
'squat-focused leg workout emphasizing quadriceps development with additional work for glutes and hamstrings.', 
'{"quadriceps": [{"sets": 5, "reps_per_set": 5}, {"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "hamstrings": [{"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "calves": [{"sets": 4, "reps_per_set": 15}]}', 
75, 
now(), 
now()),

(gen_random_uuid(), 
'posterior chain power', 
'deadlift-focused workout targeting hamstrings, glutes, and lower back for complete posterior chain development.', 
'{"hamstrings": [{"sets": 5, "reps_per_set": 5}, {"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "glutes": [{"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "back": [{"sets": 3, "reps_per_set": 12}], "calves": [{"sets": 3, "reps_per_set": 15}]}', 
70, 
now(), 
now()),

(gen_random_uuid(), 
'complete leg development', 
'balanced leg workout hitting quads, hamstrings, glutes, and calves with equal emphasis for proportional development.', 
'{"quadriceps": [{"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "hamstrings": [{"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}], "glutes": [{"sets": 4, "reps_per_set": 12}], "calves": [{"sets": 4, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 20}]}', 
80, 
now(), 
now()),

(gen_random_uuid(), 
'glute builder', 
'hip thrust and glute-focused workout designed to maximize glute development with supporting hamstring work.', 
'{"glutes": [{"sets": 5, "reps_per_set": 8}, {"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}], "hamstrings": [{"sets": 3, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}]}', 
60, 
now(), 
now()),

-- full body workouts
(gen_random_uuid(), 
'full body strength', 
'compound movement focused full body workout designed to build overall strength efficiently in one session.', 
'{"quadriceps": [{"sets": 4, "reps_per_set": 5}], "back": [{"sets": 4, "reps_per_set": 5}], "chest": [{"sets": 4, "reps_per_set": 5}], "shoulders": [{"sets": 3, "reps_per_set": 8}], "core": [{"sets": 3, "reps_per_set": 15}]}', 
60, 
now(), 
now()),

(gen_random_uuid(), 
'full body hypertrophy', 
'moderate rep range full body workout targeting muscle growth across all major muscle groups in one session.', 
'{"chest": [{"sets": 3, "reps_per_set": 10}], "back": [{"sets": 3, "reps_per_set": 10}], "quadriceps": [{"sets": 3, "reps_per_set": 10}], "shoulders": [{"sets": 3, "reps_per_set": 10}], "biceps": [{"sets": 2, "reps_per_set": 12}], "triceps": [{"sets": 2, "reps_per_set": 12}], "core": [{"sets": 3, "reps_per_set": 15}]}', 
75, 
now(), 
now()),

(gen_random_uuid(), 
'full body circuit', 
'high-intensity circuit style workout hitting all major muscle groups with minimal rest for conditioning and muscle endurance.', 
'{"full body": [{"sets": 3, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 20}], "cardiovascular": [{"sets": 3, "reps_per_set": 10}], "core": [{"sets": 3, "reps_per_set": 20}]}', 
45, 
now(), 
now()),

-- upper/lower splits
(gen_random_uuid(), 
'upper body push/pull', 
'complete upper body workout combining both pushing and pulling movements for balanced development.', 
'{"chest": [{"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "back": [{"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "shoulders": [{"sets": 3, "reps_per_set": 12}], "biceps": [{"sets": 3, "reps_per_set": 12}], "triceps": [{"sets": 3, "reps_per_set": 12}]}', 
80, 
now(), 
now()),

(gen_random_uuid(), 
'lower body power', 
'heavy lower body workout focusing on strength in the squat and deadlift patterns with accessory work.', 
'{"quadriceps": [{"sets": 5, "reps_per_set": 5}, {"sets": 3, "reps_per_set": 8}], "hamstrings": [{"sets": 4, "reps_per_set": 6}, {"sets": 3, "reps_per_set": 10}], "glutes": [{"sets": 3, "reps_per_set": 10}], "calves": [{"sets": 4, "reps_per_set": 12}]}', 
70, 
now(), 
now()),

-- arm specialization
(gen_random_uuid(), 
'arm annihilation', 
'high-volume arm workout with superset structure for maximum bicep and tricep pump and growth.', 
'{"biceps": [{"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 20}], "triceps": [{"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 20}], "forearms": [{"sets": 3, "reps_per_set": 15}]}', 
50, 
now(), 
now()),

-- core/abs workouts
(gen_random_uuid(), 
'core crusher', 
'comprehensive core workout targeting abs, obliques, and deep core stabilizers with varied movements.', 
'{"core": [{"sets": 3, "reps_per_set": 20}, {"sets": 3, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 30}, {"sets": 3, "reps_per_set": 12}], "obliques": [{"sets": 3, "reps_per_set": 20}, {"sets": 3, "reps_per_set": 15}]}', 
30, 
now(), 
now()),

(gen_random_uuid(), 
'six-pack sculptor', 
'ab-focused workout combining static holds, dynamic movements, and rotational work for complete core development.', 
'{"core": [{"sets": 4, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 20}, {"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 25}, {"sets": 2, "reps_per_set": 60}], "obliques": [{"sets": 3, "reps_per_set": 20}]}', 
35, 
now(), 
now()),

-- athletic/functional
(gen_random_uuid(), 
'athletic performance', 
'explosive, power-based workout using olympic lifts and plyometrics for athletic development.', 
'{"full body": [{"sets": 5, "reps_per_set": 3}, {"sets": 5, "reps_per_set": 5}, {"sets": 4, "reps_per_set": 5}], "quadriceps": [{"sets": 3, "reps_per_set": 8}], "core": [{"sets": 3, "reps_per_set": 12}]}', 
55, 
now(), 
now()),

(gen_random_uuid(), 
'functional fitness', 
'movement-based workout incorporating carries, crawls, and multi-planar exercises for real-world strength.', 
'{"full body": [{"sets": 3, "reps_per_set": 30}, {"sets": 4, "reps_per_set": 40}, {"sets": 3, "reps_per_set": 20}], "core": [{"sets": 3, "reps_per_set": 15}], "cardiovascular": [{"sets": 3, "reps_per_set": 10}]}', 
45, 
now(), 
now()),

-- beginner friendly
(gen_random_uuid(), 
'beginner full body', 
'perfect starter workout for beginners using basic movements to learn proper form and build foundational strength.', 
'{"chest": [{"sets": 3, "reps_per_set": 10}], "back": [{"sets": 3, "reps_per_set": 10}], "quadriceps": [{"sets": 3, "reps_per_set": 10}], "shoulders": [{"sets": 2, "reps_per_set": 10}], "core": [{"sets": 2, "reps_per_set": 15}]}', 
45, 
now(), 
now()),

(gen_random_uuid(), 
'beginner upper body', 
'simple upper body workout focusing on basic pushing and pulling movements with manageable volume.', 
'{"chest": [{"sets": 3, "reps_per_set": 10}, {"sets": 2, "reps_per_set": 12}], "back": [{"sets": 3, "reps_per_set": 10}, {"sets": 2, "reps_per_set": 12}], "shoulders": [{"sets": 2, "reps_per_set": 12}], "biceps": [{"sets": 2, "reps_per_set": 12}], "triceps": [{"sets": 2, "reps_per_set": 12}]}', 
40, 
now(), 
now()),

(gen_random_uuid(), 
'beginner lower body', 
'foundational leg workout teaching proper squat and hip hinge patterns with moderate volume.', 
'{"quadriceps": [{"sets": 3, "reps_per_set": 10}, {"sets": 2, "reps_per_set": 12}], "hamstrings": [{"sets": 3, "reps_per_set": 10}], "glutes": [{"sets": 2, "reps_per_set": 12}], "calves": [{"sets": 3, "reps_per_set": 15}]}', 
40, 
now(), 
now()),

-- conditioning/cardio
(gen_random_uuid(), 
'hiit conditioning', 
'high-intensity interval training combining full body movements for maximum calorie burn and conditioning.', 
'{"full body": [{"sets": 4, "reps_per_set": 20}, {"sets": 4, "reps_per_set": 15}], "cardiovascular": [{"sets": 4, "reps_per_set": 10}, {"sets": 4, "reps_per_set": 30}], "core": [{"sets": 3, "reps_per_set": 20}]}', 
30, 
now(), 
now()),

(gen_random_uuid(), 
'metcon madness', 
'metabolic conditioning workout using complexes and circuits to build work capacity and burn fat.', 
'{"full body": [{"sets": 5, "reps_per_set": 10}, {"sets": 4, "reps_per_set": 12}, {"sets": 4, "reps_per_set": 15}], "cardiovascular": [{"sets": 3, "reps_per_set": 20}]}', 
35, 
now(), 
now()),

-- bodyweight focused
(gen_random_uuid(), 
'bodyweight strength', 
'no equipment needed! challenging bodyweight workout building strength using progressive calisthenics.', 
'{"chest": [{"sets": 4, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 10}], "back": [{"sets": 4, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "quadriceps": [{"sets": 3, "reps_per_set": 15}], "core": [{"sets": 3, "reps_per_set": 20}, {"sets": 2, "reps_per_set": 60}]}', 
45, 
now(), 
now()),

(gen_random_uuid(), 
'calisthenics athlete', 
'advanced bodyweight workout for building impressive relative strength and movement skills.', 
'{"chest": [{"sets": 5, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 15}], "back": [{"sets": 5, "reps_per_set": 5}, {"sets": 3, "reps_per_set": 8}], "shoulders": [{"sets": 3, "reps_per_set": 10}], "core": [{"sets": 3, "reps_per_set": 15}, {"sets": 3, "reps_per_set": 12}]}', 
50, 
now(), 
now()),

-- time-efficient workouts
(gen_random_uuid(), 
'express workout - 30 min', 
'quick but effective full body workout perfect for busy schedules. hits all major muscle groups efficiently.', 
'{"chest": [{"sets": 3, "reps_per_set": 10}], "back": [{"sets": 3, "reps_per_set": 10}], "quadriceps": [{"sets": 3, "reps_per_set": 10}], "core": [{"sets": 2, "reps_per_set": 20}]}', 
30, 
now(), 
now()),

(gen_random_uuid(), 
'lunch break blast', 
'high-intensity 20-minute workout designed to maximize results in minimal time using supersets and circuits.', 
'{"full body": [{"sets": 3, "reps_per_set": 12}, {"sets": 3, "reps_per_set": 15}], "core": [{"sets": 2, "reps_per_set": 20}], "cardiovascular": [{"sets": 2, "reps_per_set": 10}]}', 
20, 
now(), 
now()),

-- powerlifting focused
(gen_random_uuid(), 
'powerlifting - squat day', 
'squat-focused powerlifting workout with heavy work sets and accessory movements to build massive leg strength.', 
'{"quadriceps": [{"sets": 5, "reps_per_set": 3}, {"sets": 3, "reps_per_set": 5}, {"sets": 3, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "hamstrings": [{"sets": 3, "reps_per_set": 10}], "core": [{"sets": 3, "reps_per_set": 15}]}', 
75, 
now(), 
now()),

(gen_random_uuid(), 
'powerlifting - bench day', 
'bench press focused workout with heavy sets and targeted accessory work for building pressing strength.', 
'{"chest": [{"sets": 5, "reps_per_set": 3}, {"sets": 3, "reps_per_set": 5}, {"sets": 3, "reps_per_set": 8}], "shoulders": [{"sets": 3, "reps_per_set": 10}], "triceps": [{"sets": 4, "reps_per_set": 10}, {"sets": 3, "reps_per_set": 12}]}', 
65, 
now(), 
now()),

(gen_random_uuid(), 
'powerlifting - deadlift day', 
'deadlift-focused powerlifting session with heavy pulls and back accessory work for total posterior chain strength.', 
'{"back": [{"sets": 5, "reps_per_set": 3}, {"sets": 3, "reps_per_set": 5}, {"sets": 3, "reps_per_set": 8}, {"sets": 3, "reps_per_set": 10}], "hamstrings": [{"sets": 3, "reps_per_set": 8}], "core": [{"sets": 3, "reps_per_set": 12}]}', 
70, 
now(), 
now());
-- +goose statementend

-- +goose Down
drop table workout_templates;
