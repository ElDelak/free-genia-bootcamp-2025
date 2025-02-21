# Backend Server Technical Specs

## Business Goal :

A language learning school wants to build a prototype of learning which will act as three things :

- Inventory of possible vocabulary that can be learned.
- Act as a Learning record store (LRS), providing correct and wrong score or practice vocabulary
- A unified lauchpad to lauch different learning apps

## Technical reauirements

- The backend will be built using Go
- The database will be SQLite3
- The API will be built using GIN
- Mage is a task runner for Go
- The API will always return JSON
- There will no authentication or authorization
- Everything will be traited as a single user

## Directory structure

```text
backend_go/
├── main.go
├── go.mod
├── go.sum
├── internal/
│   ├── database/
│   │   ├── db.go
│   │   ├── schema.sql
│   │   └── words.db
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 0001_init.sql
│   │   │   ├── 0002_create_words_table.sql
│   │   │   └── 0003_create_groups_table.sql
│   │   └── seeds/
│   │       ├── words.json
│   │       └── groups.json
│   ├── handlers/
│   │   ├── dashboard.go
│   │   ├── study_sessions.go
│   │   └── words.go
│   ├── models/
│   │   ├── study_session.go
│   │   ├── word.go
│   │   └── word_review_item.go
│   ├── repositories/
│   │   ├── study_session_repository.go
│   │   ├── word_repository.go
│   │   └── word_review_item_repository.go
│   └── services/
│       ├── study_session_service.go
│       ├── word_service.go
│       └── word_review_item_service.go
├── magefile.go
└── pkg/
    ├── ginutils/
    │   └── ginutils.go
    └── sqliteutils/
        └── sqliteutils.go
```

## Database Schema

Our database will be a single sqlite3 database called "words.db" thal will be in the root of the project folder of "backend_go"

We have the following tables :

- words - stored vocabulary words
  - id integer
  - arabic string
  - romaji string
  - english string
  - parts json
- words-groups - join table for word and groups
  many-to-many
  - id integer
  - word_id integer
  - group_id integer
- groups - thematic groups of words
  - id integer
  - name string
- study_sessions - records of study sessions grouping word_review-items
  - id integer
  - group_id integer
  - created_at datetime
  - study_activity_id integer
- study_activities - a specific study activity, linking a study session to group
  - id integer
  - study_session_id integer
  - group_id integer
  - created_at datetime
- word_review-items - a record of word practice, determining if the word was correct or not
  - word_id integer
  - study_session_id integer
  - correct boolean
  - created_at datetime

###API Endpoints

### GET /api/dashboard/last_study_session

Returns information about the most recent study session.

#### Json response

```json
{
  "id": 123, // integer
  "study_activity_id": 456, // integer
  "group_id": 789, // integer
  "created_at": "2022-01-01 12:00:00", // datetime
  "group_name": "Basic Geering"
}
```

### GET /api/dashboard/study_progress

returns study progress statistics

#### Json response

```json
{
  "total_available_words": 80, // integer
  "total_words_studied": 5 // integer
}
```

### GET /api/dashboard/quick_stats

Return auick overview statistics.

#### Json response

```json
{
  "success_rate": 80, // integer (percentage)
  "total_study_sessions": 10, // integer
  "total_active_groups": 5, // integer
  "study_streak": 5 // integer (days)
}
```

### GET /api/study_activities/:id

Returns information about a specific study activity.

#### Json response

```json
{
  {
  "id": 123, // integer
  "name": "Study Activity Name", // string
  "description": "Study Activity Description", // string
  "thumbnail": "Study Activity Thumbnail URL", // string
  "launch_url": "Study Activity Launch URL" // string
}
}
```

### GET /api/study_activities/:id/study_sessions

- pagination with 100 items per page

```json
{
  "items": [
    {
      "id": 123, // integer
      "activity_name": "Vocabulary Study", // integer
      "group_name": "Bascic Geerings", // integer
      "start_time": "2022-01-01 12:00:00", // datetime
      "end_time": "2022-01-01 13:00:00", // datetime
      "num_review_items": 10 // integer
    }
  ],
  "pagination": {
    "current_page": 1, // integer
    "total_pages": 3, // integer
    "total_items": 30, // integer
    "items_per_page": 100 // integer
  }
}
```

### POST /api/study_Activities

#### request params

- group_id integer
- study_activity_id integer

#### Json response

```json
{
  "id": 123, // integer
  "group_id": 456 // integer
}
```

### GET /api/words

- pagination with 100 items per page

#### Json response

```json
{
  "items": [
    {
      "id": 123, // integer
      "arabic": "السلام عليكم", // string
      "romaji": "salam alikom", // string
      "english": "hello", // string
      "parts": [
        {
          "type": "verb", // string
          "text": "salam" // string
        }
      ]
    }
  ],
  "pagination": {
    "current_page": 1, // integer
    "total_pages": 3, // integer
    "total_items": 30, // integer
    "items_per_page": 100 // integer
  }
}
```

