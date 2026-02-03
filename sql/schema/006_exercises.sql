-- +goose Up
create table exercises (
    id uuid primary key,
    name text not null,
    description text,
    primary_muscle_group text,
    secondary_muscle_group text,
    created_at timestamp not null,
    updated_at timestamp not null
);

create index idx_exercises_primary_muscle on exercises (primary_muscle_group);
create index idx_exercises_secondary_muscle on exercises (
    secondary_muscle_group
);

-- +goose statementbegin

-- chest exercises (15 exercises)
insert into exercises (id, name, description, primary_muscle_group, secondary_muscle_group, created_at, updated_at) values
(gen_random_uuid(), 'barbell bench press', 'lie on a flat bench with feet firmly planted on the ground. grip the barbell slightly wider than shoulder-width apart. lower the bar to your mid-chest in a controlled manner, keeping your elbows at about 45 degrees from your body. press the bar back up explosively until your arms are fully extended. keep your shoulder blades retracted throughout the movement.', 'chest', 'triceps', now(), now()),

(gen_random_uuid(), 'incline barbell bench press', 'set a bench to a 30-45 degree incline. lie back and grip the barbell slightly wider than shoulder-width. lower the bar to your upper chest, then press back up. this variation emphasizes the upper pectorals more than the flat bench press. keep your feet planted and avoid bouncing the bar off your chest.', 'chest', 'shoulders', now(), now()),

(gen_random_uuid(), 'decline barbell bench press', 'set a bench to a 15-30 degree decline and secure your legs. grip the barbell slightly wider than shoulder-width. lower to your lower chest, then press back up. this variation emphasizes the lower pectorals. have a spotter help you get the bar into position.', 'chest', 'triceps', now(), now()),

(gen_random_uuid(), 'dumbbell bench press', 'lie on a flat bench holding dumbbells at shoulder level. press the dumbbells up until arms are extended, bringing them slightly together at the top. lower with control back to the starting position. dumbbells allow for a greater range of motion and independent arm movement compared to barbells.', 'chest', 'triceps', now(), now()),

(gen_random_uuid(), 'incline dumbbell press', 'set a bench to a 30-45 degree incline. sit with dumbbells on your thighs, then lie back and position them at shoulder level. press the dumbbells up and slightly together until arms are extended. lower slowly back to starting position. this targets the upper portion of the chest.', 'chest', 'shoulders', now(), now()),

(gen_random_uuid(), 'decline dumbbell press', 'set a bench to a 15-30 degree decline and secure your legs. start with dumbbells at chest level. press them up until arms are extended. lower with control. this emphasizes the lower chest while allowing for natural arm movement and a deep stretch.', 'chest', 'triceps', now(), now()),

(gen_random_uuid(), 'push-ups', 'start in a plank position with hands slightly wider than shoulder-width. keep your body in a straight line from head to heels. lower your body until your chest nearly touches the ground, keeping elbows at about 45 degrees. push back up to starting position. maintain core tension throughout.', 'chest', 'triceps', now(), now()),

(gen_random_uuid(), 'diamond push-ups', 'start in a push-up position but place your hands close together so your thumbs and index fingers form a diamond shape. lower your body while keeping elbows close to your sides. this variation places greater emphasis on the triceps and inner chest. push back up to starting position.', 'triceps', 'chest', now(), now()),

(gen_random_uuid(), 'wide grip push-ups', 'start in a push-up position with hands placed significantly wider than shoulder-width. lower your body until chest nearly touches the ground. this variation emphasizes the chest muscles more than standard push-ups. push back up to starting position.', 'chest', 'shoulders', now(), now()),

(gen_random_uuid(), 'dumbbell flyes', 'lie on a flat bench holding dumbbells above your chest with a slight bend in your elbows. lower the weights out to the sides in an arc motion until you feel a stretch in your chest. bring the weights back together above your chest, squeezing your pecs. keep the elbow angle constant throughout.', 'chest', 'shoulders', now(), now()),

(gen_random_uuid(), 'incline dumbbell flyes', 'set a bench to 30-45 degrees. lie back holding dumbbells above your chest with slightly bent elbows. lower weights out to sides in an arc until you feel a stretch. bring them back together at the top. this targets the upper chest fibers more effectively than flat flyes.', 'chest', 'shoulders', now(), now()),

(gen_random_uuid(), 'cable crossover', 'stand between two cable stations set at shoulder height. grab handles and step forward with a slight lean. with elbows slightly bent, bring hands together in front of your chest in an arcing motion. squeeze your pecs at the point of maximum contraction. return to starting position with control.', 'chest', null, now(), now()),

