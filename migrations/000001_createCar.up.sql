BEGIN;

CREATE TABLE IF NOT EXISTS cars (
    id UUID NOT NULL,
    mark TEXT NOT NULL,
    model TEXT NOT NULL,
    color TEXT NOT NULL
    year INT NOT NULL,
    regNums TEXT NOT NULL,
    owner TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT pk_cars_idx PRIMARY KEY (id),
    CONSTRAINT cars_regNums_uniq_idx UNIQUE (regNums)
);

END;