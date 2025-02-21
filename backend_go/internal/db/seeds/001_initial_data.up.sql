-- Insert sample groups
INSERT INTO groups (name) VALUES
    ('Basic Vocabulary'),
    ('Common Phrases'),
    ('Greetings'),
    ('Numbers'),
    ('Colors');

-- Insert sample words
INSERT INTO words (arabic, romaji, english, parts) VALUES
    ('مرحبا', 'marhaba', 'hello', '{"type": "greeting", "formality": "neutral"}'),
    ('شكرا', 'shukran', 'thank you', '{"type": "courtesy", "formality": "neutral"}'),
    ('من فضلك', 'min fadlik', 'please', '{"type": "courtesy", "formality": "neutral"}'),
    ('واحد', 'wahid', 'one', '{"type": "number", "category": "cardinal"}'),
    ('اثنان', 'ithnan', 'two', '{"type": "number", "category": "cardinal"}'),
    ('ثلاثة', 'thalatha', 'three', '{"type": "number", "category": "cardinal"}'),
    ('أحمر', 'ahmar', 'red', '{"type": "color", "category": "basic"}'),
    ('أزرق', 'azraq', 'blue', '{"type": "color", "category": "basic"}'),
    ('أخضر', 'akhdar', 'green', '{"type": "color", "category": "basic"}'),
    ('صباح الخير', 'sabah al-khayr', 'good morning', '{"type": "greeting", "formality": "formal", "timeOfDay": "morning"}');

-- Link words to groups
INSERT INTO words_groups (word_id, group_id) 
SELECT w.id, g.id 
FROM words w, groups g 
WHERE 
    (w.english IN ('hello', 'good morning') AND g.name = 'Greetings')
    OR (w.english IN ('thank you', 'please') AND g.name = 'Common Phrases')
    OR (w.english IN ('one', 'two', 'three') AND g.name = 'Numbers')
    OR (w.english IN ('red', 'blue', 'green') AND g.name = 'Colors')
    OR g.name = 'Basic Vocabulary'; -- Add all words to Basic Vocabulary

-- Create some sample study sessions
INSERT INTO study_sessions (group_id) 
SELECT id FROM groups WHERE name IN ('Basic Vocabulary', 'Greetings', 'Numbers');

-- Create study activities for the sessions
INSERT INTO study_activities (study_session_id, group_id)
SELECT s.id, s.group_id
FROM study_sessions s;

-- Add some word review items with mixed results
INSERT INTO word_review_items (word_id, study_session_id, correct)
SELECT 
    w.id,
    s.id,
    CASE (ABS(RANDOM()) % 2) WHEN 0 THEN 0 ELSE 1 END
FROM words w
CROSS JOIN study_sessions s
WHERE w.id <= 5; -- Only add reviews for the first 5 words