### GET /api/words/:id

#### Json response

```json
{
  "id": 123, // integer
  "arabic": "السلام عليكم", // string
  "romaji": "salam alikom", // string
  "english": "hello", // string
  "stats": {
    "correct_count": 10, // integer
    "wrong_count": 5 // integer
  },
  "groups": [
    {
      "id": 456, // integer
      "name": "Basic Geering" // string
    }
  ]
}
```

### GET /api/groups

- pagination with 100 items per page

#### Json response

```json
{
  "items": [
    {
      "id": 123, // integer
      "name": "Basic Geering", // string,
      "words_count": 10 // integer
    }
  ],
  "pagination": {
    "current_page": 1, // integer
    "total_pages": 3, // integer
    "total_items": 30, // integer
    "items_per_page": 100 // integer
  }
}
```

### GET /api/groups/:id

#### Json response

```json
{
  "id": 123, // integer
  "name": "Basic Geering", // string
  "stats": {
    "total_word_count": 10 // integer
  }
}
```

### GET /api/groups/:id/words

#### Json response

```json
{
  "items": [
    {
      "arabic": "السلام عليكم", // string
      "romaji": "salam alikom", // string
      "english": "hello", // string
      "stats": {
        "correct_count": 10, // integer
        "wrong_count": 5 // integer
      }
    }
  ],
  "pagination": {
    "current_page": 1, // integer
    "total_pages": 3, // integer
    "total_items": 30, // integer
    "items_per_page": 100 // integer
  }
}
```

### GET /api/groups/:id/study_sessions

#### Json response

```json
{
  "items": [
    {
      "id": 123, // integer
      "activity_name": "Vocabulary Study", // integer
      "group_name": "Bascic Geerings", // integer
      "start_time": "2022-01-01 12:00:00", // datetime
      "end_time": "2022-01-01 13:00:00", // datetime
      "review_items_count": 10 // integer
    }
  ],
  "pagination": {
    "current_page": 1, // integer
    "total_pages": 3, // integer
    "total_items": 30, // integer
    "items_per_page": 100 // integer
  }
}
```

### GET /api/study_sessions

- pagination with 100 items per page

#### Json response

```json
{
  "items": [
    {
      "id": 123, // integer
      "activity_name": "Vocabulary Study", // integer
      "group_name": "Bascic Geerings", // integer
      "start_time": "2022-01-01 12:00:00", // datetime
      "end_time": "2022-01-01 13:00:00", // datetime
      "review_items_count": 10 // integer
    }
  ],
  "pagination": [
    {
      "current_page": 1, // integer
      "total_pages": 3, // integer
      "total_items": 30, // integer
      "items_per_page": 100 // integer
  ]
  ]
}
```

### GET /api/study_sessions/:id

#### Json response

```json
{
  "id": 123, // integer integer
  "activity_name": "Vocabulary Study", // integer
  "group_name": "Bascic Geerings", // integer
  "start_time": "2022-01-01 12:00:00", // datetime
  "end_time": "2022-01-01 13:00:00", // datetime
  "review_items_count": 10 // integer
}
```

### GET /api/study_sessions/:id/words

- pagination with 100 items per page

#### Json response

```json
{
  "items": [
    {
      "arabic": "السلام عليكم", // string
      "romaji": "salam alikom", // string
      "english": "hello", // string
      "stats": {
        "correct_count": 10, // integer
        "wrong_count": 5 // integer
      }
    }
  ],
  "pagination": {
    "current_page": 1, // integer
    "total_pages": 3, // integer
    "total_items": 30, // integer
    "items_per_page": 100 // integer
  }
}
```

### POST /api/reset_history

#### Json response

```json
{
  "success": true,
  "message": "Study history has been reset"
}
```

### POST /api/full_reset

#### Json response

```json
{
  "success": true,
  "message": "System has been fully reset"
}
```

### POST /api/stydy_sessions/:id/words/:word_id/review

#### request params

- id (study_session_id) integer
- word_id integer
- correct boolean

#### Request Payload

```json
{
  "success": true
}
```

#### Json response

```json
{
"success": true,
"word_id": 1,
"study_session_id": 1,
"correct": true
"created_at": "2022-01-01 12:00:00"
}
```

## Task Runner Tasks

Mage is a task manager for GO.

### initialize database

this task will initialize the sqlite3 database called "words.db"

### migrate database

this task will run a series of migration sql files on the dqtqbqse
Migrations live in the migrations folder. The migration files will be run in order of their file name. The file names should looks like this:

```sql
0001_init.sql
0002_create_words_table.sql
```

### seed database

this task will import json files and transform then into target data for our database.
All seed files live in the "seeds" folder
All seed files should be loaded
In our task we should have DSL to specific each seed file and its expected group word name.

```json
[
  {
    "arabic": "السلام عليكم", // string
    "romaji": "salam alikom", // string
    "english": "hello" // string
  }
]
```
