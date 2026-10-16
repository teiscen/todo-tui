-- Sample data for todo-tui. Safe to re-run: it first clears the todo tables
-- (labels and entries are left alone). Uses explicit ids.
DELETE FROM links;
DELETE FROM time_log;
DELETE FROM step_notes;
DELETE FROM steps;
DELETE FROM completions;
DELETE FROM recurrence;
DELETE FROM todos;
-- Dates are relative to "now", so due/overdue/streak states always look current.

INSERT OR IGNORE INTO labels (name, color) VALUES
  ('personal','#81c8be'), ('work','#ca9ee6'), ('health','#ef9f76'), ('errands','#babbf1'), ('home','#f2d5cf');

-- ── One-time, small ─────────────────────────────────────────────
INSERT INTO todos (id,title,label_id,horizon,due_date,notes) VALUES
 (1,'Write lab report','personal','small',date('now','-2 days'),'Need figures from lab 3'||char(10)||'TA office hours Fri 2-3pm'||char(10)||'Template: course wiki, week 4'),
 (3,'Email supervisor','work','small',date('now'),''),
 (4,'Order gym shoes','health','small',date('now','+6 days'),'Check size 10.5 in the blue pair'),
 (5,'Book dentist',NULL,'small',date('now','+12 days'),''),
 (6,'Pay tuition','personal','small',date('now','+20 days'),'Portal closes at 4pm on the due date'),
 (9,'Laundry',NULL,'small',NULL,'');

-- ── One-time, long-term (with steps) ────────────────────────────
INSERT INTO todos (id,title,label_id,horizon,due_date,notes) VALUES
 (2,'Apply to co-op','errands','long',date('now','+3 days'),'Target 5 postings'),
 (7,'Midterm prep','personal','long',date('now','+8 days'),'Ch 3-7 are fair game'),
 (8,'Plan volunteer shifts','home','long',date('now','+21 days'),'');

-- ── Completed ──────────────────────────────────────────────────
INSERT INTO todos (id,title,label_id,horizon,done_at) VALUES
 (10,'Renew bus pass',NULL,'small',strftime('%Y-%m-%dT%H:%M:%SZ','now','-1 day')),
 (11,'Return library books',NULL,'small',strftime('%Y-%m-%dT%H:%M:%SZ','now','-3 days')),
 (16,'Submit form 3','personal','small',strftime('%Y-%m-%dT%H:%M:%SZ','now','-5 days'));

-- ── Recurring ──────────────────────────────────────────────────
INSERT INTO todos (id,title,label_id,horizon,notes) VALUES
 (12,'Gym','health','small',''),
 (13,'Read 20 minutes',NULL,'small',''),
 (14,'Review notes','personal','small',''),
 (15,'Water plants',NULL,'small','');

INSERT INTO recurrence (todo_id,kind,weekdays,interval_days,start_date) VALUES
 (12,'weekdays',42,NULL,date('now','-60 days')),   -- Mon+Wed+Fri (bit0=Sun)
 (13,'daily',0,NULL,date('now','-60 days')),
 (14,'weekdays',1,NULL,date('now','-60 days')),    -- Sundays
 (15,'interval',0,3,date('now','-27 days'));       -- every 3 days

-- ── Steps ──────────────────────────────────────────────────────
INSERT INTO steps (todo_id,pos,text,status,started_at,done_at,due_date) VALUES
 (1,1,'outline','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-9 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-9 days'),NULL),
 (1,2,'intro','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-8 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-8 days'),NULL),
 (1,3,'methods','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-7 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-7 days'),NULL),
 (1,4,'results','started',strftime('%Y-%m-%dT%H:%M:%SZ','now'),NULL,NULL),
 (1,5,'proofread','todo',NULL,NULL,date('now','+1 day')),

 (2,1,'update resume','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-4 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-3 days'),NULL),
 (2,2,'shortlist postings','todo',NULL,NULL,date('now','+1 day')),
 (2,3,'tailor cover letters','todo',NULL,NULL,date('now','+2 days')),
 (2,4,'ask for references','todo',NULL,NULL,NULL),
 (2,5,'submit applications','todo',NULL,NULL,date('now','+3 days')),

 (7,1,'gather past exams','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-3 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-3 days'),NULL),
 (7,2,'summarize ch 3-4','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-2 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-1 day'),NULL),
 (7,3,'summarize ch 5-7','started',strftime('%Y-%m-%dT%H:%M:%SZ','now'),NULL,date('now','+3 days')),
 (7,4,'practice exam 1','todo',NULL,NULL,date('now','+5 days')),
 (7,5,'practice exam 2','todo',NULL,NULL,date('now','+7 days')),
 (7,6,'review mistakes','todo',NULL,NULL,NULL),

 (8,1,'collect availability','done',strftime('%Y-%m-%dT%H:%M:%SZ','now','-2 days'),strftime('%Y-%m-%dT%H:%M:%SZ','now','-2 days'),NULL),
 (8,2,'draft rota','todo',NULL,NULL,date('now','+7 days')),
 (8,3,'send to coordinator','todo',NULL,NULL,date('now','+14 days')),
 (8,4,'confirm with volunteers','todo',NULL,NULL,date('now','+20 days'));

