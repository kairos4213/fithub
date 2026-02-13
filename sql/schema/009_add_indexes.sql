-- +goose Up
create index idx_goals_user_id on goals (user_id);
create index idx_refresh_tokens_user_id on refresh_tokens (user_id);
create index idx_body_weights_user_id on body_weights (user_id);
create index idx_muscle_masses_user_id on muscle_masses (user_id);
create index idx_body_fat_percents_user_id on body_fat_percents (user_id);
create index idx_workouts_user_id on workouts (user_id);
create index idx_workouts_exercises_workout_id on workouts_exercises (workout_id);
create index idx_workouts_exercises_exercise_id on workouts_exercises (exercise_id);

-- +goose Down
drop index idx_goals_user_id;
drop index idx_refresh_tokens_user_id;
drop index idx_body_weights_user_id;
drop index idx_muscle_masses_user_id;
drop index idx_body_fat_percents_user_id;
drop index idx_workouts_user_id;
drop index idx_workouts_exercises_workout_id;
drop index idx_workouts_exercises_exercise_id;
