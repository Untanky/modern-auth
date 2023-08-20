import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from './schema';
import { env } from '$env/dynamic/private';

const pg = postgres(env.DB_URL);

export const db = drizzle(pg, {
  schema: schema,
});