(gen_random_uuid(), 'low to high cable flyes', 'set cable pulleys at the lowest position. stand between them and grab handles. with slight bend in elbows, raise arms up and together in front of your upper chest. this variation targets the upper chest. lower with control back to starting position.', 'chest', 'shoulders', now(), now()),

(gen_random_uuid(), 'chest dips', 'position yourself on parallel bars. lower your body by bending elbows while leaning forward slightly (about 30 degrees). lower until you feel a stretch in your chest. press back up. the forward lean shifts emphasis from triceps to chest. keep movements controlled.', 'chest', 'triceps', now(), now()),

(gen_random_uuid(), 'machine chest press', 'sit at a chest press machine with back firmly against the pad. grip handles at chest level. press handles forward until arms are extended. return with control. machines provide a stable pressing motion, good for beginners or isolation work after free weights.', 'chest', 'triceps', now(), now()),

-- back exercises (20 exercises)
(gen_random_uuid(), 'conventional deadlift', 'stand with feet hip-width apart, barbell over mid-foot. bend at hips and knees to grip the bar with hands just outside your legs. keep your back flat, chest up, and shoulders back. drive through your heels to stand up, extending hips and knees simultaneously. lower the bar back down with control. this is a compound movement working the entire posterior chain.', 'back', 'hamstrings', now(), now()),

(gen_random_uuid(), 'sumo deadlift', 'stand with feet wider than shoulder-width, toes pointed out. grip the barbell with hands inside your legs. keep chest up and back flat. drive through heels to stand up, keeping the bar close to your body. this variation reduces lower back stress and emphasizes glutes and inner thighs more.', 'back', 'glutes', now(), now()),