INSERT INTO step_notes (step_id,body)
 SELECT id,'used template from course' FROM steps WHERE todo_id=1 AND pos=1;
INSERT INTO step_notes (step_id,body)
 SELECT id,'need figures from lab 3' FROM steps WHERE todo_id=1 AND pos=4;
INSERT INTO step_notes (step_id,body)
 SELECT id,'ask TA about error bars' FROM steps WHERE todo_id=1 AND pos=4;
INSERT INTO step_notes (step_id,body)
 SELECT id,'aim for postings closing first' FROM steps WHERE todo_id=2 AND pos=2;

INSERT INTO links (todo_id,label,url) VALUES
 (1,'course wiki / week 4 template','https://example.com/course/wiki/week4'),
 (1,'lab 3 data (shared drive)','https://example.com/drive/lab3'),
 (2,'co-op posting board','https://example.com/coop/postings');

INSERT INTO time_log (todo_id,step_id,date,minutes)
 SELECT 1,id,date('now','-9 days'),65 FROM steps WHERE todo_id=1 AND pos=1;
INSERT INTO time_log (todo_id,step_id,date,minutes)
 SELECT 1,id,date('now','-8 days'),55 FROM steps WHERE todo_id=1 AND pos=2;
INSERT INTO time_log (todo_id,step_id,date,minutes)
 SELECT 1,id,date('now','-7 days'),70 FROM steps WHERE todo_id=1 AND pos=3;
INSERT INTO time_log (todo_id,step_id,date,minutes)
 SELECT 1,id,date('now'),25 FROM steps WHERE todo_id=1 AND pos=4;
INSERT INTO time_log (todo_id,step_id,date,minutes)
 SELECT 7,id,date('now','-2 days'),90 FROM steps WHERE todo_id=7 AND pos=2;

-- ── Completions (last ~3 weeks) ────────────────────────────────
WITH RECURSIVE d(n) AS (SELECT 0 UNION ALL SELECT n+1 FROM d WHERE n<20)
INSERT INTO completions (todo_id,date)                       -- Gym: Mon/Wed/Fri, with 2 misses
 SELECT 12,date('now','-'||n||' days') FROM d
 WHERE strftime('%w',date('now','-'||n||' days')) IN ('1','3','5') AND n NOT IN (5,12,13);

WITH RECURSIVE d(n) AS (SELECT 0 UNION ALL SELECT n+1 FROM d WHERE n<13)
INSERT INTO completions (todo_id,date)                       -- Read: daily, one gap on day 7
 SELECT 13,date('now','-'||n||' days') FROM d WHERE n<>7;

WITH RECURSIVE d(n) AS (SELECT 0 UNION ALL SELECT n+1 FROM d WHERE n<27)
INSERT INTO completions (todo_id,date)                       -- Review notes: Sundays, skip one
 SELECT 14,date('now','-'||n||' days') FROM d
 WHERE strftime('%w',date('now','-'||n||' days'))='0' AND n NOT IN (7,8,9,10,11,12,13);

INSERT INTO completions (todo_id,date)                       -- Water plants: every 3rd day, skip one
 SELECT 15,date('now','-'||n||' days') FROM (SELECT 0 n UNION SELECT 3 UNION SELECT 6 UNION SELECT 9 UNION SELECT 15);