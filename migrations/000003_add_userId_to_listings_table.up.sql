-- intially since we hvae records already and even in prod we have to first make nullable and then ru a script or query that populates all userid to lostings and then make it not null by running another migration
-- so now this is first migration to create Nullable userid in listings table

ALTER TABLE listings
ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;



-- for effiecnt search when filtered by userid on litsings table
CREATE INDEX idx_listings_user_id ON listings(user_id); 