(gen_random_uuid(), 'barbell row', 'bend forward at the hips with a slight knee bend, keeping your back straight. grip the barbell with hands shoulder-width apart. pull the bar to your lower chest/upper abdomen, leading with your elbows. squeeze your shoulder blades together at the top. lower with control. keep your core tight throughout.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'pendlay row', 'bend forward at the hips until your torso is nearly parallel to the ground. grip the barbell with arms fully extended, bar resting on the ground. explosively pull the bar to your lower chest, then lower it back to the ground between each rep. this variation eliminates momentum and builds explosive power.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 't-bar row', 'straddle a t-bar row setup or barbell in a landmine. grip the handle with both hands. bend at hips with slight knee bend, keeping back straight. pull the weight toward your chest, squeezing shoulder blades together. lower with control. this provides a comfortable grip angle for heavy rowing.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'pull-ups', 'hang from a pull-up bar with hands slightly wider than shoulder-width, palms facing away. pull yourself up until your chin clears the bar, focusing on driving your elbows down. lower yourself with control to full extension. avoid swinging or using momentum.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'chin-ups', 'hang from a pull-up bar with hands shoulder-width apart, palms facing toward you (underhand grip). pull yourself up until your chin clears the bar. lower with control. this variation places more emphasis on the biceps compared to standard pull-ups.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'wide grip pull-ups', 'hang from a pull-up bar with hands significantly wider than shoulder-width. pull yourself up, bringing the bar toward your upper chest. this wider grip emphasizes the lats more and reduces bicep involvement. lower with control to full extension.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'neutral grip pull-ups', 'use parallel bars or neutral grip handles. hang with palms facing each other. pull yourself up until your chin clears the bars. this grip is easier on the shoulders and provides a balanced lat and bicep workout. lower with control.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'lat pulldown', 'sit at a lat pulldown machine with thighs secured under pads. grip the bar wider than shoulder-width. pull the bar down to your upper chest, leading with your elbows and squeezing your shoulder blades together. return to starting position with control. avoid leaning back excessively.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'close grip lat pulldown', 'sit at a lat pulldown machine. use a close grip attachment or grip the bar with hands 6-8 inches apart. pull down to your upper chest. this variation increases range of motion and emphasizes the lower lats. return to starting position with control.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'seated cable row', 'sit at a cable row machine with feet on the footrests and knees slightly bent. grip the handle with arms extended. pull the handle to your lower chest, keeping elbows close to your body. squeeze your shoulder blades together. return to start position with control, feeling the stretch in your lats.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'single arm dumbbell row', 'place one knee and hand on a bench for support. hold a dumbbell in the opposite hand, arm hanging straight down. pull the dumbbell up toward your hip, leading with your elbow. keep your back flat throughout. lower with control. this allows you to focus on each side independently.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'chest supported row', 'lie face down on an inclined bench set to 30-45 degrees. hold dumbbells or grip handles hanging straight down. pull the weights up toward your chest, squeezing your shoulder blades together. this variation removes lower back stress and prevents cheating with momentum.', 'back', null, now(), now()),

(gen_random_uuid(), 'face pulls', 'set a cable at upper chest height with a rope attachment. pull the rope toward your face, separating the ends as you pull and bringing them beside your ears. focus on squeezing your rear delts and upper back. this exercise is excellent for shoulder health and posture. return with control.', 'back', 'shoulders', now(), now()),

(gen_random_uuid(), 'straight arm pulldown', 'stand at a cable station with bar set high. grip the bar with straight arms. keeping arms straight, pull the bar down in an arc toward your thighs. this isolates the lats by removing bicep involvement. return to starting position with control, feeling the stretch in your lats.', 'back', null, now(), now()),

(gen_random_uuid(), 'inverted row', 'set a barbell in a rack at waist height. lie underneath and grip the bar with hands shoulder-width apart. keep body straight and pull your chest to the bar, squeezing shoulder blades together. lower with control. this is a great bodyweight back exercise, easier than pull-ups.', 'back', 'biceps', now(), now()),

(gen_random_uuid(), 'rack pulls', 'set a barbell on safety pins in a rack at knee height. stand close to the bar and grip it outside your legs. keep back flat and drive through heels to stand up straight. this partial deadlift variation emphasizes the upper back and traps while allowing heavier loads.', 'back', 'traps', now(), now()),

(gen_random_uuid(), 'shrugs', 'stand holding a barbell or dumbbells at your sides. elevate your shoulders straight up toward your ears as high as possible. hold briefly at the top, then lower with control. avoid rolling your shoulders. this exercise targets the trapezius muscles for building the upper back and neck area.', 'traps', null, now(), now()),

(gen_random_uuid(), 'hyperextensions', 'position yourself in a hyperextension bench with hips on the pad and ankles secured. cross arms over chest or behind head. lower your torso toward the ground by bending at the waist. raise back up until body forms a straight line. this targets the lower back, glutes, and hamstrings.', 'back', 'glutes', now(), now()),

-- shoulder exercises (15 exercises)
(gen_random_uuid(), 'barbell overhead press', 'stand with feet shoulder-width apart, barbell at shoulder height in front of you. press the bar straight overhead until arms are fully extended, moving your head slightly back to let the bar pass. lower with control back to shoulders. keep your core tight and avoid excessive back arching.', 'shoulders', 'triceps', now(), now()),

(gen_random_uuid(), 'seated barbell overhead press', 'sit on a bench with back support. position barbell at shoulder height. press overhead until arms are extended. lower with control. the seated position provides more stability and isolates the shoulders by reducing leg drive and momentum.', 'shoulders', 'triceps', now(), now()),

(gen_random_uuid(), 'dumbbell overhead press', 'stand or sit holding dumbbells at shoulder level, palms facing forward. press the dumbbells up until arms are fully extended. lower with control back to shoulders. dumbbells allow for natural arm movement and greater range of motion compared to barbells.', 'shoulders', 'triceps', now(), now()),

(gen_random_uuid(), 'arnold press', 'sit holding dumbbells at shoulder height with palms facing you. as you press up, rotate your hands so palms face forward at the top. reverse the motion coming down. this variation, named after arnold schwarzenegger, works all three deltoid heads through rotation.', 'shoulders', 'triceps', now(), now()),

(gen_random_uuid(), 'lateral raises', 'stand holding dumbbells at your sides. with a slight bend in your elbows, raise the weights out to the sides until they reach shoulder height. lead with your elbows, not your hands. lower with control. avoid using momentum or swinging the weights. this isolates the side deltoids.', 'shoulders', null, now(), now()),

(gen_random_uuid(), 'cable lateral raises', 'stand sideways to a cable station set at the lowest position. grab the handle with the far hand across your body. raise the cable out to the side until arm is parallel to the ground. lower with control. cables provide constant tension throughout the movement.', 'shoulders', null, now(), now()),

(gen_random_uuid(), 'front raises', 'stand holding dumbbells in front of your thighs. with straight or slightly bent arms, raise one or both weights forward and up to shoulder level. lower with control. this targets the front deltoids. avoid leaning back or using momentum.', 'shoulders', null, now(), now()),

(gen_random_uuid(), 'barbell front raises', 'stand holding a barbell at your thighs with an overhand grip. with arms straight or slightly bent, raise the bar forward and up to shoulder level. lower with control. the barbell allows you to use more weight than dumbbells for this movement.', 'shoulders', null, now(), now()),

(gen_random_uuid(), 'rear delt flyes', 'bend forward at the hips until torso is nearly parallel to ground. hold dumbbells hanging down with slight elbow bend. raise the weights out to the sides in an arc, squeezing your rear delts. lower with control. this targets the often-neglected posterior deltoids.', 'shoulders', 'back', now(), now()),

(gen_random_uuid(), 'seated rear delt flyes', 'sit on the edge of a bench, bend forward at the hips. hold dumbbells under your legs with slight elbow bend. raise the weights out to the sides, squeezing your rear delts at the top. the seated position provides stability and isolates the posterior deltoids.', 'shoulders', 'back', now(), now()),

(gen_random_uuid(), 'face pulls', 'set a cable at upper chest height with a rope attachment. pull the rope toward your face, separating the ends as you pull. focus on squeezing your rear delts and upper back. this exercise is excellent for shoulder health and posture.', 'shoulders', 'back', now(), now()),

(gen_random_uuid(), 'upright rows', 'stand holding a barbell or dumbbells in front of you. pull the weight straight up along your body until it reaches chest level, leading with your elbows. lower with control. keep the weight close to your body. this works the side and front delts as well as the traps.', 'shoulders', 'traps', now(), now()),

(gen_random_uuid(), 'machine shoulder press', 'sit at a shoulder press machine with back against the pad. grip handles at shoulder level. press up until arms are extended. lower with control. machines provide a stable pressing path, useful for beginners or high-rep finishing work.', 'shoulders', 'triceps', now(), now()),

(gen_random_uuid(), 'pike push-ups', 'start in a downward dog yoga position with hips high in the air. bend elbows to lower the top of your head toward the ground. press back up. this bodyweight exercise targets shoulders similar to an overhead press. progress toward handstand push-ups.', 'shoulders', 'triceps', now(), now()),

(gen_random_uuid(), 'plate raises', 'hold a weight plate with both hands at the 3 and 9 o''clock positions in front of your thighs. raise the plate up to shoulder level keeping arms extended. lower with control. this provides a different stimulus than dumbbells and can be done with various plate positions.', 'shoulders', null, now(), now()),

-- leg exercises (25 exercises)
(gen_random_uuid(), 'barbell back squat', 'position barbell on your upper back. stand with feet shoulder-width apart, toes slightly pointed out. descend by bending at hips and knees simultaneously, keeping your chest up and knees tracking over toes. lower until thighs are at least parallel to the ground. drive through your heels to stand back up.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'barbell front squat', 'position barbell across the front of your shoulders, elbows high. stand with feet shoulder-width apart. descend into a squat while keeping torso upright and elbows up. this variation places more emphasis on the quadriceps and requires better mobility and core stability.', 'quadriceps', 'core', now(), now()),

(gen_random_uuid(), 'goblet squat', 'hold a dumbbell or kettlebell at chest level with both hands. stand with feet shoulder-width apart. squat down, keeping chest up and weight close to body. this is excellent for learning squat mechanics and building quad strength. drive through heels to stand back up.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'bulgarian split squat', 'stand a few feet in front of a bench. place the top of one foot on the bench behind you. lower your body by bending the front knee until the back knee nearly touches the ground. drive through the front heel to return. this single-leg exercise builds balance and leg strength.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'romanian deadlift', 'hold a barbell at hip level with a slight knee bend. hinge at the hips, pushing them back while keeping your back straight. lower the bar along your legs until you feel a stretch in your hamstrings. return to starting position by driving your hips forward. this primarily targets the hamstrings and glutes.', 'hamstrings', 'glutes', now(), now()),

(gen_random_uuid(), 'dumbbell romanian deadlift', 'hold dumbbells at your sides with a slight knee bend. hinge at the hips, pushing them back while keeping your back straight. lower the dumbbells along your legs until you feel a stretch in your hamstrings. return by driving hips forward. dumbbells allow for natural arm position.', 'hamstrings', 'glutes', now(), now()),

(gen_random_uuid(), 'leg press', 'sit in the leg press machine with feet shoulder-width apart on the platform. release the safety handles and lower the platform by bending your knees until they reach about 90 degrees. press back up through your heels until legs are nearly extended, but don''t lock your knees.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'hack squat', 'position yourself in a hack squat machine with back against the pad and shoulders under the supports. place feet shoulder-width apart on the platform. lower by bending your knees until thighs are parallel to the platform. press back up. this machine variation targets quads with less back stress.', 'quadriceps', null, now(), now()),

(gen_random_uuid(), 'walking lunges', 'stand with feet together. step forward with one leg and lower your body until both knees are bent at 90 degrees. your rear knee should nearly touch the ground. push through your front heel to step forward with the opposite leg. continue alternating legs.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'reverse lunges', 'stand with feet together. step backward with one leg and lower your body until both knees are bent at 90 degrees. push through your front heel to return to starting position. this variation is easier on the knees than forward lunges and emphasizes the glutes more.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'stationary lunges', 'start in a split stance with one foot forward and one back. lower your body by bending both knees until the rear knee nearly touches the ground. push back up to starting position without moving your feet. complete all reps on one side before switching.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'leg curl', 'lie face down on a leg curl machine with ankles under the pad. curl your legs up toward your glutes, squeezing your hamstrings at the top. lower with control to the starting position. avoid lifting your hips off the bench. this isolates the hamstrings.', 'hamstrings', null, now(), now()),

(gen_random_uuid(), 'seated leg curl', 'sit in a leg curl machine with the pad against your lower calves. curl your legs down and under, squeezing your hamstrings. return to starting position with control. the seated position changes the hamstring activation compared to lying curls.', 'hamstrings', null, now(), now()),

(gen_random_uuid(), 'leg extension', 'sit in a leg extension machine with shins behind the pad. extend your legs until they are straight, squeezing your quadriceps at the top. lower with control back to starting position. this exercise isolates the quadriceps muscles.', 'quadriceps', null, now(), now()),

(gen_random_uuid(), 'step-ups', 'stand facing a bench or box. step up onto the platform with one foot, driving through that heel to bring the other foot up. step back down with control. alternate legs or complete all reps on one side first. this builds single-leg strength and balance.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'box squats', 'set up a box or bench behind you at a height where thighs are parallel when seated. perform a squat, sitting back onto the box briefly before driving back up. this teaches proper squat mechanics and can help build explosive power from a dead stop.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'sumo squat', 'stand with feet wider than shoulder-width, toes pointed out at 45 degrees. hold a dumbbell or kettlebell at your chest. squat down, keeping chest up and knees tracking over toes. this wide stance emphasizes the inner thighs and glutes more than standard squats.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'hip thrust', 'sit on the ground with your upper back against a bench. place a barbell across your hips. drive through your heels to lift your hips up until your body forms a straight line from shoulders to knees. squeeze glutes at the top. lower with control. this is the best exercise for glute development.', 'glutes', 'hamstrings', now(), now()),

(gen_random_uuid(), 'glute bridge', 'lie on your back with knees bent and feet flat on the ground. drive through your heels to lift your hips up until body forms a straight line from shoulders to knees. squeeze glutes at the top. lower with control. this bodyweight exercise targets glutes and hamstrings.', 'glutes', 'hamstrings', now(), now()),

(gen_random_uuid(), 'single leg glute bridge', 'lie on your back with one knee bent, foot flat, and the other leg extended. drive through the grounded heel to lift your hips up. keep hips level throughout. this single-leg variation increases difficulty and helps address imbalances.', 'glutes', 'hamstrings', now(), now()),

(gen_random_uuid(), 'standing calf raises', 'stand on a raised surface with the balls of your feet on the edge and heels hanging off. rise up onto your toes as high as possible, squeezing your calves at the top. lower your heels below the level of the step for a full stretch. can be done on a machine or with dumbbells for added weight.', 'calves', null, now(), now()),

(gen_random_uuid(), 'seated calf raises', 'sit at a calf raise machine with the pads on your thighs and balls of feet on the platform. lower your heels as far as possible for a stretch. raise up onto your toes as high as possible. the seated position emphasizes the soleus muscle of the calf.', 'calves', null, now(), now()),

(gen_random_uuid(), 'good mornings', 'stand with a barbell across your upper back. with a slight knee bend, hinge forward at the hips, pushing them back while keeping your back straight. lower until torso is nearly parallel to ground. return to starting position. this targets hamstrings, glutes, and lower back.', 'hamstrings', 'back', now(), now()),

(gen_random_uuid(), 'sissy squat', 'stand with feet shoulder-width apart, rise onto your toes. lean back while bending your knees forward, lowering your body. keep your torso and thighs in a straight line. this advanced exercise intensely isolates the quadriceps. hold something for balance if needed.', 'quadriceps', null, now(), now()),

(gen_random_uuid(), 'wall sit', 'lean against a wall and slide down until your thighs are parallel to the ground, as if sitting in an invisible chair. keep your back flat against the wall. hold this position for time. this isometric exercise builds quad endurance and mental toughness.', 'quadriceps', null, now(), now()),

-- arm exercises (15 exercises)
(gen_random_uuid(), 'barbell curl', 'stand holding a barbell with an underhand grip at hip level. keeping your elbows stationary at your sides, curl the weight up toward your shoulders. squeeze your biceps at the top, then lower with control. avoid swinging or using your back.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'ez bar curl', 'stand holding an ez bar with an underhand grip. curl the bar up toward your shoulders, keeping elbows stationary. the angled grip of the ez bar is easier on the wrists than a straight bar. lower with control, fully extending arms at the bottom.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'dumbbell curl', 'stand holding dumbbells at your sides with palms facing forward. curl the weights up toward your shoulders, keeping elbows stationary. squeeze biceps at the top. lower with control. can be done alternating arms or both together.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'alternating dumbbell curl', 'stand holding dumbbells at your sides. curl one dumbbell up while keeping the other at your side. lower it, then curl the other. this allows you to focus on each arm independently and use slightly heavier weights than doing both together.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'hammer curl', 'stand holding dumbbells at your sides with palms facing each other (neutral grip). curl the weights up toward your shoulders, maintaining the neutral grip throughout. this variation targets the biceps and forearms. lower with control.', 'biceps', 'forearms', now(), now()),

(gen_random_uuid(), 'preacher curl', 'sit at a preacher curl bench with arms resting on the pad. hold a barbell or dumbbells with an underhand grip. curl the weight up, squeezing your biceps at the top. lower with control. the pad prevents momentum and isolates the biceps.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'concentration curl', 'sit on a bench with legs spread. hold a dumbbell in one hand, resting your elbow against the inside of your thigh. curl the weight up toward your shoulder. lower with control. this seated isolation exercise allows maximum focus on the bicep.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'cable curl', 'stand at a cable station with a straight bar attachment set at the low position. grip the bar with an underhand grip. curl the bar up toward your shoulders. lower with control. cables provide constant tension throughout the entire range of motion.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'incline dumbbell curl', 'sit on an incline bench set to 45-60 degrees with dumbbells hanging at your sides. curl the weights up, keeping elbows back. the incline puts the biceps in a stretched position at the start, increasing the range of motion and targeting the long head.', 'biceps', null, now(), now()),

(gen_random_uuid(), 'close grip bench press', 'lie on a flat bench. grip the barbell with hands about shoulder-width or slightly narrower. lower the bar to your chest, keeping elbows closer to your body than regular bench press. press back up. this variation emphasizes the triceps while still working the chest.', 'triceps', 'chest', now(), now()),

(gen_random_uuid(), 'tricep dips', 'position yourself on parallel bars or between two benches. lower your body by bending your elbows until your upper arms are roughly parallel to the ground. press back up to starting position. keep your body as upright as possible to emphasize triceps. lean forward slightly to involve more chest.', 'triceps', 'chest', now(), now()),

(gen_random_uuid(), 'bench dips', 'sit on the edge of a bench with hands gripping the edge beside your hips. slide your hips forward off the bench with legs extended. lower your body by bending your elbows to about 90 degrees. press back up. this bodyweight tricep exercise can be made harder by elevating your feet.', 'triceps', null, now(), now()),

(gen_random_uuid(), 'tricep pushdown', 'stand at a cable machine with a straight or rope attachment set at upper chest height. with elbows pinned to your sides, push the handle down until your arms are fully extended. squeeze your triceps at the bottom. return to starting position with control, stopping when your forearms are about parallel to the ground.', 'triceps', null, now(), now()),

(gen_random_uuid(), 'rope tricep pushdown', 'stand at a cable machine with a rope attachment. grip the rope with palms facing each other. push down until arms are extended, then pull the rope ends apart at the bottom for maximum tricep contraction. return to starting position with control.', 'triceps', null, now(), now()),

(gen_random_uuid(), 'overhead tricep extension', 'hold a dumbbell or use a cable attachment overhead with both hands. keeping your upper arms stationary and elbows pointed forward, lower the weight behind your head by bending your elbows. extend your arms back to starting position, squeezing your triceps. can be done standing or seated.', 'triceps', null, now(), now()),

(gen_random_uuid(), 'skull crushers', 'lie on a flat bench holding an ez bar or barbell above your chest with arms extended. keeping upper arms stationary, bend elbows to lower the bar toward your forehead. extend arms back to starting position. this exercise intensely targets the triceps, particularly the long head.', 'triceps', null, now(), now()),

-- core exercises (20 exercises)
(gen_random_uuid(), 'plank', 'start in a forearm plank position with elbows directly under shoulders. keep your body in a straight line from head to heels, engaging your core. hold this position without letting your hips sag or rise. focus on breathing steadily throughout. this isometric exercise builds core endurance.', 'core', null, now(), now()),

(gen_random_uuid(), 'side plank', 'lie on your side with forearm on the ground, elbow under shoulder. stack your feet and lift your hips so your body forms a straight line. hold this position, keeping hips up. this targets the obliques and helps build lateral core stability. switch sides and repeat.', 'core', 'obliques', now(), now()),

(gen_random_uuid(), 'plank with shoulder taps', 'start in a high plank position (hands instead of forearms). while maintaining a stable core and hips, lift one hand to tap the opposite shoulder. alternate sides. this dynamic plank variation adds anti-rotation work and shoulder stability.', 'core', 'shoulders', now(), now()),

(gen_random_uuid(), 'mountain climbers', 'start in a high plank position. quickly alternate bringing each knee toward your chest in a running motion. keep your hips down and core engaged. this dynamic exercise combines core work with cardiovascular conditioning.', 'core', 'cardiovascular', now(), now()),

(gen_random_uuid(), 'crunches', 'lie on your back with knees bent and feet flat on the ground. place hands behind your head or across your chest. curl your upper body toward your knees, lifting your shoulder blades off the ground. squeeze your abs at the top, then lower with control. avoid pulling on your neck.', 'core', null, now(), now()),

(gen_random_uuid(), 'bicycle crunches', 'lie on your back with hands behind your head and legs lifted with knees bent. bring your right elbow toward your left knee while extending your right leg. alternate sides in a pedaling motion. this dynamic exercise targets the entire abdominal region, especially the obliques.', 'core', 'obliques', now(), now()),

(gen_random_uuid(), 'reverse crunches', 'lie on your back with hands at your sides or under your hips. bend knees at 90 degrees. curl your hips up off the ground, bringing knees toward your chest. lower with control. this variation targets the lower abs more than regular crunches.', 'core', null, now(), now()),

(gen_random_uuid(), 'russian twists', 'sit on the ground with knees bent and feet elevated. lean back slightly to engage your core. hold a weight or medicine ball at chest level. rotate your torso from side to side, touching the weight to the ground on each side. keep your core tight throughout. this targets the obliques.', 'core', 'obliques', now(), now()),

(gen_random_uuid(), 'sit-ups', 'lie on your back with knees bent and feet flat or anchored. place hands behind your head or across your chest. curl your entire torso up toward your knees until you reach a seated position. lower back down with control. this full range motion works the entire abdominal wall.', 'core', null, now(), now()),

(gen_random_uuid(), 'hanging leg raises', 'hang from a pull-up bar with arms fully extended. keeping your legs straight or slightly bent, raise them up toward your chest by flexing your hips and core. lower with control to starting position. avoid swinging or using momentum. this is an advanced core exercise.', 'core', 'hip flexors', now(), now()),

(gen_random_uuid(), 'hanging knee raises', 'hang from a pull-up bar with arms fully extended. bend your knees and pull them up toward your chest. lower with control. this is an easier variation of leg raises, good for building toward the full movement. focus on using your abs rather than just hip flexors.', 'core', 'hip flexors', now(), now()),

(gen_random_uuid(), 'dead bug', 'lie on your back with arms extended toward the ceiling and knees bent at 90 degrees. slowly lower your right arm overhead while extending your left leg, hovering just above the ground. return to starting position and repeat on opposite side. keep your lower back pressed to the ground throughout. this teaches core stability.', 'core', null, now(), now()),

(gen_random_uuid(), 'bird dog', 'start on hands and knees in a tabletop position. extend your right arm forward and left leg back simultaneously, forming a straight line. hold briefly, then return to start. alternate sides. this exercise builds core stability and teaches proper spine position under movement.', 'core', 'back', now(), now()),

(gen_random_uuid(), 'ab wheel rollout', 'kneel on the ground holding an ab wheel. roll the wheel forward, extending your body while keeping your core tight. roll out as far as you can while maintaining a neutral spine. pull yourself back to starting position using your abs. this is an advanced, intense core exercise.', 'core', null, now(), now()),

(gen_random_uuid(), 'cable woodchop', 'stand sideways to a cable station set at shoulder height. grip the handle with both hands. pull the cable down and across your body in a diagonal chopping motion, rotating your torso. return with control. this dynamic movement targets the obliques and teaches rotational power.', 'core', 'obliques', now(), now()),

(gen_random_uuid(), 'pallof press', 'stand sideways to a cable station set at chest height. hold the handle at your chest with both hands. press the handle straight out in front of you, resisting the cable trying to pull you sideways. pull back to chest. this anti-rotation exercise builds core stability.', 'core', null, now(), now()),

(gen_random_uuid(), 'v-ups', 'lie flat on your back with arms extended overhead. simultaneously lift your legs and upper body, reaching your hands toward your feet, forming a v shape. lower back down with control. this advanced exercise works the entire abdominal region dynamically.', 'core', null, now(), now()),

(gen_random_uuid(), 'toe touches', 'lie on your back with legs extended straight up toward the ceiling. reach your hands up toward your toes, lifting your shoulder blades off the ground. lower with control. this exercise targets the upper abs while the fixed leg position provides constant tension.', 'core', null, now(), now()),

(gen_random_uuid(), 'flutter kicks', 'lie on your back with hands under your hips and legs extended. lift your legs slightly off the ground. alternate kicking your legs up and down in small, quick motions. keep your lower back pressed to the ground. this targets the lower abs and hip flexors.', 'core', 'hip flexors', now(), now()),

(gen_random_uuid(), 'hollow body hold', 'lie on your back and lift your shoulders and legs off the ground, pressing your lower back into the floor. extend arms overhead or by your sides. hold this hollow position. this gymnastics-based exercise builds incredible core strength and body awareness.', 'core', null, now(), now()),

-- compound/functional/full body (15 exercises)
(gen_random_uuid(), 'burpees', 'start standing. drop into a squat and place your hands on the ground. jump your feet back into a plank position. perform a push-up. jump your feet back to your hands. explosively jump up with arms overhead. this full-body exercise combines strength and cardio for conditioning.', 'full body', 'cardiovascular', now(), now()),

(gen_random_uuid(), 'kettlebell swing', 'stand with feet shoulder-width apart, holding a kettlebell with both hands. hinge at your hips to swing the kettlebell between your legs. explosively drive your hips forward, swinging the kettlebell to shoulder height. let momentum bring it back down. this is a powerful posterior chain and conditioning exercise.', 'glutes', 'hamstrings', now(), now()),

(gen_random_uuid(), 'box jumps', 'stand facing a sturdy box or platform. bend your knees and swing your arms back. explosively jump onto the box, landing softly with knees bent. step back down with control. this develops explosive power in the legs. start with a lower height and progress gradually.', 'quadriceps', 'glutes', now(), now()),

(gen_random_uuid(), 'thrusters', 'hold dumbbells or a barbell at shoulder height. perform a front squat, then as you stand up, press the weight overhead in one fluid motion. lower back to shoulders and repeat. this combines a squat and press for a full-body conditioning exercise.', 'full body', 'shoulders', now(), now()),

(gen_random_uuid(), 'clean and press', 'start with a barbell on the ground. explosively lift it to your shoulders in one motion (the clean). then press it overhead. lower to shoulders, then to ground. this olympic lift variation builds total body power and coordination.', 'full body', 'shoulders', now(), now()),

(gen_random_uuid(), 'power clean', 'start with a barbell on the ground. in one explosive movement, pull the bar up while simultaneously dropping under it to catch it at shoulder height in a quarter squat. stand up fully. this olympic lift develops explosive power throughout the entire body.', 'full body', 'back', now(), now()),

(gen_random_uuid(), 'turkish get-up', 'lie on your back holding a kettlebell or dumbbell straight up. through a series of movements, stand up while keeping the weight overhead the entire time. reverse the movement to return to lying. this complex movement builds stability, mobility, and full-body strength.', 'full body', 'shoulders', now(), now()),

(gen_random_uuid(), 'farmer''s walk', 'hold heavy dumbbells, kettlebells, or farmers walk handles at your sides. walk forward with good posture, taking controlled steps. this simple but brutal exercise builds grip strength, core stability, and overall conditioning. keep shoulders back and core tight.', 'full body', 'forearms', now(), now()),

(gen_random_uuid(), 'bear crawl', 'start on hands and knees, then lift knees slightly off the ground. crawl forward by moving opposite hand and foot together. keep your hips low and core tight. this primal movement pattern works the entire body and builds conditioning.', 'full body', 'core', now(), now()),

(gen_random_uuid(), 'medicine ball slams', 'hold a medicine ball overhead. explosively throw it down to the ground as hard as possible, engaging your entire body. catch it on the bounce or pick it up and repeat. this develops power and provides a great conditioning stimulus.', 'full body', 'core', now(), now()),

(gen_random_uuid(), 'battle ropes', 'hold the ends of heavy ropes with both hands. create waves by alternating raising and lowering each arm rapidly. maintain athletic stance with slight knee bend. this conditioning exercise works the arms, shoulders, and core while elevating heart rate.', 'full body', 'shoulders', now(), now()),

(gen_random_uuid(), 'sled push', 'place hands on a weighted sled at shoulder height. drive through your legs to push the sled forward, keeping your core tight and body at an angle. this develops lower body power and conditioning without the eccentric stress of other exercises.', 'full body', 'quadriceps', now(), now()),

(gen_random_uuid(), 'sled pull', 'attach a rope or strap to a weighted sled. hold the rope and walk backward, pulling the sled toward you. this targets the posterior chain and builds conditioning. can also be done facing forward for different muscle emphasis.', 'full body', 'back', now(), now()),

(gen_random_uuid(), 'wall balls', 'hold a medicine ball at chest level. perform a squat, then explosively stand and throw the ball to a target on the wall above you (typically 10 feet high). catch it on the return and immediately go into the next rep. this crossfit staple combines leg and shoulder work with conditioning.', 'full body', 'shoulders', now(), now()),

(gen_random_uuid(), 'rowing machine', 'sit at a rowing machine with feet strapped in. push with legs first, then lean back slightly and pull the handle to your chest. reverse the motion with control. this full-body cardio exercise works the legs, back, arms, and core while providing excellent conditioning.', 'full body', 'back', now(), now())
-- +goose statementend

-- +goose Down
drop table exercises;
